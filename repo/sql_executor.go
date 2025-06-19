package repo

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"postgres-example/tools"
)

// SQLExecutor provides an abstraction for executing parsed SQL statements
type SQLExecutor struct {
	DB            *sql.DB
	repo          Database
	executed      map[string]bool
	createdTables []string
}

// NewSQLExecutor creates a new SQL executor instance
func NewSQLExecutor(db *sql.DB, repo Database) *SQLExecutor {
	return &SQLExecutor{
		DB:       db,
		repo:     repo,
		executed: make(map[string]bool),
	}
}

// ExecuteFromFile parses and executes SQL statements from a file
func (e *SQLExecutor) ExecuteFromFile(filename string) (*ExecutionResult, error) {
	parser := tools.NewSQLParser(filename)
	statements, err := parser.ParseSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to parse SQL file %s: %w", filename, err)
	}

	return e.ExecuteStatements(statements)
}

// ExecutionResult contains the results of executing SQL statements
type ExecutionResult struct {
	TotalStatements  int
	CreateStatements int
	InsertStatements int
	SelectStatements int
	DropStatements   int
	ExecutedCount    int
	SkippedCount     int
	Errors           []error
	InsertedRecords  []InsertResult
	QueryResults     []QueryResult
}

// InsertResult contains information about an insert operation
type InsertResult struct {
	Statement string
	ID        int64
	Error     error
}

// QueryResult contains information about a select operation
type QueryResult struct {
	Statement string
	RowCount  int
	Error     error
}

// ExecuteStatements executes a list of SQL statements in the proper order
func (e *SQLExecutor) ExecuteStatements(statements []string) (*ExecutionResult, error) {
	result := &ExecutionResult{
		TotalStatements: len(statements),
	}

	// Categorize statements
	var createStmts, insertStmts, selectStmts, dropStmts []string

	for _, stmt := range statements {
		upperStmt := strings.ToUpper(strings.TrimSpace(stmt))
		switch {
		case strings.HasPrefix(upperStmt, "CREATE"):
			createStmts = append(createStmts, stmt)
			result.CreateStatements++
		case strings.HasPrefix(upperStmt, "INSERT"):
			insertStmts = append(insertStmts, stmt)
			result.InsertStatements++
		case strings.HasPrefix(upperStmt, "SELECT"):
			selectStmts = append(selectStmts, stmt)
			result.SelectStatements++
		case strings.HasPrefix(upperStmt, "DROP"):
			dropStmts = append(dropStmts, stmt)
			result.DropStatements++
		}
	}

	// Execute in proper order: DROP, CREATE, INSERT, SELECT

	// 1. Execute DROP statements first (for cleanup)
	for _, stmt := range dropStmts {
		if err := e.executeDrop(stmt); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("DROP failed: %w", err))
		} else {
			result.ExecutedCount++
		}
	}

	// 2. Execute CREATE statements
	for _, stmt := range createStmts {
		if err := e.executeCreate(stmt); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("CREATE failed: %w", err))
		} else {
			result.ExecutedCount++
		}
	}

	// 3. Execute INSERT statements
	for _, stmt := range insertStmts {
		insertResult := e.executeInsert(stmt)
		result.InsertedRecords = append(result.InsertedRecords, insertResult)
		if insertResult.Error == nil {
			result.ExecutedCount++
		} else {
			result.Errors = append(result.Errors, fmt.Errorf("INSERT failed: %w", insertResult.Error))
		}
	}

	// 4. Execute SELECT statements
	for _, stmt := range selectStmts {
		queryResult := e.executeSelect(stmt)
		result.QueryResults = append(result.QueryResults, queryResult)
		if queryResult.Error == nil {
			result.ExecutedCount++
		} else {
			result.Errors = append(result.Errors, fmt.Errorf("SELECT failed: %w", queryResult.Error))
		}
	}

	result.SkippedCount = result.TotalStatements - result.ExecutedCount

	return result, nil
}

// executeDrop executes a DROP statement
func (e *SQLExecutor) executeDrop(stmt string) error {
	// Clean up the statement
	cleanStmt := e.cleanStatement(stmt)

	// Execute the drop statement
	_, err := e.DB.Exec(cleanStmt)
	if err != nil {
		// For DROP statements, be lenient about "not found" errors
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "not found") ||
			strings.Contains(errStr, "does not exist") ||
			strings.Contains(errStr, "unknown table") {
			// Log but don't fail for missing objects
			return nil
		}
		return fmt.Errorf("executing DROP statement: %w", err)
	}

	return nil
}

// executeCreate executes a CREATE statement
func (e *SQLExecutor) executeCreate(stmt string) error {
	// Clean up the statement
	cleanStmt := e.cleanStatement(stmt)

	// Track table names for cleanup
	if tableName := e.extractTableName(cleanStmt); tableName != "" {
		e.createdTables = append(e.createdTables, tableName)
	}

	// Execute the create statement
	_, err := e.DB.Exec(cleanStmt)
	if err != nil {
		return fmt.Errorf("executing CREATE statement: %w", err)
	}

	return nil
}

// executeInsert executes an INSERT statement
func (e *SQLExecutor) executeInsert(stmt string) InsertResult {
	result := InsertResult{Statement: stmt}

	// Clean up the statement
	cleanStmt := e.cleanStatement(stmt)

	// Handle parameterized queries by replacing with sample data
	populatedStmt, err := e.populateInsertParameters(cleanStmt)
	if err != nil {
		result.Error = fmt.Errorf("populating INSERT parameters: %w", err)
		return result
	}

	// For simplicity, always use regular insert (THEN RETURN syntax often causes issues)
	_, err = e.DB.Exec(populatedStmt)
	if err != nil {
		result.Error = fmt.Errorf("executing INSERT: %w", err)
		return result
	}

	// Try to get the last insert ID if possible (this is database-specific)
	// For Spanner, we'll just mark it as successful without trying to get the ID
	result.ID = 1 // Placeholder to indicate success

	return result
}

// executeSelect executes a SELECT statement
func (e *SQLExecutor) executeSelect(stmt string) QueryResult {
	result := QueryResult{Statement: stmt}

	// Clean up the statement
	cleanStmt := e.cleanStatement(stmt)

	// Execute the select and count rows
	rows, err := e.DB.Query(cleanStmt)
	if err != nil {
		result.Error = fmt.Errorf("executing SELECT: %w", err)
		return result
	}
	defer rows.Close()

	// Count the rows
	for rows.Next() {
		result.RowCount++
	}

	if err = rows.Err(); err != nil {
		result.Error = fmt.Errorf("reading SELECT results: %w", err)
		return result
	}

	return result
}

// cleanStatement removes extra whitespace and ensures proper formatting
func (e *SQLExecutor) cleanStatement(stmt string) string {
	// Remove leading/trailing whitespace
	cleaned := strings.TrimSpace(stmt)

	// Normalize internal whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned = re.ReplaceAllString(cleaned, " ")

	return cleaned
}

// extractTableName extracts the table name from a CREATE statement
func (e *SQLExecutor) extractTableName(stmt string) string {
	upperStmt := strings.ToUpper(stmt)

	// Pattern for CREATE TABLE statements
	if strings.Contains(upperStmt, "CREATE TABLE") {
		re := regexp.MustCompile(`CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?(\w+)`)
		matches := re.FindStringSubmatch(upperStmt)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

// populateInsertParameters replaces parameter placeholders with sample data
func (e *SQLExecutor) populateInsertParameters(stmt string) (string, error) {
	// Replace named parameters (@param) with positional parameters ($1, $2, etc.)
	namedParamRe := regexp.MustCompile(`@(\w+)`)
	paramIndex := 1
	paramMap := make(map[string]string)

	populated := namedParamRe.ReplaceAllStringFunc(stmt, func(match string) string {
		paramName := match[1:] // Remove @
		if replacement, exists := paramMap[paramName]; exists {
			return replacement
		}

		// Generate sample data based on parameter name
		value := e.getSampleValueForParameter(paramName)
		paramMap[paramName] = value
		paramIndex++
		return value
	})

	// Replace positional parameters ($1, $2, etc.) with sample values
	positionalParamRe := regexp.MustCompile(`\$\d+`)
	argIndex := 1
	populated = positionalParamRe.ReplaceAllStringFunc(populated, func(match string) string {
		value := e.getSampleValueForIndex(argIndex)
		argIndex++
		return value
	})

	return populated, nil
}

// getSampleValueForParameter returns sample data based on parameter name
func (e *SQLExecutor) getSampleValueForParameter(paramName string) string {
	lower := strings.ToLower(paramName)

	switch {
	case strings.Contains(lower, "name"):
		if strings.Contains(lower, "first") {
			return "'John'"
		} else if strings.Contains(lower, "last") {
			return "'Doe'"
		} else if strings.Contains(lower, "dept") {
			return "'Engineering'"
		} else if strings.Contains(lower, "project") {
			return "'Test Project'"
		}
		return "'Sample Name'"
	case strings.Contains(lower, "email"):
		return "'test@example.com'"
	case strings.Contains(lower, "location"):
		return "'New York'"
	case strings.Contains(lower, "date"):
		return "'2024-01-01'"
	case strings.Contains(lower, "salary"):
		return "75000"
	case strings.Contains(lower, "budget"):
		return "50000"
	case strings.Contains(lower, "hours"):
		return "40"
	case strings.Contains(lower, "id"):
		return "1"
	case strings.Contains(lower, "status"):
		return "'ACTIVE'"
	case strings.Contains(lower, "role"):
		return "'Developer'"
	case strings.Contains(lower, "phone"):
		return "'555-1234'"
	default:
		return "'sample_value'"
	}
}

// getSampleValueForIndex returns sample data based on parameter index
func (e *SQLExecutor) getSampleValueForIndex(index int) string {
	samples := []string{
		"'Sample Name'",
		"'New York'",
		"'test@example.com'",
		"'2024-01-01'",
		"75000",
		"1",
		"'ACTIVE'",
		"'Developer'",
		"40",
		"'555-1234'",
	}

	if index-1 < len(samples) {
		return samples[index-1]
	}
	return "'sample_value'"
}

// GetCreatedTables returns the list of tables created during execution
func (e *SQLExecutor) GetCreatedTables() []string {
	return e.createdTables
}

// Cleanup drops all created tables and objects
func (e *SQLExecutor) Cleanup() error {
	// Use the repository's cleanup method if available
	if e.repo != nil {
		return e.repo.CleanupDB()
	}

	// Fallback: drop created tables in reverse order
	for i := len(e.createdTables) - 1; i >= 0; i-- {
		tableName := e.createdTables[i]
		_, err := e.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
		if err != nil {
			return fmt.Errorf("dropping table %s: %w", tableName, err)
		}
	}

	return nil
}
