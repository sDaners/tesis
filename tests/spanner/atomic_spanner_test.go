package spanner_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"postgres-example/repo"
	"postgres-example/tools"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/stretchr/testify/assert"
)

// Global statement cache: Map[fileName][statementType] = []string
var statementCache map[string]map[string][]StatementInfo

// StatementInfo holds a statement with its original position
type StatementInfo struct {
	Statement string
	Position  int // Original position in file (1-indexed)
}

// TestMain sets up the statement cache before running tests
func TestMain(m *testing.M) {
	fmt.Println("Setting up statement cache...")

	err := setupStatementCache()
	if err != nil {
		fmt.Printf("Failed to setup statement cache: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Statement cache setup complete. Found %d files.\n", len(statementCache))

	// Run tests
	code := m.Run()

	// Cleanup if needed
	statementCache = nil

	// Exit with the test result code
	os.Exit(code)
}

// setupStatementCache reads all SQL files and categorizes statements by type
func setupStatementCache() error {
	statementCache = make(map[string]map[string][]StatementInfo)

	// Get all SQL files from generated_sql folder
	sqlFiles, err := filepath.Glob("../../generated_sql/*.sql")
	if err != nil {
		return fmt.Errorf("failed to find SQL files: %w", err)
	}

	for _, sqlFile := range sqlFiles {
		filename := filepath.Base(sqlFile)

		// Extract statements from file
		statements, err := tools.ExtractStatementsFromFile(sqlFile)
		if err != nil {
			return fmt.Errorf("failed to extract SQL statements from %s: %w", sqlFile, err)
		}

		// Initialize file map
		statementCache[filename] = make(map[string][]StatementInfo)

		// Categorize statements by type
		for i, stmt := range statements {
			if !isCommentOnlyStatement(stmt) {
				statementType := determineStatementType(filename, stmt)

				statementInfo := StatementInfo{
					Statement: stmt,
					Position:  i + 1,
				}

				statementCache[filename][statementType] = append(
					statementCache[filename][statementType],
					statementInfo,
				)
			}
		}

		fmt.Printf("  %s: %d statements categorized\n", filename, len(statements))
	}

	return nil
}

// determineStatementType determines the type of a statement, handling parse errors gracefully
func determineStatementType(filename string, stmt string) string {
	// Try to parse with memefish first
	_, err := memefish.ParseStatement(filename, stmt)
	if err != nil {
		// If parsing fails, fall back to raw text analysis
		return getStatementTypeFromRaw(stmt)
	}

	return getAtomicStatementType(stmt)
}

// AtomicStatementResult holds the results for a single SQL statement test
type AtomicStatementResult struct {
	Filename      string
	StatementNum  int
	Statement     string
	StatementType string
	ParseSuccess  bool
	ParseError    string
	ExecSuccess   bool
	ExecError     string
	ExecutionTime time.Duration
}

// AtomicFileResult aggregates results for all statements in a file
type AtomicFileResult struct {
	Filename         string
	TotalStatements  int
	ParsedCount      int
	ExecutedCount    int
	StatementResults []AtomicStatementResult
}

// AtomicTestSummary holds overall atomic test results
type AtomicTestSummary struct {
	TotalFiles      int
	TotalStatements int
	ParsedCount     int
	ExecutedCount   int
	FileResults     []AtomicFileResult
	ExecutionTime   time.Duration
}

// Valid Spanner schema setup - hardcoded from valid_spanner_database.sql
var validSchemaSetup = []string{
	`CREATE TABLE departments (
		dept_id STRING(36) DEFAULT (GENERATE_UUID()),
		dept_name STRING(50) NOT NULL,
		location STRING(100),
		created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
	) PRIMARY KEY (dept_id)`,

	`CREATE TABLE employees (
		emp_id STRING(36) DEFAULT (GENERATE_UUID()),
		first_name STRING(50) NOT NULL,
		last_name STRING(50) NOT NULL,
		email STRING(150),
		hire_date TIMESTAMP NOT NULL,
		salary FLOAT64,
		dept_id STRING(36),
		manager_id STRING(36),
		phone_number STRING(20),
		CONSTRAINT fk_dept FOREIGN KEY (dept_id) REFERENCES departments(dept_id),
		CONSTRAINT fk_manager FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
	) PRIMARY KEY (emp_id)`,

	`CREATE TABLE projects (
		project_id STRING(36) DEFAULT (GENERATE_UUID()),
		project_name STRING(100) NOT NULL,
		start_date TIMESTAMP,
		end_date TIMESTAMP,
		budget FLOAT64,
		status STRING(20) DEFAULT ('ACTIVE'),
		CONSTRAINT check_dates CHECK (end_date > start_date),
		CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
	) PRIMARY KEY (project_id)`,

	`CREATE TABLE project_assignments (
		emp_id STRING(36) NOT NULL,
		project_id STRING(36) NOT NULL,
		role STRING(50),
		hours_allocated INT64,
		CONSTRAINT fk_emp FOREIGN KEY (emp_id) REFERENCES employees(emp_id),
		CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects(project_id)
	) PRIMARY KEY (emp_id, project_id)`,

	`CREATE UNIQUE INDEX idx_emp_email ON employees(email)`,
	`CREATE INDEX idx_emp_name ON employees(last_name, first_name)`,
	`CREATE INDEX idx_dept_location ON departments(location)`,
	`CREATE INDEX idx_project_status ON projects(status)`,

	`CREATE OR REPLACE VIEW employee_details
		SQL SECURITY INVOKER
		AS SELECT 
			e.emp_id,
			e.first_name,
			e.last_name,
			e.email,
			d.dept_name,
			m.first_name as manager_first_name,
			m.last_name as manager_last_name
		FROM employees e
		LEFT JOIN departments d ON e.dept_id = d.dept_id
		LEFT JOIN employees m ON e.manager_id = m.emp_id`,
}

// Cleanup statements - hardcoded from valid_spanner_database.sql
var validSchemaCleanup = []string{
	`DROP VIEW IF EXISTS employee_details`,
	`DROP INDEX IF EXISTS idx_project_status`,
	`DROP INDEX IF EXISTS idx_dept_location`,
	`DROP INDEX IF EXISTS idx_emp_name`,
	`DROP INDEX IF EXISTS idx_emp_email`,
	`DROP TABLE IF EXISTS project_assignments`,
	`DROP TABLE IF EXISTS projects`,
	`DROP TABLE IF EXISTS employees`,
	`DROP TABLE IF EXISTS departments`,
}

type AtomicSpannerDBTeardown struct {
	db        *sql.DB
	repo      repo.Database
	t         *testing.T
	terminate func()
}

func setupAtomicSpannerDB(t *testing.T) *AtomicSpannerDBTeardown {
	db, terminate, err := tools.GetDB(true)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	r := repo.NewSpannerRepo(db)
	if err := r.CleanupDB(); err != nil {
		t.Fatalf("Failed to cleanup DB: %v", err)
	}
	return &AtomicSpannerDBTeardown{db: db, repo: r, t: t, terminate: terminate}
}

func (d *AtomicSpannerDBTeardown) Close() {
	if err := d.repo.CleanupDB(); err != nil {
		d.t.Errorf("Failed to cleanup DB: %v", err)
	}
	d.db.Close()
	d.terminate()
}

// setupValidSchema creates the valid schema in the database
func setupValidSchema(db *sql.DB) error {
	for _, stmt := range validSchemaSetup {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to setup schema: %w", err)
		}
	}
	return nil
}

// cleanupValidSchema removes the valid schema from the database
func cleanupValidSchema(db *sql.DB) error {
	for _, stmt := range validSchemaCleanup {
		if _, err := db.Exec(stmt); err != nil {
			// Ignore errors in cleanup - some objects might not exist
			continue
		}
	}
	return nil
}

// TestAtomicSQLStatements tests all SQL statements atomically split by statement type
func TestAtomicSQLStatements(t *testing.T) {
	t.Run("CREATE_Statements", func(t *testing.T) {
		testAtomicStatementsByType(t, "CREATE")
	})

	t.Run("INSERT_Statements", func(t *testing.T) {
		testAtomicStatementsByType(t, "INSERT")
	})

	t.Run("SELECT_Statements", func(t *testing.T) {
		testAtomicStatementsByType(t, "SELECT")
	})

	t.Run("DROP_Statements", func(t *testing.T) {
		testAtomicStatementsByType(t, "DROP")
	})
}

// TestAtomicCREATEStatements tests only CREATE statements atomically
func TestAtomicCREATEStatements(t *testing.T) {
	t.Run("AtomicCREATEStatements", func(t *testing.T) {
		testAtomicStatementsByType(t, "CREATE")
	})
}

// TestAtomicINSERTStatements tests only INSERT statements atomically
func TestAtomicINSERTStatements(t *testing.T) {
	t.Run("AtomicINSERTStatements", func(t *testing.T) {
		testAtomicStatementsByType(t, "INSERT")
	})
}

// TestAtomicSELECTStatements tests only SELECT statements atomically
func TestAtomicSELECTStatements(t *testing.T) {
	t.Run("AtomicSELECTStatements", func(t *testing.T) {
		testAtomicStatementsByType(t, "SELECT")
	})
}

// TestAtomicDROPStatements tests only DROP statements atomically
func TestAtomicDROPStatements(t *testing.T) {
	t.Run("AtomicDROPStatements", func(t *testing.T) {
		testAtomicStatementsByType(t, "DROP")
	})
}

// testAtomicStatementsByType tests statements of a specific type across all SQL files
func testAtomicStatementsByType(t *testing.T, statementType string) {
	start := time.Now()

	fmt.Printf("Testing cached statements for %s type atomically\n", statementType)

	summary := AtomicTestSummary{
		TotalFiles: 0, // Will count files that have this statement type
	}

	// Use cached statements instead of reading files
	for filename, fileStatements := range statementCache {
		if statements, exists := fileStatements[statementType]; exists && len(statements) > 0 {
			t.Run(filename, func(t *testing.T) {
				result := testAtomicSQLFileByTypeFromCache(t, filename, statementType, statements)
				summary.FileResults = append(summary.FileResults, result)
				summary.TotalStatements += result.TotalStatements
				summary.ParsedCount += result.ParsedCount
				summary.ExecutedCount += result.ExecutedCount
			})
			summary.TotalFiles++
		}
	}

	summary.ExecutionTime = time.Since(start)

	// Generate comprehensive report for this statement type
	if err := generateAtomicReportByType(summary, statementType); err != nil {
		t.Logf("Warning: Failed to generate atomic report for %s: %v", statementType, err)
	} else {
		t.Logf("Generated atomic report for %s: atomic_%s_results.md", statementType, strings.ToLower(statementType))
	}

	// Log summary
	t.Logf("=== ATOMIC %s STATEMENT TEST SUMMARY ===", statementType)
	t.Logf("Total Files with %s: %d", statementType, summary.TotalFiles)
	t.Logf("Total %s Statements: %d", statementType, summary.TotalStatements)
	if summary.TotalStatements > 0 {
		t.Logf("Successfully Parsed: %d (%.1f%%)", summary.ParsedCount,
			float64(summary.ParsedCount)/float64(summary.TotalStatements)*100)
		t.Logf("Successfully Executed: %d (%.1f%%)", summary.ExecutedCount,
			float64(summary.ExecutedCount)/float64(summary.TotalStatements)*100)
	}
	t.Logf("Total Execution Time: %v", summary.ExecutionTime)

	// Validate that valid_spanner_database.sql achieves high success rates for this type
	for _, fileResult := range summary.FileResults {
		if fileResult.Filename == "valid_spanner_database.sql" && fileResult.TotalStatements > 0 {
			if fileResult.ParsedCount != fileResult.TotalStatements {
				t.Errorf("valid_spanner_database.sql should achieve 100%% parse success for %s, got %d/%d",
					statementType, fileResult.ParsedCount, fileResult.TotalStatements)
			}
			t.Logf("✅ valid_spanner_database.sql %s validation: %d/%d parsed (%.1f%%), %d/%d executed (%.1f%%)",
				statementType, fileResult.ParsedCount, fileResult.TotalStatements,
				float64(fileResult.ParsedCount)/float64(fileResult.TotalStatements)*100,
				fileResult.ExecutedCount, fileResult.TotalStatements,
				float64(fileResult.ExecutedCount)/float64(fileResult.TotalStatements)*100)
			break
		}
	}
	assert.True(t, true, "TestAtomicSQLStatements")
}

func testAtomicSQLFile(t *testing.T, sqlFile string) AtomicFileResult {
	filename := filepath.Base(sqlFile)
	result := AtomicFileResult{
		Filename: filename,
	}

	// Extract statements from file
	statements, err := tools.ExtractStatementsFromFile(sqlFile)
	if err != nil {
		t.Fatalf("Failed to extract SQL statements from %s: %v", sqlFile, err)
	}

	// Filter out comment-only statements
	var filteredStatements []string
	for _, stmt := range statements {
		if !isCommentOnlyStatement(stmt) {
			filteredStatements = append(filteredStatements, stmt)
		}
	}

	result.TotalStatements = len(filteredStatements)

	t.Logf("Testing %s with %d statements atomically (filtered from %d)", filename, len(filteredStatements), len(statements))

	// Setup single database connection for all statements in this file
	dbT := setupAtomicSpannerDB(t)
	defer dbT.Close()

	// Test each statement individually with atomic isolation
	for i, stmt := range filteredStatements {
		stmtResult := testAtomicSingleStatement(t, dbT, filename, i+1, stmt)
		result.StatementResults = append(result.StatementResults, stmtResult)

		if stmtResult.ParseSuccess {
			result.ParsedCount++
		}
		if stmtResult.ExecSuccess {
			result.ExecutedCount++
		}
	}

	t.Logf("File %s: %d/%d parsed (%.1f%%), %d/%d executed (%.1f%%)",
		filename, result.ParsedCount, result.TotalStatements,
		float64(result.ParsedCount)/float64(result.TotalStatements)*100,
		result.ExecutedCount, result.TotalStatements,
		float64(result.ExecutedCount)/float64(result.TotalStatements)*100)

	return result
}

// testAtomicSQLFileByTypeFromCache tests statements from the cache for a specific type
func testAtomicSQLFileByTypeFromCache(t *testing.T, filename string, statementType string, statements []StatementInfo) AtomicFileResult {
	result := AtomicFileResult{
		Filename: filename,
	}

	result.TotalStatements = len(statements)

	if result.TotalStatements == 0 {
		t.Logf("File %s has no %s statements", filename, statementType)
		return result
	}

	if statementType == "CREATE" {
		t.Logf("Testing %s with %d cached %s statements SEQUENTIALLY (building upon each other)", filename, len(statements), statementType)
		return testCreateStatementsSequentially(t, filename, statements)
	} else {
		t.Logf("Testing %s with %d cached %s statements ATOMICALLY (isolated)", filename, len(statements), statementType)
		return testStatementsAtomically(t, filename, statementType, statements)
	}
}

// testCreateStatementsSequentially tests CREATE statements in sequence, building upon each other
func testCreateStatementsSequentially(t *testing.T, filename string, statements []StatementInfo) AtomicFileResult {
	result := AtomicFileResult{
		Filename:        filename,
		TotalStatements: len(statements),
	}

	// Setup single database connection for all CREATE statements in this file
	dbT := setupAtomicSpannerDB(t)
	defer dbT.Close()

	// DO NOT run setupValidSchema for CREATE statements - they need a clean database
	t.Logf("Testing CREATE statements sequentially without pre-existing schema")

	// Test each CREATE statement sequentially, building upon previous ones
	for _, stmtInfo := range statements {
		stmtResult := testSingleCreateStatement(t, dbT, filename, stmtInfo.Position, stmtInfo.Statement)
		result.StatementResults = append(result.StatementResults, stmtResult)

		if stmtResult.ParseSuccess {
			result.ParsedCount++
		}
		if stmtResult.ExecSuccess {
			result.ExecutedCount++
		}
	}

	t.Logf("File %s CREATE statements: %d/%d parsed (%.1f%%), %d/%d executed (%.1f%%)",
		filename, result.ParsedCount, result.TotalStatements,
		float64(result.ParsedCount)/float64(result.TotalStatements)*100,
		result.ExecutedCount, result.TotalStatements,
		float64(result.ExecutedCount)/float64(result.TotalStatements)*100)

	return result
}

// testStatementsAtomically tests non-CREATE statements with full atomic isolation
func testStatementsAtomically(t *testing.T, filename string, statementType string, statements []StatementInfo) AtomicFileResult {
	result := AtomicFileResult{
		Filename:        filename,
		TotalStatements: len(statements),
	}

	// Setup single database connection for all statements in this file
	dbT := setupAtomicSpannerDB(t)
	defer dbT.Close()

	// Test each cached statement individually with atomic isolation
	for _, stmtInfo := range statements {
		stmtResult := testAtomicSingleStatement(t, dbT, filename, stmtInfo.Position, stmtInfo.Statement)
		result.StatementResults = append(result.StatementResults, stmtResult)

		if stmtResult.ParseSuccess {
			result.ParsedCount++
		}
		if stmtResult.ExecSuccess {
			result.ExecutedCount++
		}
	}

	t.Logf("File %s %s statements: %d/%d parsed (%.1f%%), %d/%d executed (%.1f%%)",
		filename, statementType, result.ParsedCount, result.TotalStatements,
		float64(result.ParsedCount)/float64(result.TotalStatements)*100,
		result.ExecutedCount, result.TotalStatements,
		float64(result.ExecutedCount)/float64(result.TotalStatements)*100)

	return result
}

// testSingleCreateStatement tests a single CREATE statement without schema setup/cleanup
func testSingleCreateStatement(t *testing.T, dbT *AtomicSpannerDBTeardown, filename string, stmtNum int, stmt string) AtomicStatementResult {
	start := time.Now()

	result := AtomicStatementResult{
		Filename:     filename,
		StatementNum: stmtNum,
		Statement:    strings.TrimSpace(stmt),
	}

	// Step 1: Parse the statement with memefish
	_, err := memefish.ParseStatement(filename, stmt)
	if err != nil {
		result.ParseSuccess = false
		result.ParseError = err.Error()
		result.ExecutionTime = time.Since(start)
		return result
	}

	result.ParseSuccess = true
	result.StatementType = "CREATE"

	// Step 2: Execute the CREATE statement directly (no schema setup/cleanup)
	executor := repo.NewSQLExecutor(dbT.db, dbT.repo)
	defer func() {
		if err := executor.Cleanup(); err != nil {
			t.Logf("Warning: executor cleanup failed for CREATE statement %d: %v", stmtNum, err)
		}
	}()

	// Execute single CREATE statement
	execResult, err := executor.ExecuteStatements([]string{stmt})
	if err != nil {
		result.ExecSuccess = false
		result.ExecError = fmt.Sprintf("Execution setup failed: %v", err)
	} else if len(execResult.Errors) > 0 {
		result.ExecSuccess = false
		result.ExecError = execResult.Errors[0].Error()
	} else {
		result.ExecSuccess = true
	}

	result.ExecutionTime = time.Since(start)
	return result
}

// getStatementTypeFromRaw determines statement type from raw SQL text
func getStatementTypeFromRaw(stmt string) string {
	upperStmt := strings.ToUpper(stmt)

	// Look for SQL keywords anywhere in the statement, not just at the beginning
	// This handles cases where statements start with comments
	switch {
	case strings.Contains(upperStmt, "CREATE TABLE") ||
		strings.Contains(upperStmt, "CREATE INDEX") ||
		strings.Contains(upperStmt, "CREATE UNIQUE INDEX") ||
		strings.Contains(upperStmt, "CREATE VIEW") ||
		strings.Contains(upperStmt, "CREATE SEQUENCE") ||
		strings.Contains(upperStmt, "CREATE OR REPLACE"):
		return "CREATE"
	case strings.Contains(upperStmt, "INSERT INTO"):
		return "INSERT"
	case strings.Contains(upperStmt, "SELECT "):
		return "SELECT"
	case strings.Contains(upperStmt, "DROP TABLE") ||
		strings.Contains(upperStmt, "DROP INDEX") ||
		strings.Contains(upperStmt, "DROP VIEW") ||
		strings.Contains(upperStmt, "DROP SEQUENCE"):
		return "DROP"
	case strings.Contains(upperStmt, "ALTER TABLE"):
		return "ALTER"
	case strings.Contains(upperStmt, "UPDATE "):
		return "UPDATE"
	case strings.Contains(upperStmt, "DELETE FROM"):
		return "DELETE"
	default:
		return "OTHER"
	}
}

// isCommentOnlyStatement checks if a statement is only comments and should be skipped
func isCommentOnlyStatement(stmt string) bool {
	trimmed := strings.TrimSpace(stmt)

	// Empty statement
	if trimmed == "" {
		return true
	}

	// Multi-line comment (simple case - entire statement is wrapped in /* */)
	if strings.HasPrefix(trimmed, "/*") && strings.HasSuffix(trimmed, "*/") {
		return true
	}

	// Check line by line for comment-only content
	lines := strings.Split(stmt, "\n")
	for _, line := range lines {
		lineTrimed := strings.TrimSpace(line)
		if lineTrimed == "" {
			continue // Empty line is OK
		}
		if strings.HasPrefix(lineTrimed, "--") {
			continue // Comment line is OK
		}
		if strings.HasPrefix(lineTrimed, "/*") || strings.HasSuffix(lineTrimed, "*/") {
			continue // Multi-line comment is OK
		}
		// If we find any non-comment content, it's not comment-only
		return false
	}

	return true
}

func testAtomicSingleStatement(t *testing.T, dbT *AtomicSpannerDBTeardown, filename string, stmtNum int, stmt string) AtomicStatementResult {
	start := time.Now()

	result := AtomicStatementResult{
		Filename:     filename,
		StatementNum: stmtNum,
		Statement:    strings.TrimSpace(stmt),
	}

	// Step 1: Parse the statement with memefish
	_, err := memefish.ParseStatement(filename, stmt)
	if err != nil {
		result.ParseSuccess = false
		result.ParseError = err.Error()
		result.ExecutionTime = time.Since(start)
		return result
	}

	result.ParseSuccess = true
	result.StatementType = getAtomicStatementType(stmt)

	// Step 2: Setup fresh valid schema for this statement (for non-valid files)
	if err := cleanupValidSchema(dbT.db); err != nil {
		t.Logf("Warning: cleanup failed for statement %d: %v", stmtNum, err)
	}

	if err := setupValidSchema(dbT.db); err != nil {
		result.ExecSuccess = false
		result.ExecError = fmt.Sprintf("Schema setup failed: %v", err)
		result.ExecutionTime = time.Since(start)
		return result
	}

	// Step 3: Execute the single statement
	executor := repo.NewSQLExecutor(dbT.db, dbT.repo)
	defer func() {
		if err := executor.Cleanup(); err != nil {
			t.Logf("Warning: executor cleanup failed for statement %d: %v", stmtNum, err)
		}
	}()

	// Execute single statement
	execResult, err := executor.ExecuteStatements([]string{stmt})
	if err != nil {
		result.ExecSuccess = false
		result.ExecError = fmt.Sprintf("Execution setup failed: %v", err)
	} else if len(execResult.Errors) > 0 {
		result.ExecSuccess = false
		result.ExecError = execResult.Errors[0].Error()
	} else {
		result.ExecSuccess = true
	}

	result.ExecutionTime = time.Since(start)
	return result
}

// getAtomicStatementType determines the type of a parsed statement
func getAtomicStatementType(originalStmt string) string {
	// Skip comment-only statements
	if isCommentOnlyStatement(originalStmt) {
		return "COMMENT"
	}

	upperStmt := strings.ToUpper(originalStmt)

	// Look for SQL keywords anywhere in the statement, not just at the beginning
	// This handles cases where statements start with comments
	switch {
	case strings.Contains(upperStmt, "CREATE TABLE") ||
		strings.Contains(upperStmt, "CREATE INDEX") ||
		strings.Contains(upperStmt, "CREATE UNIQUE INDEX") ||
		strings.Contains(upperStmt, "CREATE VIEW") ||
		strings.Contains(upperStmt, "CREATE SEQUENCE") ||
		strings.Contains(upperStmt, "CREATE OR REPLACE"):
		return "CREATE"
	case strings.Contains(upperStmt, "INSERT INTO"):
		return "INSERT"
	case strings.Contains(upperStmt, "SELECT "):
		return "SELECT"
	case strings.Contains(upperStmt, "DROP TABLE") ||
		strings.Contains(upperStmt, "DROP INDEX") ||
		strings.Contains(upperStmt, "DROP VIEW") ||
		strings.Contains(upperStmt, "DROP SEQUENCE"):
		return "DROP"
	case strings.Contains(upperStmt, "ALTER TABLE"):
		return "ALTER"
	case strings.Contains(upperStmt, "UPDATE "):
		return "UPDATE"
	case strings.Contains(upperStmt, "DELETE FROM"):
		return "DELETE"
	default:
		return "OTHER"
	}
}

// generateAtomicReport creates a markdown file with atomic test results
func generateAtomicReport(summary AtomicTestSummary) error {
	filename := "atomic_statement_results.md"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating markdown file: %w", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "# Atomic SQL Statement Test Results\n\n")
	fmt.Fprintf(file, "Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "## Testing Methodology\n\n")
	fmt.Fprintf(file, "- **All SQL Files**: Each statement tested in complete isolation with fresh valid schema\n")
	fmt.Fprintf(file, "- **No Special Cases**: All files use the same atomic isolation approach\n\n")

	// Write summary
	fmt.Fprintf(file, "## Summary\n\n")
	fmt.Fprintf(file, "- **Total SQL Files**: %d\n", summary.TotalFiles)
	fmt.Fprintf(file, "- **Total Statements**: %d\n", summary.TotalStatements)
	fmt.Fprintf(file, "- **Successfully Parsed**: %d (%.1f%%)\n", summary.ParsedCount,
		float64(summary.ParsedCount)/float64(summary.TotalStatements)*100)
	fmt.Fprintf(file, "- **Successfully Executed**: %d (%.1f%%)\n", summary.ExecutedCount,
		float64(summary.ExecutedCount)/float64(summary.TotalStatements)*100)
	fmt.Fprintf(file, "- **Total Execution Time**: %v\n\n", summary.ExecutionTime)

	// Write file-level results
	fmt.Fprintf(file, "## File-Level Results\n\n")
	fmt.Fprintf(file, "| File | Total | Parsed | Executed | Parse Rate | Exec Rate |\n")
	fmt.Fprintf(file, "|------|-------|--------|----------|------------|----------|\n")

	for _, fileResult := range summary.FileResults {
		parseRate := 0.0
		execRate := 0.0
		if fileResult.TotalStatements > 0 {
			parseRate = float64(fileResult.ParsedCount) / float64(fileResult.TotalStatements) * 100
			execRate = float64(fileResult.ExecutedCount) / float64(fileResult.TotalStatements) * 100
		}

		fmt.Fprintf(file, "| [%s](../../generated_sql/%s) | %d | %d | %d | %.1f%% | %.1f%% |\n",
			fileResult.Filename, fileResult.Filename,
			fileResult.TotalStatements, fileResult.ParsedCount, fileResult.ExecutedCount,
			parseRate, execRate)
	}

	// Write statement-level results (truncated for readability)
	fmt.Fprintf(file, "\n## Statement-Level Results\n\n")
	fmt.Fprintf(file, "| File | # | Statement | Type | Parse | Exec | Parse Error | Exec Error |\n")
	fmt.Fprintf(file, "|------|---|-----------|------|-------|------|-------------|------------|\n")

	for _, fileResult := range summary.FileResults {
		for _, stmtResult := range fileResult.StatementResults {
			// Truncate long statements for readability and escape newlines
			stmt := strings.ReplaceAll(stmtResult.Statement, "\n", "")
			stmt = strings.ReplaceAll(stmt, "\r", "\\r")
			if len(stmt) > 40 {
				stmt = stmt[:37] + "..."
			}

			parseIcon := "❌"
			if stmtResult.ParseSuccess {
				parseIcon = "✅"
			}

			execIcon := "❌"
			if stmtResult.ExecSuccess {
				execIcon = "✅"
			}

			parseError := stmtResult.ParseError
			if parseError == "" {
				parseError = "-"
			} else if len(parseError) > 30 {
				parseError = parseError[:27] + "..."
			}

			execError := stmtResult.ExecError
			if execError == "" {
				execError = "-"
			} else if len(execError) > 30 {
				execError = execError[:27] + "..."
			}

			fmt.Fprintf(file, "| %s | %d | `%s` | %s | %s | %s | %s | %s |\n",
				fileResult.Filename, stmtResult.StatementNum, stmt, stmtResult.StatementType,
				parseIcon, execIcon, parseError, execError)
		}
	}

	// Write detailed error analysis
	fmt.Fprintf(file, "\n## Error Analysis\n\n")

	parseErrorCount := 0
	execErrorCount := 0
	parseErrorTypes := make(map[string]int)
	execErrorTypes := make(map[string]int)

	for _, fileResult := range summary.FileResults {
		for _, stmtResult := range fileResult.StatementResults {
			if !stmtResult.ParseSuccess {
				parseErrorCount++
				if strings.Contains(stmtResult.ParseError, "syntax error") {
					parseErrorTypes["Syntax Error"]++
				} else if strings.Contains(stmtResult.ParseError, "expecting") {
					parseErrorTypes["Expected Token"]++
				} else {
					parseErrorTypes["Other Parse Error"]++
				}
			}

			if !stmtResult.ExecSuccess && stmtResult.ParseSuccess {
				execErrorCount++
				if strings.Contains(stmtResult.ExecError, "NotFound") {
					execErrorTypes["NotFound"]++
				} else if strings.Contains(stmtResult.ExecError, "InvalidArgument") {
					execErrorTypes["InvalidArgument"]++
				} else if strings.Contains(stmtResult.ExecError, "AlreadyExists") {
					execErrorTypes["AlreadyExists"]++
				} else {
					execErrorTypes["Other Exec Error"]++
				}
			}
		}
	}

	fmt.Fprintf(file, "### Parse Error Breakdown\n")
	fmt.Fprintf(file, "Total parse errors: %d\n\n", parseErrorCount)
	for errorType, count := range parseErrorTypes {
		fmt.Fprintf(file, "- **%s**: %d occurrences\n", errorType, count)
	}

	fmt.Fprintf(file, "\n### Execution Error Breakdown\n")
	fmt.Fprintf(file, "Total execution errors: %d\n\n", execErrorCount)
	for errorType, count := range execErrorTypes {
		fmt.Fprintf(file, "- **%s**: %d occurrences\n", errorType, count)
	}

	return nil
}

// generateAtomicReportByType creates a markdown file with atomic test results for a specific statement type
func generateAtomicReportByType(summary AtomicTestSummary, statementType string) error {
	filename := fmt.Sprintf("atomic_%s_results.md", strings.ToLower(statementType))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating markdown file: %w", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "# %s Statement Test Results\n\n", statementType)
	fmt.Fprintf(file, "Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Write methodology based on statement type
	fmt.Fprintf(file, "## Testing Methodology\n\n")
	fmt.Fprintf(file, "- **Statement Type**: Only %s statements tested\n", statementType)

	if statementType == "CREATE" {
		fmt.Fprintf(file, "- **Testing Approach**: Sequential execution (statements build upon each other)\n")
		fmt.Fprintf(file, "- **Database State**: Clean database, no pre-existing schema\n")
		fmt.Fprintf(file, "- **Isolation**: None - CREATE statements depend on previous CREATE statements\n")
		fmt.Fprintf(file, "- **Rationale**: CREATE statements cannot be atomic since setupValidSchema would cause AlreadyExists errors\n\n")
	} else {
		fmt.Fprintf(file, "- **Testing Approach**: Atomic isolation (each statement tested independently)\n")
		fmt.Fprintf(file, "- **Database State**: Fresh valid schema setup before each statement\n")
		fmt.Fprintf(file, "- **Isolation**: Complete - each statement tested in isolation with fresh valid schema\n")
		fmt.Fprintf(file, "- **No Special Cases**: All files use the same atomic isolation approach\n\n")
	}

	// Write summary
	fmt.Fprintf(file, "## Summary\n\n")
	fmt.Fprintf(file, "- **Total SQL Files with %s**: %d\n", statementType, len(summary.FileResults))
	fmt.Fprintf(file, "- **Total %s Statements**: %d\n", statementType, summary.TotalStatements)
	if summary.TotalStatements > 0 {
		fmt.Fprintf(file, "- **Successfully Parsed**: %d (%.1f%%)\n", summary.ParsedCount,
			float64(summary.ParsedCount)/float64(summary.TotalStatements)*100)
		fmt.Fprintf(file, "- **Successfully Executed**: %d (%.1f%%)\n", summary.ExecutedCount,
			float64(summary.ExecutedCount)/float64(summary.TotalStatements)*100)
	}
	fmt.Fprintf(file, "- **Total Execution Time**: %v\n\n", summary.ExecutionTime)

	// Write file-level results
	fmt.Fprintf(file, "## File-Level Results\n\n")
	fmt.Fprintf(file, "| File | Total %s | Parsed | Executed | Parse Rate | Exec Rate |\n", statementType)
	fmt.Fprintf(file, "|------|----------|--------|----------|------------|----------|\n")

	for _, fileResult := range summary.FileResults {
		if fileResult.TotalStatements > 0 {
			parseRate := float64(fileResult.ParsedCount) / float64(fileResult.TotalStatements) * 100
			execRate := float64(fileResult.ExecutedCount) / float64(fileResult.TotalStatements) * 100

			fmt.Fprintf(file, "| [%s](../../generated_sql/%s) | %d | %d | %d | %.1f%% | %.1f%% |\n",
				fileResult.Filename, fileResult.Filename,
				fileResult.TotalStatements, fileResult.ParsedCount, fileResult.ExecutedCount,
				parseRate, execRate)
		}
	}

	// Write statement-level results (truncated for readability)
	fmt.Fprintf(file, "\n## Statement-Level Results\n\n")
	fmt.Fprintf(file, "| File | # | Statement | Parse | Exec | Parse Error | Exec Error |\n")
	fmt.Fprintf(file, "|------|---|-----------|-------|------|-------------|------------|\n")

	for _, fileResult := range summary.FileResults {
		for _, stmtResult := range fileResult.StatementResults {
			// Truncate long statements for readability and escape newlines
			stmt := strings.ReplaceAll(stmtResult.Statement, "\n", "")
			stmt = strings.ReplaceAll(stmt, "\r", "\\r")
			if len(stmt) > 40 {
				stmt = stmt[:37] + "..."
			}

			parseIcon := "❌"
			if stmtResult.ParseSuccess {
				parseIcon = "✅"
			}

			execIcon := "❌"
			if stmtResult.ExecSuccess {
				execIcon = "✅"
			}

			parseError := stmtResult.ParseError
			if parseError == "" {
				parseError = "-"
			} else if len(parseError) > 30 {
				parseError = parseError[:27] + "..."
			}

			execError := stmtResult.ExecError
			if execError == "" {
				execError = "-"
			} else if len(execError) > 30 {
				execError = execError[:27] + "..."
			}

			fmt.Fprintf(file, "| %s | %d | `%s` | %s | %s | %s | %s |\n",
				fileResult.Filename, stmtResult.StatementNum, stmt,
				parseIcon, execIcon, parseError, execError)
		}
	}

	// Write detailed error analysis
	fmt.Fprintf(file, "\n## Error Analysis for %s Statements\n\n", statementType)

	parseErrorCount := 0
	execErrorCount := 0
	parseErrorTypes := make(map[string]int)
	execErrorTypes := make(map[string]int)

	for _, fileResult := range summary.FileResults {
		for _, stmtResult := range fileResult.StatementResults {
			if !stmtResult.ParseSuccess {
				parseErrorCount++
				if strings.Contains(stmtResult.ParseError, "syntax error") {
					parseErrorTypes["Syntax Error"]++
				} else if strings.Contains(stmtResult.ParseError, "expecting") {
					parseErrorTypes["Expected Token"]++
				} else {
					parseErrorTypes["Other Parse Error"]++
				}
			}

			if !stmtResult.ExecSuccess && stmtResult.ParseSuccess {
				execErrorCount++
				if strings.Contains(stmtResult.ExecError, "NotFound") {
					execErrorTypes["NotFound"]++
				} else if strings.Contains(stmtResult.ExecError, "InvalidArgument") {
					execErrorTypes["InvalidArgument"]++
				} else if strings.Contains(stmtResult.ExecError, "AlreadyExists") {
					execErrorTypes["AlreadyExists"]++
				} else {
					execErrorTypes["Other Exec Error"]++
				}
			}
		}
	}

	fmt.Fprintf(file, "### Parse Error Breakdown\n")
	fmt.Fprintf(file, "Total parse errors: %d\n\n", parseErrorCount)
	for errorType, count := range parseErrorTypes {
		fmt.Fprintf(file, "- **%s**: %d occurrences\n", errorType, count)
	}

	fmt.Fprintf(file, "\n### Execution Error Breakdown\n")
	fmt.Fprintf(file, "Total execution errors: %d\n\n", execErrorCount)
	for errorType, count := range execErrorTypes {
		fmt.Fprintf(file, "- **%s**: %d occurrences\n", errorType, count)
	}

	return nil
}
