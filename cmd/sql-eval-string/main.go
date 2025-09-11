package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"sql-parser/models"
	"sql-parser/repo"
	"sql-parser/tools"
)

// runResult aggregates the outcome for a single string evaluation
type runResult struct {
	fileResult models.TestFileResult
}

func main() {
	os.Exit(run())
}

func run() int {
	// For now, just use a placeholder string - this will be replaced with actual input mechanism
	sqlContent := `-- Example SQL content
CREATE TABLE employees (
    id INT64 NOT NULL,
    name STRING(100),
    department_id INT64
) PRIMARY KEY (id);

INSERT INTO employees (id, name, department_id) VALUES (1, 'John Doe', 1);
SELECT * FROM employees;`

	result, err := evaluateSQLString(sqlContent, "mock")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
		return 2
	}

	// Generate the report as a string instead of printing directly
	reportString := generateTerminalReportString(result.fileResult)

	// TODO: In the future, this string can be:
	// - Returned from a function
	// - Sent over HTTP
	// - Stored in a database
	// - Processed by other systems
	// For now, print the string to console (will be removed later)
	fmt.Print(reportString)

	if len(result.fileResult.ParseErrors) == 0 && len(result.fileResult.ExecutionErrors) == 0 {
		return 0
	}
	return 1
}

func evaluateSQLString(content string, filename string) (runResult, error) {
	start := time.Now()

	fr := models.TestFileResult{
		Filename:        filename,
		ParseErrorCodes: make(map[string]int),
		ErrorCodes:      make(map[string]int),
		ErrorCategories: make(map[string]int),
	}

	statements, err := tools.ExtractStatementsFromString(content)
	if err != nil {
		return runResult{}, fmt.Errorf("extract statements: %w", err)
	}
	fr.TotalStatements = len(statements)

	parseResults := tools.ParseStatementsWithMemefish(statements, filename)

	var validStatements []string
	for _, pr := range parseResults {
		if pr.Parsed {
			fr.ParsedCount++
			validStatements = append(validStatements, pr.Statement)
			switch strings.ToUpper(pr.Type) {
			case "CREATE":
				fr.CreateStatements++
			case "INSERT":
				fr.InsertStatements++
			case "SELECT":
				fr.SelectStatements++
			case "DROP":
				fr.DropStatements++
			}
		} else {
			errMsg := pr.Error.Error()
			fr.ParseErrors = append(fr.ParseErrors, errMsg)
			fr.ParseErrorDetails = append(fr.ParseErrorDetails, models.ParseError{Statement: pr.Statement, Error: errMsg})
			errType := tools.CategorizeMemefishError(errMsg)
			fr.ParseErrorCodes[errType]++
		}
	}

	if len(validStatements) > 0 {
		db, terminate, err := tools.GetDB(true)
		if err != nil {
			return runResult{}, fmt.Errorf("connect DB: %w", err)
		}
		defer func() {
			_ = db.Close()
			if terminate != nil {
				terminate()
			}
		}()

		r := repo.NewSpannerRepo(db)
		_ = r.CleanupDB()

		executor := repo.NewSQLExecutor(db, r)
		defer func() { _ = executor.Cleanup() }()

		execResult, _ := executor.ExecuteStatements(validStatements)
		if execResult != nil {
			fr.ExecutedCount = execResult.ExecutedCount
			fr.FailedCount = len(execResult.Errors)
			for _, e := range execResult.Errors {
				errMsg := e.Error()
				fr.ExecutionErrors = append(fr.ExecutionErrors, errMsg)
				code := tools.ExtractSpannerErrorCode(errMsg)
				if code != "" {
					fr.ErrorCodes[code]++
					if code == "InvalidArgument" {
						category := tools.CategorizeInvalidArgumentError(errMsg)
						if category != "" {
							fr.ErrorCategories[category]++
						}
					} else {
						fr.ErrorCategories[code]++
					}
				}
			}
		}
	}

	fr.ExecutionTime = time.Since(start)
	if fr.TotalStatements > 0 {
		totalErrors := len(fr.ParseErrors) + fr.FailedCount
		fr.ErrorRate = float64(totalErrors) / float64(fr.TotalStatements) * 100
	}

	return runResult{fileResult: fr}, nil
}

func generateTerminalReportString(fr models.TestFileResult) string {
	var report strings.Builder

	report.WriteString("The generated sql code has gone through some testing, here are the results:\n\n")
	report.WriteString(fmt.Sprintf("Total statements: %d\n", fr.TotalStatements))
	report.WriteString(fmt.Sprintf("Successfully parsed: %d\n", fr.ParsedCount))
	report.WriteString(fmt.Sprintf("Parse errors: %d\n", len(fr.ParseErrors)))
	report.WriteString(fmt.Sprintf("Successfully executed: %d\n", fr.ExecutedCount))
	report.WriteString(fmt.Sprintf("Execution errors: %d\n", len(fr.ExecutionErrors)))

	if fr.TotalStatements > 0 {
		parseRate := float64(fr.ParsedCount) / float64(fr.TotalStatements) * 100
		execRate := 0.0
		if fr.ParsedCount > 0 {
			execRate = float64(fr.ExecutedCount) / float64(fr.ParsedCount) * 100
		}
		overall := float64(fr.ExecutedCount) / float64(fr.TotalStatements) * 100
		report.WriteString(fmt.Sprintf("Parse success rate: %.1f%%\n", parseRate))
		report.WriteString(fmt.Sprintf("Execution success rate (of parsed): %.1f%%\n", execRate))
		report.WriteString(fmt.Sprintf("Overall success rate: %.1f%%\n", overall))
	}
	report.WriteString(fmt.Sprintf("Total time: %v\n", fr.ExecutionTime.Round(time.Millisecond)))

	if len(fr.ParseErrorCodes) > 0 {
		report.WriteString("\n")
		report.WriteString("Parse Error Summary:\n")
		for typ, cnt := range fr.ParseErrorCodes {
			desc := tools.GetParseErrorDescription(typ)
			report.WriteString(fmt.Sprintf("- %s: %d (%s)\n", typ, cnt, desc))
		}
	}

	if len(fr.ErrorCodes) > 0 {
		report.WriteString("\n")
		report.WriteString("Execution Error Code Summary:\n")
		for code, cnt := range fr.ErrorCodes {
			desc := tools.GetErrorCodeDescription(code)
			report.WriteString(fmt.Sprintf("- %s: %d (%s)\n", code, cnt, desc))
		}
	}

	if len(fr.ErrorCategories) > 0 {
		report.WriteString("\n")
		report.WriteString("Execution Error Categories:\n")
		for cat, cnt := range fr.ErrorCategories {
			desc := tools.GetErrorCategoryDescription(cat)
			report.WriteString(fmt.Sprintf("- %s: %d (%s)\n", cat, cnt, desc))
		}
	}

	if len(fr.ParseErrors) > 0 {
		report.WriteString("\n")
		report.WriteString("Parse Errors:\n")
		for _, e := range fr.ParseErrorDetails {
			stmt := e.Statement
			if len(stmt) > 200 {
				stmt = stmt[:100] + "..." + stmt[len(stmt)-100:]
			}
			report.WriteString(fmt.Sprintf("- %s\n  Statement: %s\n", e.Error, stmt))
		}
	}

	if len(fr.ExecutionErrors) > 0 {
		report.WriteString("\n")
		report.WriteString("Execution Errors:\n")
		for i, e := range fr.ExecutionErrors {
			report.WriteString(fmt.Sprintf("%d. %s\n", i+1, e))
		}
	}

	// Add AI-specific recommendations
	recommendations := tools.GetAIRecommendations(fr)
	if len(recommendations) > 0 {
		report.WriteString("\n")
		report.WriteString("=== AI AGENT RECOMMENDATIONS ===\n")
		report.WriteString("ALWAYS: The response to this message should be a the entire sql code with fixes applied to it.\n")
		for _, rec := range recommendations {
			report.WriteString(rec + "\n")
		}
		report.WriteString("\n")
		report.WriteString("TIP: Focus on fixing parse errors first, as they prevent execution.\n")
		report.WriteString("TIP: Refer to Spanner SQL documentation for supported syntax and functions.\n")
	}

	return report.String()
}
