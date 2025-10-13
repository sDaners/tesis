package repo

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

// SQLExecutor provides an abstraction for executing parsed SQL statements
type SQLExecutor struct {
	DB            *sql.DB
	repo          Database
	executed      map[string]bool
	createdTables []string
	tableSchemas  map[string]map[string]string // table -> column -> type
	uuidRegistry  map[string]string            // logical_key -> actual_uuid
}

// NewSQLExecutor creates a new SQL executor instance
func NewSQLExecutor(db *sql.DB, repo Database) *SQLExecutor {
	return &SQLExecutor{
		DB:           db,
		repo:         repo,
		executed:     make(map[string]bool),
		tableSchemas: make(map[string]map[string]string),
		uuidRegistry: make(map[string]string),
	}
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
	ID        any // int64 or string
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
	var createStmts, insertStmts, selectStmts, dropStmts, otherStmts []string

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
		default:
			// Handle other statement types (UPDATE, DELETE, ALTER, etc.)
			otherStmts = append(otherStmts, stmt)
		}
	}

	// Execute in proper order: DROP, CREATE, INSERT (with dependency order), SELECT

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

	// 3. Execute INSERT statements in dependency order
	sortedInserts := e.sortInsertsByDependencies(insertStmts)
	for _, stmt := range sortedInserts {
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

	// 5. Execute other statement types (UPDATE, DELETE, ALTER, etc.)
	for _, stmt := range otherStmts {
		if err := e.executeOther(stmt); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("statement failed: %w", err))
		} else {
			result.ExecutedCount++
		}
	}

	result.SkippedCount = result.TotalStatements - result.ExecutedCount

	return result, nil
}

// executeOther executes other statement types (UPDATE, DELETE, ALTER, etc.)
func (e *SQLExecutor) executeOther(stmt string) error {
	// Clean up the statement
	cleanStmt := e.cleanStatement(stmt)

	// Check if this is a supported statement type
	upperStmt := strings.ToUpper(strings.TrimSpace(cleanStmt))

	// Support common statement types
	if strings.HasPrefix(upperStmt, "UPDATE") ||
		strings.HasPrefix(upperStmt, "DELETE") ||
		strings.HasPrefix(upperStmt, "ALTER") ||
		strings.HasPrefix(upperStmt, "GRANT") ||
		strings.HasPrefix(upperStmt, "REVOKE") ||
		strings.HasPrefix(upperStmt, "SET") {

		// Execute the statement
		_, err := e.DB.Exec(cleanStmt)
		if err != nil {
			return fmt.Errorf("executing %s statement: %w", strings.Fields(upperStmt)[0], err)
		}
		return nil
	}

	// For unsupported statement types, return an error
	stmtType := "UNKNOWN"
	if words := strings.Fields(upperStmt); len(words) > 0 {
		stmtType = words[0]
	}
	return fmt.Errorf("unsupported statement type: %s", stmtType)
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

	// Track table names and schemas for better parameter substitution
	if tableName := e.extractTableName(cleanStmt); tableName != "" {
		e.createdTables = append(e.createdTables, tableName)
		// Extract table schema for parameter substitution
		e.extractTableSchema(cleanStmt, tableName)
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

	// Get table name for ID tracking
	tableName := e.getTableFromInsert(populatedStmt)

	// Check if this is a THEN RETURN statement
	if strings.Contains(strings.ToUpper(populatedStmt), "THEN RETURN") {
		// Execute and capture the returned ID
		var returnedID string
		err = e.DB.QueryRow(populatedStmt).Scan(&returnedID)
		if err != nil {
			result.Error = fmt.Errorf("executing INSERT with THEN RETURN: %w", err)
			return result
		}

		// Store the returned ID in our registry for foreign key references
		switch tableName {
		case "departments":
			e.uuidRegistry["departments_pk"] = "'" + returnedID + "'"
		case "employees":
			e.uuidRegistry["employees_pk"] = "'" + returnedID + "'"
		case "projects":
			e.uuidRegistry["projects_pk"] = "'" + returnedID + "'"
		}

		result.ID = returnedID
	} else {
		// Regular INSERT without THEN RETURN
		// Execute the insert
		_, err = e.DB.Exec(populatedStmt)
		if err != nil {
			result.Error = fmt.Errorf("executing INSERT: %w", err)
			return result
		}

		result.ID = "success" // Use string to indicate success
	}

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

func (e *SQLExecutor) cleanStatement(stmt string) string {
	cleaned := strings.TrimSpace(stmt)

	cleaned = e.stripComments(cleaned)

	// Normalize internal whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned = re.ReplaceAllString(cleaned, " ")

	return cleaned
}

func (e *SQLExecutor) stripComments(stmt string) string {
	lines := strings.Split(stmt, "\n")
	var sqlLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "--") {
			continue
		}

		// Handle inline comments (remove everything after --)
		if commentIdx := strings.Index(trimmed, "--"); commentIdx != -1 {
			trimmed = strings.TrimSpace(trimmed[:commentIdx])
			if trimmed == "" {
				continue
			}
		}

		sqlLines = append(sqlLines, trimmed)
	}

	return strings.Join(sqlLines, " ")
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

func (e *SQLExecutor) extractTableSchema(stmt, tableName string) {
	if e.tableSchemas[tableName] == nil {
		e.tableSchemas[tableName] = make(map[string]string)
	}

	// Extract column definitions
	// Look for patterns like: column_name TYPE
	colRe := regexp.MustCompile(`(\w+)\s+(STRING\(\d+\)|INT64|FLOAT64|TIMESTAMP|DATE|BOOL)\s*(?:NOT\s+NULL|DEFAULT|,|\))`)
	matches := colRe.FindAllStringSubmatch(stmt, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			columnName := strings.ToLower(match[1])
			columnType := strings.ToUpper(match[2])
			e.tableSchemas[tableName][columnName] = columnType
		}
	}
}

func (e *SQLExecutor) getTableFromInsert(stmt string) string {
	re := regexp.MustCompile(`(?i)INSERT\s+INTO\s+(\w+)`)
	matches := re.FindStringSubmatch(stmt)
	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}
	return ""
}

// populateInsertParameters replaces parameter placeholders with sample data
func (e *SQLExecutor) populateInsertParameters(stmt string) (string, error) {
	// Get table name to understand column types
	tableName := e.getTableFromInsert(stmt)

	// Replace named parameters (@param) with sample data
	namedParamRe := regexp.MustCompile(`@(\w+)`)
	paramMap := make(map[string]string)

	populated := namedParamRe.ReplaceAllStringFunc(stmt, func(match string) string {
		paramName := match[1:] // Remove @
		if replacement, exists := paramMap[paramName]; exists {
			return replacement
		}

		// Generate sample data based on parameter name and table schema
		value := e.getSampleValueForParameter(paramName, tableName)
		paramMap[paramName] = value
		return value
	})

	// Replace positional parameters ($1, $2, etc.) with sample values
	positionalParamRe := regexp.MustCompile(`\$\d+`)
	argIndex := 1
	populated = positionalParamRe.ReplaceAllStringFunc(populated, func(match string) string {
		value := e.getSampleValueForIndex(argIndex, tableName)
		argIndex++
		return value
	})

	return populated, nil
}

// getSampleValueForParameter returns sample data based on parameter name and table schema
func (e *SQLExecutor) getSampleValueForParameter(paramName, tableName string) string {
	lower := strings.ToLower(paramName)

	// Check if we have schema information for this table
	if schema, exists := e.tableSchemas[tableName]; exists {
		if columnType, exists := schema[lower]; exists {
			return e.generateValueForType(paramName, columnType, tableName)
		}
	}

	// Fallback to name-based logic
	switch {
	case strings.Contains(lower, "id"):
		// Check if this looks like a UUID column
		if strings.Contains(lower, "dept") || strings.Contains(lower, "emp") || strings.Contains(lower, "project") {
			return e.generateUUIDForParameter(paramName, tableName)
		}
		return "1"
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
	case strings.Contains(lower, "start_date"):
		return "'2024-01-01'"
	case strings.Contains(lower, "end_date"):
		return "'2024-12-31'"
	case strings.Contains(lower, "hire_date"):
		return "'2024-01-01'"
	case strings.Contains(lower, "date"):
		return "'2024-06-01'"
	case strings.Contains(lower, "salary"):
		return "75000.0"
	case strings.Contains(lower, "budget"):
		return "50000.0"
	case strings.Contains(lower, "hours"):
		return "40"
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

// generateValueForType generates appropriate sample data based on column type
func (e *SQLExecutor) generateValueForType(paramName, columnType, tableName string) string {
	lower := strings.ToLower(paramName)

	switch {
	case strings.HasPrefix(columnType, "STRING"):
		// Check if this is likely a UUID column
		if strings.Contains(lower, "id") {
			return e.generateUUIDForParameter(paramName, tableName)
		}
		// Handle other string types based on name
		return e.getSampleValueForParameter(paramName, "")
	case columnType == "INT64":
		if strings.Contains(lower, "salary") || strings.Contains(lower, "budget") {
			return "75000"
		} else if strings.Contains(lower, "hours") {
			return "40"
		}
		return "1"
	case columnType == "FLOAT64":
		if strings.Contains(lower, "salary") || strings.Contains(lower, "budget") {
			return "75000.0"
		}
		return "1.0"
	case columnType == "TIMESTAMP":
		if strings.Contains(lower, "start") {
			return "'2024-01-01T00:00:00Z'"
		} else if strings.Contains(lower, "end") {
			return "'2024-12-31T23:59:59Z'"
		}
		return "'2024-06-01T12:00:00Z'"
	case columnType == "DATE":
		if strings.Contains(lower, "start") {
			return "'2024-01-01'"
		} else if strings.Contains(lower, "end") {
			return "'2024-12-31'"
		}
		return "'2024-06-01'"
	case columnType == "BOOL":
		return "true"
	default:
		return "'sample_value'"
	}
}

// generateUUIDForParameter generates UUID with proper foreign key relationships
func (e *SQLExecutor) generateUUIDForParameter(paramName, tableName string) string {
	lower := strings.ToLower(paramName)

	// Handle foreign key references by reusing captured UUIDs from THEN RETURN
	switch {
	case strings.Contains(lower, "dept_id"):
		if tableName != "departments" {
			// This is a foreign key reference to departments
			key := "departments_pk"
			if uuid, exists := e.uuidRegistry[key]; exists {
				return uuid
			}
		}
		return "NULL" // Let database generate UUID
	case strings.Contains(lower, "emp_id"):
		if tableName != "employees" {
			// This is a foreign key reference to employees
			key := "employees_pk"
			if uuid, exists := e.uuidRegistry[key]; exists {
				return uuid
			}
		}
		// For employees table or if no UUID captured yet, let database generate
		return "NULL" // Let database generate UUID
	case strings.Contains(lower, "project_id"):
		if tableName != "projects" {
			// This is a foreign key reference to projects
			key := "projects_pk"
			if uuid, exists := e.uuidRegistry[key]; exists {
				return uuid
			}
		}
		// For projects table or if no UUID captured yet, let database generate
		return "NULL" // Let database generate UUID
	case strings.Contains(lower, "manager_id"):
		// Manager is a self-reference to employees table
		key := "employees_pk"
		if uuid, exists := e.uuidRegistry[key]; exists {
			return uuid
		}
		// If no employee UUID captured yet, this will likely fail (which is expected)
		return "NULL"
	default:
		// For other ID columns, let database generate
		return "NULL"
	}
}

// getSampleValueForIndex returns sample data based on parameter index and table schema
func (e *SQLExecutor) getSampleValueForIndex(index int, tableName string) string {
	if schema, exists := e.tableSchemas[tableName]; exists {
		columnNames := make([]string, 0, len(schema))
		for colName := range schema {
			columnNames = append(columnNames, colName)
		}
		if index-1 < len(columnNames) {
			colName := columnNames[index-1]
			colType := schema[colName]
			return e.generateValueForType(colName, colType, tableName)
		}
	}

	// Fallback to hardcoded samples
	switch index {
	case 1:
		return "'Sample Name'"
	case 2:
		return "'New York'"
	case 3:
		return "'test@example.com'"
	case 4:
		return "'2024-01-01'"
	case 5:
		return "75000.0"
	case 6:
		// For UUID columns, let database generate
		return "NULL"
	case 7:
		return "'ACTIVE'"
	case 8:
		return "'Developer'"
	case 9:
		return "40"
	case 10:
		return "'555-1234'"
	default:
		return "'sample_value'"
	}
}

// Cleanup drops all created tables and objects
func (e *SQLExecutor) Cleanup() error {
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

// sortInsertsByDependencies sorts INSERT statements to respect foreign key dependencies
func (e *SQLExecutor) sortInsertsByDependencies(insertStmts []string) []string {
	// Define dependency order for common table patterns
	dependencyOrder := map[string]int{
		"departments":         1, // No dependencies
		"projects":            1, // No dependencies
		"employees":           2, // Depends on departments
		"project_assignments": 3, // Depends on employees and projects
	}

	// Create a slice of statements with their priority
	type stmtWithPriority struct {
		stmt     string
		priority int
		table    string
	}

	var stmtsWithPriority []stmtWithPriority

	for _, stmt := range insertStmts {
		tableName := e.getTableFromInsert(stmt)
		priority := dependencyOrder[tableName]
		if priority == 0 {
			priority = 999 // Unknown tables go last
		}

		stmtsWithPriority = append(stmtsWithPriority, stmtWithPriority{
			stmt:     stmt,
			priority: priority,
			table:    tableName,
		})
	}

	// Sort by priority
	for i := 0; i < len(stmtsWithPriority); i++ {
		for j := i + 1; j < len(stmtsWithPriority); j++ {
			if stmtsWithPriority[i].priority > stmtsWithPriority[j].priority {
				stmtsWithPriority[i], stmtsWithPriority[j] = stmtsWithPriority[j], stmtsWithPriority[i]
			}
		}
	}

	// Extract sorted statements
	var sortedStmts []string
	for _, s := range stmtsWithPriority {
		sortedStmts = append(sortedStmts, s.stmt)
	}

	return sortedStmts
}
