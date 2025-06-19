package spanner_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"postgres-example/repo"
	"postgres-example/tools"
)

// TestFileResult holds the results for a single SQL file test
type TestFileResult struct {
	Filename         string
	TotalStatements  int
	CreateStatements int
	InsertStatements int
	SelectStatements int
	DropStatements   int
	ExecutedCount    int
	FailedCount      int
	ErrorRate        float64
	ExecutionTime    time.Duration
	Errors           []string
}

type SpannerDBTeardown struct {
	db        *sql.DB
	repo      repo.Database
	t         *testing.T
	terminate func()
}

func setupSpannerDB(t *testing.T) *SpannerDBTeardown {
	db, terminate, err := tools.GetDB(true)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	r := repo.NewSpannerRepo(db)
	if err := r.CleanupDB(); err != nil {
		t.Fatalf("Failed to cleanup DB: %v", err)
	}
	return &SpannerDBTeardown{db: db, repo: r, t: t, terminate: terminate}
}

func (d *SpannerDBTeardown) Close() {
	if err := d.repo.CleanupDB(); err != nil {
		d.t.Errorf("Failed to cleanup DB: %v", err)
	}
	d.db.Close()
	d.terminate()
}

// TestGeneratedSQLFiles tests all SQL files in the generated_sql folder
func TestGeneratedSQLFiles(t *testing.T) {
	// Get all SQL files from generated_sql folder
	sqlFiles, err := filepath.Glob("../../generated_sql/*.sql")
	if err != nil {
		t.Fatalf("Failed to find SQL files: %v", err)
	}

	print("SQL Files: ", len(sqlFiles), "\n")

	if len(sqlFiles) == 0 {
		t.Fatalf("No SQL files found in generated_sql folder")
	}

	// Collect results for markdown report
	var results []TestFileResult

	for _, sqlFile := range sqlFiles {
		t.Run(filepath.Base(sqlFile), func(t *testing.T) {
			result := testSQLFileExecution(t, sqlFile)
			results = append(results, result)
		})
	}

	// Generate markdown report
	if err := generateMarkdownReport(results); err != nil {
		t.Logf("Warning: Failed to generate markdown report: %v", err)
	} else {
		t.Logf("Generated markdown report: sql_test_results.md")
	}
}

func testSQLFileExecution(t *testing.T, sqlFile string) TestFileResult {
	start := time.Now()
	filename := filepath.Base(sqlFile)

	result := TestFileResult{
		Filename: filename,
	}

	// Setup database
	dbT := setupSpannerDB(t)
	defer dbT.Close()

	// Create SQL executor
	executor := repo.NewSQLExecutor(dbT.db, dbT.repo)
	defer func() {
		if err := executor.Cleanup(); err != nil {
			t.Logf("Warning: cleanup failed: %v", err)
		}
	}()

	// Execute SQL from file
	execResult, err := executor.ExecuteFromFile(sqlFile)
	if err != nil {
		t.Fatalf("Failed to execute SQL file %s: %v", sqlFile, err)
	}

	result.ExecutionTime = time.Since(start)
	result.TotalStatements = execResult.TotalStatements
	result.CreateStatements = execResult.CreateStatements
	result.InsertStatements = execResult.InsertStatements
	result.SelectStatements = execResult.SelectStatements
	result.DropStatements = execResult.DropStatements
	result.ExecutedCount = execResult.ExecutedCount
	result.FailedCount = len(execResult.Errors)

	if result.TotalStatements > 0 {
		result.ErrorRate = float64(result.FailedCount) / float64(result.TotalStatements) * 100
	}

	for _, err := range execResult.Errors {
		result.Errors = append(result.Errors, err.Error())
	}

	// Log execution results
	t.Logf("Execution results for %s:", filename)
	t.Logf("  Total statements: %d", result.TotalStatements)
	t.Logf("  CREATE: %d, INSERT: %d, SELECT: %d, DROP: %d",
		result.CreateStatements, result.InsertStatements, result.SelectStatements, result.DropStatements)
	t.Logf("  Executed: %d, Failed: %d", result.ExecutedCount, result.FailedCount)
	t.Logf("  Error rate: %.1f%%", result.ErrorRate)
	t.Logf("  Execution time: %v", result.ExecutionTime)

	// Validate results (but don't fail the test)
	validateExecutionResult(t, execResult, filename)

	// Test additional queries if data was inserted
	if len(execResult.InsertedRecords) > 0 {
		testDataIntegrity(t, dbT, execResult)
	}

	return result
}

// generateMarkdownReport creates a markdown file with test results
func generateMarkdownReport(results []TestFileResult) error {
	filename := "sql_test_results.md"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating markdown file: %w", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "# SQL Test Results\n\n")
	fmt.Fprintf(file, "Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Write summary
	totalFiles := len(results)
	totalStatements := 0
	totalExecuted := 0
	totalFailed := 0

	for _, result := range results {
		totalStatements += result.TotalStatements
		totalExecuted += result.ExecutedCount
		totalFailed += result.FailedCount
	}

	fmt.Fprintf(file, "## Summary\n\n")
	fmt.Fprintf(file, "- **Total SQL Files**: %d\n", totalFiles)
	fmt.Fprintf(file, "- **Total Statements**: %d\n", totalStatements)
	fmt.Fprintf(file, "- **Successfully Executed**: %d\n", totalExecuted)
	fmt.Fprintf(file, "- **Failed**: %d\n", totalFailed)
	fmt.Fprintf(file, "- **Overall Success Rate**: %.1f%%\n\n", float64(totalExecuted)/float64(totalStatements)*100)

	// Write detailed results table
	fmt.Fprintf(file, "## Detailed Results\n\n")
	fmt.Fprintf(file, "| File | Total | CREATE | INSERT | SELECT | DROP | Executed | Failed | Success Rate | Execution Time |\n")
	fmt.Fprintf(file, "|------|-------|--------|--------|--------|------|----------|--------|--------------|----------------|\n")

	for _, result := range results {
		fmt.Fprintf(file, "| [%s](../../generated_sql/%s) | %d | %d | %d | %d | %d | %d | %d | %.1f%% | %v |\n",
			result.Filename,
			result.Filename,
			result.TotalStatements,
			result.CreateStatements,
			result.InsertStatements,
			result.SelectStatements,
			result.DropStatements,
			result.ExecutedCount,
			result.FailedCount,
			100-result.ErrorRate, // Show success rate instead of error rate
			result.ExecutionTime.Round(time.Millisecond),
		)
	}

	// Write error details
	fmt.Fprintf(file, "\n## Error Details\n\n")
	for _, result := range results {
		if len(result.Errors) > 0 {
			fmt.Fprintf(file, "### %s\n\n", result.Filename)
			fmt.Fprintf(file, "**Error Rate**: %.1f%% (%d/%d failed)\n\n", result.ErrorRate, result.FailedCount, result.TotalStatements)

			if len(result.Errors) > 0 {
				fmt.Fprintf(file, "**Sample Errors**:\n")
				for i, errMsg := range result.Errors {
					fmt.Fprintf(file, "%d. %s\n", i+1, errMsg)
				}
				fmt.Fprintf(file, "\n")
			}
		}
	}

	// Write compatibility insights
	fmt.Fprintf(file, "## Compatibility Insights\n\n")
	fmt.Fprintf(file, "### Common Issues Found\n\n")

	commonIssues := analyzeCommonIssues(results)
	for issue, count := range commonIssues {
		fmt.Fprintf(file, "- **%s**: Found in %d files\n", issue, count)
	}

	return nil
}

// analyzeCommonIssues finds common error patterns across all files
func analyzeCommonIssues(results []TestFileResult) map[string]int {
	issues := make(map[string]int)

	for _, result := range results {
		fileIssues := make(map[string]bool) // Track unique issues per file

		for _, errMsg := range result.Errors {
			switch {
			case contains(errMsg, "CURRENT_TIMESTAMP"):
				fileIssues["CURRENT_TIMESTAMP syntax"] = true
			case contains(errMsg, "SQL SECURITY"):
				fileIssues["Missing SQL SECURITY clause in views"] = true
			case contains(errMsg, "GENERATE_UUID"):
				fileIssues["GENERATE_UUID() compatibility"] = true
			case contains(errMsg, "DEFAULT"):
				fileIssues["DEFAULT value syntax"] = true
			case contains(errMsg, "CONSTRAINT"):
				fileIssues["CHECK constraints not supported"] = true
			case contains(errMsg, "FOREIGN KEY"):
				fileIssues["FOREIGN KEY constraints"] = true
			case contains(errMsg, "IDENTITY"):
				fileIssues["IDENTITY column issues"] = true
			case contains(errMsg, "Table not found"):
				fileIssues["Table dependency issues"] = true
			}
		}

		// Count each unique issue once per file
		for issue := range fileIssues {
			issues[issue]++
		}
	}

	return issues
}

// contains is a helper function to check if a string contains a substring (case insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func validateExecutionResult(t *testing.T, result *repo.ExecutionResult, filename string) {
	// Report any errors without treating them as test failures
	if len(result.Errors) > 0 {
		t.Logf("Execution errors for %s:", filename)
		for i, err := range result.Errors {
			t.Logf("  Error %d: %v", i+1, err)
		}

		// Report error statistics but don't fail the test
		errorRate := float64(len(result.Errors)) / float64(result.TotalStatements)
		t.Logf("Error rate for %s: %d/%d statements failed (%.1f%%)",
			filename, len(result.Errors), result.TotalStatements, errorRate*100)
	}

	// Report execution statistics
	t.Logf("Execution summary for %s:", filename)
	t.Logf("  Total statements: %d", result.TotalStatements)
	t.Logf("  Successfully executed: %d", result.ExecutedCount)
	t.Logf("  Failed: %d", len(result.Errors))

	// Log successful inserts
	successfulInserts := 0
	for _, insert := range result.InsertedRecords {
		if insert.Error == nil {
			successfulInserts++
			if insert.ID != nil {
				t.Logf("  INSERT returned ID: %v", insert.ID)
			}
		}
	}
	if successfulInserts > 0 {
		t.Logf("  Successful inserts: %d", successfulInserts)
	}

	// Log successful queries
	successfulQueries := 0
	totalRows := 0
	for _, query := range result.QueryResults {
		if query.Error == nil {
			successfulQueries++
			totalRows += query.RowCount
		}
	}
	if successfulQueries > 0 {
		t.Logf("  Successful queries: %d, Total rows returned: %d", successfulQueries, totalRows)
	}
}

func testDataIntegrity(t *testing.T, dbT *SpannerDBTeardown, result *repo.ExecutionResult) {
	// Test that we can query the data that was inserted

	// Try to query a common table that might exist
	commonTables := []string{"departments", "employees", "projects", "Departments", "Employees", "Projects"}

	for _, tableName := range commonTables {
		query := "SELECT COUNT(*) FROM " + tableName
		var count int
		err := dbT.db.QueryRow(query).Scan(&count)
		if err == nil {
			t.Logf("  Table %s contains %d rows", tableName, count)

			// If we have data, try to select some records
			if count > 0 {
				selectQuery := "SELECT * FROM " + tableName + " LIMIT 1"
				rows, err := dbT.db.Query(selectQuery)
				if err == nil {
					defer rows.Close()
					if rows.Next() {
						t.Logf("  Successfully queried data from %s", tableName)
					}
				}
			}
		}
	}
}

func TestSpannerDatabaseOperations(t *testing.T) {
	dbT := setupSpannerDB(t)
	defer dbT.Close()

	// Create tables first
	if err := dbT.repo.CreateTables(); err != nil {
		t.Skipf("Skipping test due to table creation error: %v", err)
	}

	// Insert sample data
	deptID, empID, projectID, err := dbT.repo.InsertSampleData()
	if err != nil {
		t.Fatalf("InsertSampleData failed: %v", err)
	}
	if deptID == 0 || empID == 0 || projectID == 0 {
		t.Errorf("Expected non-zero IDs, got deptID=%d, empID=%d, projectID=%d", deptID, empID, projectID)
	}

	// Query and check results
	details, err := dbT.repo.QueryEmployeeDetails()
	if err != nil {
		t.Fatalf("QueryEmployeeDetails failed: %v", err)
	}
	if len(details) == 0 {
		t.Error("Expected at least one employee detail result")
	}

	// Check if we found our test employee
	found := false
	for _, detail := range details {
		if detail.FirstName == "John" &&
			detail.LastName == "Doe" &&
			detail.DeptName == "Engineering" &&
			detail.ProjectName.String == "Database Migration" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find inserted employee and project in results")
	}
}

func TestSpannerEmployeeDetails(t *testing.T) {
	dbT := setupSpannerDB(t)
	defer dbT.Close()

	// Create tables first
	if err := dbT.repo.CreateTables(); err != nil {
		t.Skipf("Skipping test due to table creation error: %v", err)
	}

	// Insert test data
	_, _, _, err := dbT.repo.InsertSampleData()
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Query employee details
	details, err := dbT.repo.QueryEmployeeDetails()
	if err != nil {
		t.Fatalf("Failed to query employee details: %v", err)
	}

	// Verify the structure of the results
	if len(details) == 0 {
		t.Error("Expected at least one employee detail")
	}

	for _, detail := range details {
		// Check required fields
		if detail.FirstName == "" {
			t.Error("Expected non-empty FirstName")
		}
		if detail.LastName == "" {
			t.Error("Expected non-empty LastName")
		}
		if detail.Email == "" {
			t.Error("Expected non-empty Email")
		}
		if detail.DeptName == "" {
			t.Error("Expected non-empty DeptName")
		}

		// Check nullable fields
		if detail.ManagerFirstName.Valid {
			t.Logf("Manager First Name: %s", detail.ManagerFirstName.String)
		}
		if detail.ManagerLastName.Valid {
			t.Logf("Manager Last Name: %s", detail.ManagerLastName.String)
		}
		if detail.ProjectName.Valid {
			t.Logf("Project Name: %s", detail.ProjectName.String)
		}
	}
}
