package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	integration "sql-parser/openai_integration"
	"sql-parser/tools"
)

func main() {
	os.Exit(run())
}

func run() int {
	start := time.Now()
	var (
		mode            = flag.String("mode", "iterative", "Mode: 'single' or 'iterative'")
		maxIterations   = flag.Int("iterations", 2, "Maximum iterations for iterative mode")
		outputFile      = flag.String("output", "", "Output file for generated SQL (optional)")
		verbose         = flag.Bool("verbose", false, "Verbose output")
		debugPrompt     = flag.Bool("debug-prompt", false, "Save prompts to file for debugging")
		saveAccumulated = flag.Bool("save-results", true, "Save results to accumulated JSON file for graphing")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: go run ./openai_integration/cmd/openai-pipeline [options]\n\n")
		fmt.Fprintf(os.Stderr, "This tool integrates OpenAI GPT-4o mini with SQL testing pipeline.\n\n")
		fmt.Fprintf(os.Stderr, "Environment variables required:\n")
		fmt.Fprintf(os.Stderr, "  OPENAI_API_KEY  - Your OpenAI API key\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDebugging:\n")
		fmt.Fprintf(os.Stderr, "  --debug-prompt saves all prompts to debug_prompts_<timestamp>.txt\n")
		fmt.Fprintf(os.Stderr, "\nData Collection:\n")
		fmt.Fprintf(os.Stderr, "  --save-results=false disables saving to pipeline_results.json for graphing\n")
	}

	flag.Parse()

	// Check for API key
	config := tools.Get()
	if config.OpenAIAPIKey == "" {
		fmt.Fprintf(os.Stderr, "Error: OPENAI_API_KEY environment variable not set\n")
		return 2
	}

	// Get base path (parent directory)
	basePath, err := getBasePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting base path: %v\n", err)
		return 2
	}

	// Create pipeline
	pipelineStart := time.Now()
	pipeline, err := integration.NewPipeline(basePath, *maxIterations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating pipeline: %v\n", err)
		return 2
	}
	pipeline.SetVerbose(*verbose)
	pipeline.SetDebugPrompt(*debugPrompt)
	fmt.Printf("[%.3fs] Created pipeline\n", time.Since(pipelineStart).Seconds())

	fmt.Printf("\nStarting OpenAI integration pipeline in %s mode...\n", *mode)
	fmt.Printf("=== ITERATION PROGRESS ===\n")

	var result *integration.PipelineResult

	// Run pipeline based on mode
	executionStart := time.Now()
	switch *mode {
	case "single":
		result, err = pipeline.RunSingleShot()
	case "iterative":
		result, err = pipeline.RunIterative()
	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid mode '%s'. Use 'single' or 'iterative'\n", *mode)
		return 2
	}
	executionTime := time.Since(executionStart)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Pipeline error: %v\n", err)
		return 2
	}

	fmt.Printf("\n[%.3fs] Pipeline execution completed\n", executionTime.Seconds())

	// Save results to accumulated JSON file (if enabled)
	if *saveAccumulated {
		saveAccumStart := time.Now()
		if err := pipeline.AddResultToAccumulated(result, *mode); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save to accumulated results: %v\n", err)
		} else {
			fmt.Printf("[%.3fs] Results added to accumulated file\n", time.Since(saveAccumStart).Seconds())
		}
	}

	// Print results
	printResults(result, *verbose)

	// Save output file if specified
	if *outputFile != "" {
		saveStart := time.Now()
		if err := pipeline.SaveResultToFile(result, *outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save output file: %v\n", err)
		} else {
			fmt.Printf("[%.3fs] Generated SQL saved to: %s\n", time.Since(saveStart).Seconds(), filepath.Join("generated_sql", *outputFile))
		}
	}

	fmt.Printf("[%.3fs] Total execution time\n", time.Since(start).Seconds())

	// Return appropriate exit code
	if result.Success {
		return 0
	}
	return 1
}

func printResults(result *integration.PipelineResult, verbose bool) {
	fmt.Printf("\n=== PIPELINE RESULTS ===\n")
	fmt.Printf("Success: %v\n", result.Success)
	fmt.Printf("Iterations: %d\n", result.Iterations)
	fmt.Printf("Total time: %v\n", result.TotalTime.Round(time.Millisecond))
	fmt.Printf("Session ID: %s\n\n", result.SessionID)

	fmt.Printf("=== TEST RESULTS ===\n")
	fmt.Printf("Total statements: %d\n", result.TestResults.TotalStatements)
	fmt.Printf("Successfully parsed: %d\n", result.TestResults.ParsedCount)
	fmt.Printf("Parse errors: %d\n", len(result.TestResults.ParseErrors))
	fmt.Printf("Successfully executed: %d\n", result.TestResults.ExecutedCount)
	fmt.Printf("Execution errors: %d\n", len(result.TestResults.ExecutionErrors))

	if result.TestResults.TotalStatements > 0 {
		parseRate := float64(result.TestResults.ParsedCount) / float64(result.TestResults.TotalStatements) * 100
		execRate := 0.0
		if result.TestResults.ParsedCount > 0 {
			execRate = float64(result.TestResults.ExecutedCount) / float64(result.TestResults.ParsedCount) * 100
		}
		overall := float64(result.TestResults.ExecutedCount) / float64(result.TestResults.TotalStatements) * 100
		fmt.Printf("Parse success rate: %.1f%%\n", parseRate)
		fmt.Printf("Execution success rate (of parsed): %.1f%%\n", execRate)
		fmt.Printf("Overall success rate: %.1f%%\n", overall)
	}

	if verbose {
		if len(result.TestResults.ParseErrors) > 0 {
			fmt.Printf("\n=== PARSE ERRORS ===\n")
			for _, e := range result.TestResults.ParseErrorDetails {
				stmt := e.Statement
				if len(stmt) > 100 {
					stmt = stmt[:50] + "..." + stmt[len(stmt)-50:]
				}
				fmt.Printf("- %s\n  Statement: %s\n", e.Error, stmt)
			}
		}

		if len(result.TestResults.ExecutionErrors) > 0 {
			fmt.Printf("\n=== EXECUTION ERRORS ===\n")
			for i, e := range result.TestResults.ExecutionErrors {
				fmt.Printf("%d. %s\n", i+1, e)
			}
		}

		fmt.Printf("\n=== GENERATED SQL ===\n")
		fmt.Println(result.GeneratedSQL)

		fmt.Printf("\n=== CONVERSATION HISTORY ===\n")
		for i, msg := range result.Messages {
			fmt.Printf("\n--- Message %d (%s) ---\n", i+1, msg.Role)
			content := msg.Content
			if len(content) > 500 && !verbose {
				content = content[:500] + "... [truncated]"
			}
			fmt.Println(content)
		}

		fmt.Printf("\n=== DETAILED RESULTS (JSON) ===\n")
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonData))
	}
}

func getBasePath() (string, error) {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// If we're in the openai_integration directory, go up one level
	if filepath.Base(wd) == "openai_integration" {
		return filepath.Dir(wd), nil
	}

	return wd, nil
}
