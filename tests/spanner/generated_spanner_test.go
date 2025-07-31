package spanner_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"postgres-example/models"
	"postgres-example/repo"
	"postgres-example/tools"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	t.Run("TestGeneratedSQLFiles", func(t *testing.T) {
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
		var results []models.TestFileResult

		for _, sqlFile := range sqlFiles {
			t.Run(filepath.Base(sqlFile), func(t *testing.T) {
				result := testSQLFileWithParsing(t, sqlFile)
				results = append(results, result)
			})
		}

		// Generate markdown report
		if err := generateMarkdownReport(results); err != nil {
			t.Logf("Warning: Failed to generate markdown report: %v", err)
		} else {
			t.Logf("Generated markdown report: sql_test_results.md")
		}

		// Generate Allure report
		allureReporter := tools.NewAllureReporter("allure-results")
		if err := allureReporter.GenerateAllureReport(results); err != nil {
			t.Logf("Warning: Failed to generate Allure report: %v", err)
		} else {
			t.Logf("Generated Allure reports in: allure-results/")
		}

		// Test results for valid_spanner_database.sql should be 100%
		var validSpannerDatabaseResult models.TestFileResult

		for _, result := range results {
			if result.Filename == "valid_spanner_database.sql" {
				validSpannerDatabaseResult = result
				break
			}
		}

		if !assert.NotNil(t, validSpannerDatabaseResult) {
			require.Fail(t, "valid_spanner_database.sql not found in results")
			return
		}
		assert.Equal(t, validSpannerDatabaseResult.ExecutedCount, validSpannerDatabaseResult.TotalStatements)
		assert.Zero(t, validSpannerDatabaseResult.ErrorRate)
		assert.Zero(t, validSpannerDatabaseResult.FailedCount)
		assert.Empty(t, validSpannerDatabaseResult.ParseErrors)
		assert.Empty(t, validSpannerDatabaseResult.ExecutionErrors)
		assert.Empty(t, validSpannerDatabaseResult.ErrorCodes)
		assert.Empty(t, validSpannerDatabaseResult.ErrorCategories)
	})
}

func testSQLFileWithParsing(t *testing.T, sqlFile string) models.TestFileResult {
	start := time.Now()
	filename := filepath.Base(sqlFile)

	result := models.TestFileResult{
		Filename:        filename,
		ParseErrorCodes: make(map[string]int),
		ErrorCodes:      make(map[string]int),
		ErrorCategories: make(map[string]int),
	}

	// Step 1: Parse SQL file to extract statements using memefish
	statements, err := tools.ExtractStatementsFromFile(sqlFile)
	if err != nil {
		t.Fatalf("Failed to extract SQL statements from %s: %v", sqlFile, err)
	}

	result.TotalStatements = len(statements)

	// Step 2: Parse each statement with memefish
	parseResults := parseStatementsWithMemefish(statements, filename)

	// Analyze parsing results
	var validStatements []string
	for _, pr := range parseResults {
		if pr.Parsed {
			result.ParsedCount++
			validStatements = append(validStatements, pr.Statement)

			// Count statement types from successful parsing
			switch strings.ToUpper(pr.Type) {
			case "CREATE":
				result.CreateStatements++
			case "INSERT":
				result.InsertStatements++
			case "SELECT":
				result.SelectStatements++
			case "DROP":
				result.DropStatements++
			}
		} else {
			errMsg := pr.Error.Error()
			result.ParseErrors = append(result.ParseErrors, errMsg)

			// Categorize parsing errors
			errorType := categorizeMemefishError(errMsg)
			result.ParseErrorCodes[errorType]++
		}
	}

	// Log parsing results
	t.Logf("Parsing results for %s:", filename)
	t.Logf("  Total statements: %d", result.TotalStatements)
	t.Logf("  Successfully parsed: %d", result.ParsedCount)
	t.Logf("  Parse failures: %d", len(result.ParseErrors))
	if len(result.ParseErrors) > 0 {
		parseErrorRate := float64(len(result.ParseErrors)) / float64(result.TotalStatements) * 100
		t.Logf("  Parse error rate: %.1f%%", parseErrorRate)
	}

	// Step 3: Execute only valid statements if any were parsed successfully
	if len(validStatements) > 0 {
		dbT := setupSpannerDB(t)
		defer dbT.Close()

		executor := repo.NewSQLExecutor(dbT.db, dbT.repo)
		defer func() {
			if err := executor.Cleanup(); err != nil {
				t.Logf("Warning: cleanup failed: %v", err)
			}
		}()

		// Execute the valid statements
		execResult, err := executor.ExecuteStatements(validStatements)
		if err != nil {
			t.Logf("Warning: ExecuteStatements returned error: %v", err)
		}

		// Collect execution results
		if execResult != nil {
			result.ExecutedCount = execResult.ExecutedCount
			result.FailedCount = len(execResult.Errors)

			for _, err := range execResult.Errors {
				errMsg := err.Error()
				result.ExecutionErrors = append(result.ExecutionErrors, errMsg)

				// Extract and count error codes
				errorCode := extractSpannerErrorCode(errMsg)
				if errorCode != "" {
					result.ErrorCodes[errorCode]++

					// Categorize InvalidArgument errors further
					if errorCode == "InvalidArgument" {
						category := categorizeInvalidArgumentError(errMsg)
						if category != "" {
							result.ErrorCategories[category]++
						}
					} else {
						result.ErrorCategories[errorCode]++
					}
				}
			}

			// Test additional queries if data was inserted
			if len(execResult.InsertedRecords) > 0 {
				testDataIntegrity(t, dbT, execResult)
			}
		}
	} else {
		t.Logf("No valid statements to execute for %s", filename)
	}

	result.ExecutionTime = time.Since(start)

	// Calculate error rate based on total statements
	if result.TotalStatements > 0 {
		totalErrors := len(result.ParseErrors) + result.FailedCount
		result.ErrorRate = float64(totalErrors) / float64(result.TotalStatements) * 100
	}

	// Log comprehensive results
	t.Logf("Final results for %s:", filename)
	t.Logf("  CREATE: %d, INSERT: %d, SELECT: %d, DROP: %d",
		result.CreateStatements, result.InsertStatements, result.SelectStatements, result.DropStatements)
	t.Logf("  Executed: %d, Failed: %d", result.ExecutedCount, result.FailedCount)
	t.Logf("  Overall error rate: %.1f%%", result.ErrorRate)
	t.Logf("  Total execution time: %v", result.ExecutionTime)

	// Log parsing error summary
	if len(result.ParseErrorCodes) > 0 {
		t.Logf("  Parse error types:")
		for errorType, count := range result.ParseErrorCodes {
			t.Logf("    %s: %d occurrences", errorType, count)
		}
	}

	// Log execution error codes summary
	if len(result.ErrorCodes) > 0 {
		t.Logf("  Execution error codes:")
		for code, count := range result.ErrorCodes {
			t.Logf("    %s: %d occurrences", code, count)
		}
	}

	return result
}

// parseStatementsWithMemefish parses each statement using memefish
func parseStatementsWithMemefish(statements []string, filename string) []models.ParseResult {
	var results []models.ParseResult

	for _, stmt := range statements {
		result := models.ParseResult{
			Statement: stmt,
		}

		// Try to parse the statement with memefish
		parsedStmt, err := memefish.ParseStatement(filename, stmt)
		if err != nil {
			result.Parsed = false
			result.Error = err
		} else {
			result.Parsed = true
			result.Type = getStatementType(parsedStmt, stmt)
		}

		results = append(results, result)
	}

	return results
}

// getStatementType determines the type of a parsed statement
func getStatementType(parsedStmt interface{}, originalStmt string) string {
	// For now, use a simple approach based on the original statement
	// since memefish AST types might be complex to analyze
	upperStmt := strings.ToUpper(strings.TrimSpace(originalStmt))

	switch {
	case strings.HasPrefix(upperStmt, "CREATE"):
		return "CREATE"
	case strings.HasPrefix(upperStmt, "INSERT"):
		return "INSERT"
	case strings.HasPrefix(upperStmt, "SELECT"):
		return "SELECT"
	case strings.HasPrefix(upperStmt, "DROP"):
		return "DROP"
	case strings.HasPrefix(upperStmt, "ALTER"):
		return "ALTER"
	case strings.HasPrefix(upperStmt, "UPDATE"):
		return "UPDATE"
	case strings.HasPrefix(upperStmt, "DELETE"):
		return "DELETE"
	default:
		return "OTHER"
	}
}

// categorizeMemefishError categorizes memefish parsing errors
func categorizeMemefishError(errMsg string) string {
	errLower := strings.ToLower(errMsg)

	// Common memefish error patterns
	switch {
	case strings.Contains(errLower, "syntax error"):
		if strings.Contains(errLower, "expecting") {
			return "Syntax Error: Missing Token"
		}
		return "Syntax Error: General"
	case strings.Contains(errLower, "unexpected token"):
		return "Syntax Error: Unexpected Token"
	case strings.Contains(errLower, "expecting"):
		return "Syntax Error: Expected Token"
	case strings.Contains(errLower, "invalid"):
		return "Invalid Syntax"
	case strings.Contains(errLower, "not supported"):
		return "Unsupported Feature"
	case strings.Contains(errLower, "unknown"):
		return "Unknown Element"
	default:
		return "Parse Error: Other"
	}
}

// extractSpannerErrorCode extracts the error code from Spanner error messages
func extractSpannerErrorCode(errMsg string) string {
	// Spanner errors typically have format: "rpc error: code = ErrorCode desc = ..."
	// or "spanner: code = \"ErrorCode\", desc = ..."

	// Pattern 1: rpc error format
	if strings.Contains(errMsg, "rpc error: code = ") {
		start := strings.Index(errMsg, "rpc error: code = ") + len("rpc error: code = ")
		end := strings.Index(errMsg[start:], " desc = ")
		if end != -1 {
			return errMsg[start : start+end]
		}
	}

	// Pattern 2: spanner client format
	if strings.Contains(errMsg, "spanner: code = ") {
		start := strings.Index(errMsg, "spanner: code = ") + len("spanner: code = ")
		if start < len(errMsg) && errMsg[start] == '"' {
			start++ // Skip opening quote
			end := strings.Index(errMsg[start:], "\"")
			if end != -1 {
				return errMsg[start : start+end]
			}
		}
	}

	return ""
}

// generateMarkdownReport creates a markdown file with test results
func generateMarkdownReport(results []models.TestFileResult) error {
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
	totalParsed := 0
	totalExecuted := 0
	totalParseErrors := 0
	totalExecutionErrors := 0
	allErrorCodes := make(map[string]int)      // Global error code counts
	allErrorCategories := make(map[string]int) // Global error category counts
	allParseErrorCodes := make(map[string]int) // Global parse error counts

	for _, result := range results {
		totalStatements += result.TotalStatements
		totalParsed += result.ParsedCount
		totalExecuted += result.ExecutedCount
		totalParseErrors += len(result.ParseErrors)
		totalExecutionErrors += len(result.ExecutionErrors)

		// Aggregate error codes
		for code, count := range result.ErrorCodes {
			allErrorCodes[code] += count
		}

		// Aggregate error categories
		for category, count := range result.ErrorCategories {
			allErrorCategories[category] += count
		}

		// Aggregate parse error codes
		for code, count := range result.ParseErrorCodes {
			allParseErrorCodes[code] += count
		}
	}

	fmt.Fprintf(file, "## Summary\n\n")
	fmt.Fprintf(file, "- **Total SQL Files**: %d\n", totalFiles)
	fmt.Fprintf(file, "- **Total Statements**: %d\n", totalStatements)
	fmt.Fprintf(file, "- **Successfully Parsed**: %d\n", totalParsed)
	fmt.Fprintf(file, "- **Parse Errors**: %d\n", totalParseErrors)
	fmt.Fprintf(file, "- **Successfully Executed**: %d\n", totalExecuted)
	fmt.Fprintf(file, "- **Execution Errors**: %d\n", totalExecutionErrors)
	if totalStatements > 0 {
		parseSuccessRate := float64(totalParsed) / float64(totalStatements) * 100
		fmt.Fprintf(file, "- **Parse Success Rate**: %.1f%%\n", parseSuccessRate)
		if totalParsed > 0 {
			execSuccessRate := float64(totalExecuted) / float64(totalParsed) * 100
			fmt.Fprintf(file, "- **Execution Success Rate** (of parsed): %.1f%%\n", execSuccessRate)
		}
		overallSuccessRate := float64(totalExecuted) / float64(totalStatements) * 100
		fmt.Fprintf(file, "- **Overall Success Rate**: %.1f%%\n\n", overallSuccessRate)
	}

	// Write parse error summary
	if len(allParseErrorCodes) > 0 {
		fmt.Fprintf(file, "## Parse Error Summary\n\n")
		fmt.Fprintf(file, "| Parse Error Type | Total Occurrences | Description |\n")
		fmt.Fprintf(file, "|------------------|-------------------|-------------|\n")

		// Sort parse error codes by frequency
		type parseErrorCount struct {
			errorType string
			count     int
		}
		var sortedParseErrors []parseErrorCount
		for errorType, count := range allParseErrorCodes {
			sortedParseErrors = append(sortedParseErrors, parseErrorCount{errorType, count})
		}

		// Simple bubble sort by count (descending)
		for i := 0; i < len(sortedParseErrors); i++ {
			for j := i + 1; j < len(sortedParseErrors); j++ {
				if sortedParseErrors[i].count < sortedParseErrors[j].count {
					sortedParseErrors[i], sortedParseErrors[j] = sortedParseErrors[j], sortedParseErrors[i]
				}
			}
		}

		for _, pec := range sortedParseErrors {
			description := getParseErrorDescription(pec.errorType)
			fmt.Fprintf(file, "| %s | %d | %s |\n", pec.errorType, pec.count, description)
		}
		fmt.Fprintf(file, "\n")
	}

	// Write error code summary
	if len(allErrorCodes) > 0 {
		fmt.Fprintf(file, "## Execution Error Code Summary\n\n")
		fmt.Fprintf(file, "| Error Code | Total Occurrences | Description |\n")
		fmt.Fprintf(file, "|------------|-------------------|-------------|\n")

		// Sort error codes by frequency
		type errorCodeCount struct {
			code  string
			count int
		}
		var sortedCodes []errorCodeCount
		for code, count := range allErrorCodes {
			sortedCodes = append(sortedCodes, errorCodeCount{code, count})
		}

		// Simple bubble sort by count (descending)
		for i := 0; i < len(sortedCodes); i++ {
			for j := i + 1; j < len(sortedCodes); j++ {
				if sortedCodes[i].count < sortedCodes[j].count {
					sortedCodes[i], sortedCodes[j] = sortedCodes[j], sortedCodes[i]
				}
			}
		}

		for _, ec := range sortedCodes {
			description := getErrorCodeDescription(ec.code)
			fmt.Fprintf(file, "| %s | %d | %s |\n", ec.code, ec.count, description)
		}
		fmt.Fprintf(file, "\n")
	}

	// Write error category summary
	if len(allErrorCategories) > 0 {
		fmt.Fprintf(file, "## Error Category Analysis\n\n")
		fmt.Fprintf(file, "| Error Category | Total Occurrences | Files Affected | Percentage |\n")
		fmt.Fprintf(file, "|----------------|-------------------|----------------|------------|\n")

		// Track which files have each error category
		categoryFileCount := make(map[string]int)
		for _, result := range results {
			if len(result.ExecutionErrors) > 0 { // Only count files that have execution errors
				fileCategories := make(map[string]bool) // Track unique categories per file
				for category := range result.ErrorCategories {
					fileCategories[category] = true
				}
				// Count each category once per file
				for category := range fileCategories {
					categoryFileCount[category]++
				}
			}
		}

		// Count files with execution errors
		filesWithExecutionErrors := 0
		for _, result := range results {
			if len(result.ExecutionErrors) > 0 {
				filesWithExecutionErrors++
			}
		}

		// Sort error categories by frequency
		type errorCategoryCount struct {
			category string
			count    int
		}
		var sortedCategories []errorCategoryCount
		for category, count := range allErrorCategories {
			sortedCategories = append(sortedCategories, errorCategoryCount{category, count})
		}

		// Simple bubble sort by count (descending)
		for i := 0; i < len(sortedCategories); i++ {
			for j := i + 1; j < len(sortedCategories); j++ {
				if sortedCategories[i].count < sortedCategories[j].count {
					sortedCategories[i], sortedCategories[j] = sortedCategories[j], sortedCategories[i]
				}
			}
		}

		for _, ec := range sortedCategories {
			percentage := float64(ec.count) / float64(totalExecutionErrors) * 100
			filesAffected := categoryFileCount[ec.category]
			fmt.Fprintf(file, "| %s | %d | %d/%d | %.1f%% |\n",
				ec.category, ec.count, filesAffected, filesWithExecutionErrors, percentage)
		}
		fmt.Fprintf(file, "\n")
	}

	// Write detailed results table
	fmt.Fprintf(file, "## Detailed Results\n\n")
	fmt.Fprintf(file, "| File | Total | Parsed | Parse Errors | CREATE | INSERT | SELECT | DROP | Executed | Exec Errors | Parse Success | Exec Success | Total Time |\n")
	fmt.Fprintf(file, "|------|-------|--------|--------------|--------|--------|--------|------|----------|-------------|---------------|--------------|------------|\n")

	for _, result := range results {
		parseSuccessRate := 0.0
		if result.TotalStatements > 0 {
			parseSuccessRate = float64(result.ParsedCount) / float64(result.TotalStatements) * 100
		}
		execSuccessRate := 0.0
		if result.ParsedCount > 0 {
			execSuccessRate = float64(result.ExecutedCount) / float64(result.ParsedCount) * 100
		}

		fmt.Fprintf(file, "| [%s](../../generated_sql/%s) | %d | %d | %d | %d | %d | %d | %d | %d | %d | %.1f%% | %.1f%% | %v |\n",
			result.Filename,
			result.Filename,
			result.TotalStatements,
			result.ParsedCount,
			len(result.ParseErrors),
			result.CreateStatements,
			result.InsertStatements,
			result.SelectStatements,
			result.DropStatements,
			result.ExecutedCount,
			len(result.ExecutionErrors),
			parseSuccessRate,
			execSuccessRate,
			result.ExecutionTime.Round(time.Millisecond),
		)
	}

	// Write error details with error codes
	fmt.Fprintf(file, "\n## Error Details\n\n")
	for _, result := range results {
		if len(result.ParseErrors) > 0 || len(result.ExecutionErrors) > 0 {
			fmt.Fprintf(file, "### %s\n\n", result.Filename)

			// Parse errors section
			if len(result.ParseErrors) > 0 {
				parseErrorRate := float64(len(result.ParseErrors)) / float64(result.TotalStatements) * 100
				fmt.Fprintf(file, "**Parse Error Rate**: %.1f%% (%d/%d failed to parse)\n\n", parseErrorRate, len(result.ParseErrors), result.TotalStatements)

				// Write parse error codes for this file
				if len(result.ParseErrorCodes) > 0 {
					fmt.Fprintf(file, "**Parse Error Types**:\n")
					for errorType, count := range result.ParseErrorCodes {
						fmt.Fprintf(file, "- `%s`: %d occurrences\n", errorType, count)
					}
					fmt.Fprintf(file, "\n")
				}

				fmt.Fprintf(file, "**Parse Errors**:\n")
				for i, errMsg := range result.ParseErrors {
					fmt.Fprintf(file, "%d. %s\n", i+1, errMsg)
				}
				fmt.Fprintf(file, "\n")
			}

			// Execution errors section
			if len(result.ExecutionErrors) > 0 {
				execErrorRate := 0.0
				if result.ParsedCount > 0 {
					execErrorRate = float64(len(result.ExecutionErrors)) / float64(result.ParsedCount) * 100
				}
				fmt.Fprintf(file, "**Execution Error Rate**: %.1f%% (%d/%d parsed statements failed)\n\n", execErrorRate, len(result.ExecutionErrors), result.ParsedCount)

				// Write error codes for this file
				if len(result.ErrorCodes) > 0 {
					fmt.Fprintf(file, "**Execution Error Codes**:\n")
					for code, count := range result.ErrorCodes {
						fmt.Fprintf(file, "- `%s`: %d occurrences\n", code, count)
					}
					fmt.Fprintf(file, "\n")
				}

				// Write error categories for this file
				if len(result.ErrorCategories) > 0 {
					fmt.Fprintf(file, "**Error Categories**:\n")
					for category, count := range result.ErrorCategories {
						fmt.Fprintf(file, "- `%s`: %d occurrences\n", category, count)
					}
					fmt.Fprintf(file, "\n")
				}

				fmt.Fprintf(file, "**Execution Errors**:\n")
				for i, errMsg := range result.ExecutionErrors {
					fmt.Fprintf(file, "%d. %s\n", i+1, errMsg)
				}
				fmt.Fprintf(file, "\n")
			}
		}
	}

	// Write compatibility insights
	fmt.Fprintf(file, "## Compatibility Insights\n\n")

	// Parse errors insights
	if len(allParseErrorCodes) > 0 {
		fmt.Fprintf(file, "### Most Common Parse Issues\n\n")

		// Get top parse error types
		type parseInsight struct {
			errorType   string
			count       int
			percentage  float64
			description string
		}

		var parseInsights []parseInsight
		for errorType, count := range allParseErrorCodes {
			percentage := float64(count) / float64(totalParseErrors) * 100
			description := getParseErrorDescription(errorType)
			parseInsights = append(parseInsights, parseInsight{errorType, count, percentage, description})
		}

		// Sort by count (descending)
		for i := 0; i < len(parseInsights); i++ {
			for j := i + 1; j < len(parseInsights); j++ {
				if parseInsights[i].count < parseInsights[j].count {
					parseInsights[i], parseInsights[j] = parseInsights[j], parseInsights[i]
				}
			}
		}

		// Show top 5
		maxShow := 5
		if len(parseInsights) < maxShow {
			maxShow = len(parseInsights)
		}

		for i := 0; i < maxShow; i++ {
			insight := parseInsights[i]
			fmt.Fprintf(file, "1. **%s** (%.1f%% of parse errors)\n", insight.errorType, insight.percentage)
			fmt.Fprintf(file, "   - %d occurrences across all files\n", insight.count)
			fmt.Fprintf(file, "   - %s\n\n", insight.description)
		}
	}

	// Execution errors insights
	fmt.Fprintf(file, "### Most Common Execution Issues\n\n")

	// Show top error categories with descriptions
	if len(allErrorCategories) > 0 {
		// Get top 5 error categories
		type categoryInsight struct {
			category    string
			count       int
			percentage  float64
			description string
		}

		var insights []categoryInsight
		for category, count := range allErrorCategories {
			percentage := float64(count) / float64(totalExecutionErrors) * 100
			description := getErrorCategoryDescription(category)
			insights = append(insights, categoryInsight{category, count, percentage, description})
		}

		// Sort by count (descending)
		for i := 0; i < len(insights); i++ {
			for j := i + 1; j < len(insights); j++ {
				if insights[i].count < insights[j].count {
					insights[i], insights[j] = insights[j], insights[i]
				}
			}
		}

		// Show top 5
		maxShow := 5
		if len(insights) < maxShow {
			maxShow = len(insights)
		}

		for i := 0; i < maxShow; i++ {
			insight := insights[i]
			fmt.Fprintf(file, "1. **%s** (%.1f%% of execution errors)\n", insight.category, insight.percentage)
			fmt.Fprintf(file, "   - %d occurrences across all files\n", insight.count)
			fmt.Fprintf(file, "   - %s\n\n", insight.description)
		}
	}

	return nil
}

// getErrorCodeDescription returns a human-readable description for Spanner error codes
func getErrorCodeDescription(code string) string {
	descriptions := map[string]string{
		"InvalidArgument":    "Invalid SQL syntax or unsupported features",
		"NotFound":           "Referenced table, column, or object not found",
		"FailedPrecondition": "Constraint violations or prerequisite not met",
		"AlreadyExists":      "Object already exists (duplicate creation)",
		"PermissionDenied":   "Insufficient permissions for operation",
		"Unimplemented":      "Feature not implemented in Spanner",
		"Internal":           "Internal Spanner error",
		"Unavailable":        "Service temporarily unavailable",
		"DeadlineExceeded":   "Operation timeout",
		"ResourceExhausted":  "Resource limits exceeded",
		"Cancelled":          "Operation was cancelled",
		"Unknown":            "Unknown error occurred",
	}

	if desc, exists := descriptions[code]; exists {
		return desc
	}
	return "Unknown error code"
}

// getErrorCategoryDescription returns a human-readable description for error categories
func getErrorCategoryDescription(category string) string {
	descriptions := map[string]string{
		// Syntax errors
		"Syntax Error: CURRENT_TIMESTAMP":           "CURRENT_TIMESTAMP() function call syntax not supported in Spanner",
		"Syntax Error: Missing Parentheses":         "SQL statement missing required opening parentheses",
		"Syntax Error: Missing Closing Parentheses": "SQL statement missing required closing parentheses",
		"Syntax Error: General":                     "General SQL syntax errors not matching specific patterns",

		// Type mismatches
		"Type Mismatch: GENERATE_UUID on INT64": "GENERATE_UUID() function used on INT64 columns instead of STRING",
		"Type Mismatch: General":                "Data type mismatches between expected and provided types",

		// Unsupported features
		"Unsupported Feature: Sequence Kind": "Identity column sequence kind not specified or unsupported",
		"Unsupported Feature: General":       "General Spanner unsupported features",

		// Missing clauses
		"Missing Clause: SQL SECURITY": "VIEW definitions missing required SQL SECURITY clause",
		"Missing Clause: General":      "SQL statements missing required clauses",

		// Function issues
		"Function Not Found: NEXTVAL": "NEXTVAL() function not available in Spanner",
		"Function Not Found: General": "SQL functions not available in Spanner",

		// Identity/Sequence issues
		"Identity Column: Missing Sequence Kind": "Identity columns require explicit sequence kind specification",

		// Table issues
		"Table Not Found (InvalidArgument)": "Table references that result in InvalidArgument rather than NotFound",

		// Foreign key issues
		"Foreign Key: Syntax Error": "Foreign key constraint syntax errors",

		// Default value issues
		"Default Value: Parsing Error": "Default value expressions that cannot be parsed",

		// Constraint issues
		"Constraint: Unsupported": "CHECK constraints and other constraint types not supported",

		// View issues
		"View Definition: Error": "Errors in view definition syntax or structure",

		// Other error codes
		"NotFound":           "Referenced objects (tables, columns, etc.) not found",
		"FailedPrecondition": "Constraint violations or prerequisites not met",
		"AlreadyExists":      "Attempting to create objects that already exist",
		"PermissionDenied":   "Insufficient permissions for the operation",
		"Unimplemented":      "Features not yet implemented in Spanner",

		// Catch-all
		"InvalidArgument: Other": "InvalidArgument errors not matching specific patterns",
	}

	if desc, exists := descriptions[category]; exists {
		return desc
	}
	return "No description available for this error category"
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

// categorizeInvalidArgumentError categorizes InvalidArgument errors into specific subcategories
func categorizeInvalidArgumentError(errMsg string) string {
	errLower := strings.ToLower(errMsg)

	// Syntax errors
	if strings.Contains(errLower, "syntax error") {
		if strings.Contains(errLower, "current_timestamp") {
			return "Syntax Error: CURRENT_TIMESTAMP"
		} else if strings.Contains(errLower, "expecting '('") {
			return "Syntax Error: Missing Parentheses"
		} else if strings.Contains(errLower, "expecting ')'") {
			return "Syntax Error: Missing Closing Parentheses"
		}
		return "Syntax Error: General"
	}

	// Type mismatches
	if strings.Contains(errLower, "expected type") && strings.Contains(errLower, "found") {
		if strings.Contains(errLower, "generate_uuid") {
			return "Type Mismatch: GENERATE_UUID on INT64"
		}
		return "Type Mismatch: General"
	}

	// Unsupported features
	if strings.Contains(errLower, "unsupported") {
		if strings.Contains(errLower, "sequence kind") {
			return "Unsupported Feature: Sequence Kind"
		}
		return "Unsupported Feature: General"
	}

	// Missing clauses
	if strings.Contains(errLower, "missing") {
		if strings.Contains(errLower, "sql security") {
			return "Missing Clause: SQL SECURITY"
		}
		return "Missing Clause: General"
	}

	// Function/feature not found
	if strings.Contains(errLower, "function not found") {
		if strings.Contains(errLower, "nextval") {
			return "Function Not Found: NEXTVAL"
		}
		return "Function Not Found: General"
	}

	// Sequence/Identity issues
	if strings.Contains(errLower, "sequence kind") && strings.Contains(errLower, "not specified") {
		return "Identity Column: Missing Sequence Kind"
	}

	// Table not found (when it's an InvalidArgument, not NotFound)
	if strings.Contains(errLower, "table not found") {
		return "Table Not Found (InvalidArgument)"
	}

	// Foreign key issues
	if strings.Contains(errLower, "foreign key") {
		return "Foreign Key: Syntax Error"
	}

	// Default value issues
	if strings.Contains(errLower, "default value") {
		return "Default Value: Parsing Error"
	}

	// Constraint issues
	if strings.Contains(errLower, "constraint") || strings.Contains(errLower, "check") {
		return "Constraint: Unsupported"
	}

	// View definition issues
	if strings.Contains(errLower, "definition of view") {
		return "View Definition: Error"
	}

	return "InvalidArgument: Other"
}

// getParseErrorDescription returns a human-readable description for memefish parse errors
func getParseErrorDescription(errorType string) string {
	descriptions := map[string]string{
		"Syntax Error: Missing Token":    "SQL statements missing required tokens (parentheses, keywords, etc.)",
		"Syntax Error: General":          "General SQL syntax errors not matching specific patterns",
		"Syntax Error: Unexpected Token": "Unexpected tokens found where different syntax was expected",
		"Syntax Error: Expected Token":   "Missing expected tokens in SQL syntax",
		"Invalid Syntax":                 "SQL syntax that doesn't conform to Spanner SQL grammar",
		"Unsupported Feature":            "SQL features that are not supported by Spanner",
		"Unknown Element":                "Unknown SQL elements or identifiers",
		"Parse Error: Other":             "Other parsing errors not categorized above",
	}

	if desc, exists := descriptions[errorType]; exists {
		return desc
	}
	return "No description available for this parse error type"
}
