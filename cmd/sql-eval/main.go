package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sql-parser/models"
	"sql-parser/repo"
	"sql-parser/tools"
)

// runResult aggregates the outcome for a single file evaluation
type runResult struct {
	fileResult models.TestFileResult
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: go run ./cmd/sql-eval <path-to-sql-file>\n")
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		return 2
	}

	sqlFile := flag.Arg(0)
	if err := validatePath(sqlFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	result, err := evaluateSQLFile(sqlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
		return 2
	}

	printTerminalReport(result.fileResult)

	if len(result.fileResult.ParseErrors) == 0 && len(result.fileResult.ExecutionErrors) == 0 {
		return 0
	}
	return 1
}

func validatePath(p string) error {
	if strings.TrimSpace(p) == "" {
		return errors.New("path must not be empty")
	}
	info, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("cannot stat path: %w", err)
	}
	if info.IsDir() {
		return errors.New("path points to a directory; expected a .sql file")
	}
	if !strings.HasSuffix(strings.ToLower(info.Name()), ".sql") {
		return errors.New("file must have .sql extension")
	}
	return nil
}

func evaluateSQLFile(sqlFile string) (runResult, error) {
	start := time.Now()
	filename := filepath.Base(sqlFile)

	fr := models.TestFileResult{
		Filename:        filename,
		ParseErrorCodes: make(map[string]int),
		ErrorCodes:      make(map[string]int),
		ErrorCategories: make(map[string]int),
	}

	statements, err := tools.ExtractStatementsFromFile(sqlFile)
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

func printTerminalReport(fr models.TestFileResult) {
	fmt.Printf("SQL Evaluation Results\n\n")
	fmt.Printf("File: %s\n", fr.Filename)
	fmt.Printf("Total statements: %d\n", fr.TotalStatements)
	fmt.Printf("Successfully parsed: %d\n", fr.ParsedCount)
	fmt.Printf("Parse errors: %d\n", len(fr.ParseErrors))
	fmt.Printf("Executed: %d\n", fr.ExecutedCount)
	fmt.Printf("Execution errors: %d\n", len(fr.ExecutionErrors))
	if fr.TotalStatements > 0 {
		parseRate := float64(fr.ParsedCount) / float64(fr.TotalStatements) * 100
		execRate := 0.0
		if fr.ParsedCount > 0 {
			execRate = float64(fr.ExecutedCount) / float64(fr.ParsedCount) * 100
		}
		overall := float64(fr.ExecutedCount) / float64(fr.TotalStatements) * 100
		fmt.Printf("Parse success rate: %.1f%%\n", parseRate)
		fmt.Printf("Execution success rate (of parsed): %.1f%%\n", execRate)
		fmt.Printf("Overall success rate: %.1f%%\n", overall)
	}
	fmt.Printf("Total time: %v\n", fr.ExecutionTime.Round(time.Millisecond))

	if len(fr.ParseErrorCodes) > 0 {
		fmt.Println()
		fmt.Println("Parse Error Summary:")
		for typ, cnt := range fr.ParseErrorCodes {
			desc := tools.GetParseErrorDescription(typ)
			fmt.Printf("- %s: %d (%s)\n", typ, cnt, desc)
		}
	}

	if len(fr.ErrorCodes) > 0 {
		fmt.Println()
		fmt.Println("Execution Error Code Summary:")
		for code, cnt := range fr.ErrorCodes {
			desc := tools.GetErrorCodeDescription(code)
			fmt.Printf("- %s: %d (%s)\n", code, cnt, desc)
		}
	}

	if len(fr.ErrorCategories) > 0 {
		fmt.Println()
		fmt.Println("Execution Error Categories:")
		for cat, cnt := range fr.ErrorCategories {
			desc := tools.GetErrorCategoryDescription(cat)
			fmt.Printf("- %s: %d (%s)\n", cat, cnt, desc)
		}
	}

	if len(fr.ParseErrors) > 0 {
		fmt.Println()
		fmt.Println("Parse Errors:")
		for _, e := range fr.ParseErrorDetails {
			stmt := e.Statement
			if len(stmt) > 200 {
				stmt = stmt[:100] + "..." + stmt[len(stmt)-100:]
			}
			fmt.Printf("- %s\n  Statement: %s\n", e.Error, stmt)
		}
	}

	if len(fr.ExecutionErrors) > 0 {
		fmt.Println()
		fmt.Println("Execution Errors:")
		for i, e := range fr.ExecutionErrors {
			fmt.Printf("%d. %s\n", i+1, e)
		}
	}
}
