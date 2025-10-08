package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"sql-parser/tools"
)

// PipelineConfig holds configuration for running a pipeline instance
type PipelineConfig struct {
	Mode            string
	MaxIterations   int
	OutputFile      string
	Verbose         bool
	DebugPrompt     bool
	SaveAccumulated bool
	ShortPrompts    bool
	RAGEnabled      bool
	UniqueID        string
	Model           string
}

// PipelineRunner encapsulates the logic for running a single pipeline instance
type PipelineRunner struct {
	config   PipelineConfig
	basePath string
}

// NewPipelineRunner creates a new pipeline runner with the given configuration
func NewPipelineRunner(config PipelineConfig, basePath string) *PipelineRunner {
	return &PipelineRunner{
		config:   config,
		basePath: basePath,
	}
}

// Run executes the pipeline with the configured settings
func (pr *PipelineRunner) Run() (*PipelineResult, error) {
	start := time.Now()

	// Create pipeline
	pipelineStart := time.Now()
	pipeline, err := NewPipelineWithModel(pr.basePath, pr.config.MaxIterations, pr.config.Model, pr.config.Verbose)
	if err != nil {
		return nil, fmt.Errorf("error creating pipeline: %w", err)
	}
	pipeline.SetDebugPrompt(pr.config.DebugPrompt)
	pipeline.SetShortPrompts(pr.config.ShortPrompts)
	pipeline.SetRAGEnabled(pr.config.RAGEnabled)

	// Set unique ID if provided (for concurrent execution)
	if pr.config.UniqueID != "" {
		pipeline.SetUniqueID(pr.config.UniqueID)
	}

	if pr.config.Verbose {
		fmt.Printf("[%.3fs] Created pipeline\n", time.Since(pipelineStart).Seconds())
		fmt.Printf("\nStarting OpenAI integration pipeline in %s mode...\n", pr.config.Mode)
		fmt.Printf("=== ITERATION PROGRESS ===\n")
	}

	var result *PipelineResult

	// Run pipeline based on mode
	executionStart := time.Now()
	switch pr.config.Mode {
	case "single":
		result, err = pipeline.RunSingleShot()
	case "iterative":
		result, err = pipeline.RunIterative()
	default:
		return nil, fmt.Errorf("invalid mode '%s'. Use 'single' or 'iterative'", pr.config.Mode)
	}
	executionTime := time.Since(executionStart)

	if err != nil {
		return nil, fmt.Errorf("pipeline error: %w", err)
	}

	if pr.config.Verbose {
		fmt.Printf("\n[%.3fs] Pipeline execution completed\n", executionTime.Seconds())
	}

	// Save results to accumulated JSON file (if enabled)
	if pr.config.SaveAccumulated {
		saveAccumStart := time.Now()
		if err := pipeline.AddResultToAccumulated(result, pr.config.Mode); err != nil {
			fmt.Printf("Warning: Failed to save to accumulated results: %v\n", err)
		} else if pr.config.Verbose {
			fmt.Printf("[%.3fs] Results added to accumulated file\n", time.Since(saveAccumStart).Seconds())
		}
	}

	// Save output file if specified
	if pr.config.OutputFile != "" {
		saveStart := time.Now()
		outputFile := fmt.Sprintf("%s-%s", pr.config.UniqueID, pr.config.OutputFile)
		if err := pipeline.SaveResultToFile(result, outputFile); err != nil {
			fmt.Printf("Warning: Failed to save output file: %v\n", err)
		} else if pr.config.Verbose {
			fmt.Printf("[%.3fs] Generated SQL saved to: %s\n", time.Since(saveStart).Seconds(), filepath.Join("generated_sql", outputFile))
		}
	}

	if pr.config.Verbose {
		fmt.Printf("[%.3fs] Total execution time\n", time.Since(start).Seconds())
	}

	return result, nil
}

// RunWithResults runs the pipeline and returns formatted results
func (pr *PipelineRunner) RunWithResults() (result *PipelineResult, exitCode int, err error) {
	// Check for API key
	config := tools.Get()
	if config.OpenAIAPIKey == "" {
		return nil, 2, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	result, err = pr.Run()
	if err != nil {
		fmt.Printf("An error occurred while running the pipeline: %v\n", err)
		return result, 2, err
	}

	// Return appropriate exit code
	exitCode = 0
	if result != nil && !result.Success {
		exitCode = 1
	}

	return result, exitCode, nil
}

// PrintResults prints the pipeline results in a standard format
func PrintResults(result *PipelineResult, verbose bool) {
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

	// Verbose output (if requested)
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
	}
}

// GetBasePath returns the base path for the project
func GetBasePath() (string, error) {
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
