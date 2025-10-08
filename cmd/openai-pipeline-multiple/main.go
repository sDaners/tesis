package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	integration "sql-parser/openai_integration"
)

func main() {
	os.Exit(runMultiple())
}

func runMultiple() int {
	var (
		mode            = flag.String("mode", "iterative", "Mode: 'single' or 'iterative'")
		maxIterations   = flag.Int("iterations", 1, "Maximum iterations for iterative mode")
		numConcurrent   = flag.Int("concurrent", 3, "Number of concurrent pipeline instances to run")
		saveAccumulated = flag.Bool("save-results", true, "Save results to accumulated JSON file for graphing")
		shortPrompts    = flag.Bool("short-prompts", false, "Generate shorter iterative prompts by removing summaries and truncating error details")
		ragEnabled      = flag.Bool("rag", false, "Enable RAG mode: combine prompt.txt with spanner_sql_generation_guidelines.txt")
		saveOutput      = flag.Bool("save-output", false, "Save output to file")
		verbose         = flag.Bool("verbose", false, "Verbose output for each pipeline")
		model           = flag.String("model", "chatgpt-4o-latest", "OpenAI model to use")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: go run ./openai_integration/cmd/main_multiple.go [options]\n\n")
		fmt.Fprintf(os.Stderr, "This tool runs multiple concurrent OpenAI pipeline instances.\n\n")
		fmt.Fprintf(os.Stderr, "Environment variables required:\n")
		fmt.Fprintf(os.Stderr, "  OPENAI_API_KEY  - Your OpenAI API key\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "  --debug-prompt saves all prompts to debug_prompts_<timestamp>.txt\n")
		fmt.Fprintf(os.Stderr, "  --short-prompts generates shorter iterative prompts by removing summaries\n")
		fmt.Fprintf(os.Stderr, "  --concurrent specifies number of concurrent pipeline instances (default: 3)\n")
		fmt.Fprintf(os.Stderr, "  --save-output saves output to file\n")
	}

	flag.Parse()

	// Get base path (parent directory)
	basePath, err := integration.GetBasePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting base path: %v\n", err)
		return 2
	}

	fmt.Printf("Starting %d concurrent OpenAI pipeline instances...\n", *numConcurrent)
	fmt.Printf("Mode: %s | Model: %s | Iterations: %d | Short Prompts: %v | RAG: %v\n", *mode, *model, *maxIterations, *shortPrompts, *ragEnabled)
	fmt.Printf("=== CONCURRENT EXECUTION PROGRESS ===\n")

	start := time.Now()

	// Channel to collect results
	results := make(chan *PipelineExecutionResult, *numConcurrent)
	var wg sync.WaitGroup

	outputFile := ""
	if *saveOutput {
		outputFile = fmt.Sprintf("%s-%d-%s.sql", *model, *maxIterations, time.Now().Format("20060102150405"))
	}

	// Launch concurrent pipeline instances
	for i := 0; i < *numConcurrent; i++ {
		wg.Add(1)
		go func(instanceID int) {
			defer wg.Done()
			runPipelineInstance(instanceID, integration.PipelineConfig{
				Mode:            *mode,
				MaxIterations:   *maxIterations,
				OutputFile:      outputFile,
				Verbose:         *verbose,
				SaveAccumulated: *saveAccumulated,
				DebugPrompt:     false,
				ShortPrompts:    *shortPrompts,
				RAGEnabled:      *ragEnabled,
				UniqueID:        fmt.Sprintf("instance-%d", instanceID),
				Model:           *model,
			}, basePath, results)
		}(i + 1)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and process results
	var allResults []*PipelineExecutionResult
	for result := range results {
		allResults = append(allResults, result)
		if result.Error != nil {
			fmt.Printf("Instance %d: FAILED - %v\n", result.InstanceID, result.Error)
		} else {
			fmt.Printf("Instance %d: SUCCESS - Parse: %.1f%% | Overall: %.1f%% | Time: %v\n",
				result.InstanceID,
				result.ParseSuccessRate,
				result.OverallSuccessRate,
				result.ExecutionTime.Round(time.Second))
		}
	}

	totalTime := time.Since(start)

	// Print summary
	printConcurrentSummary(allResults, totalTime)

	// Return exit code based on results
	successCount := 0
	for _, result := range allResults {
		if result.Error == nil && result.Success {
			successCount++
		}
	}

	if successCount == len(allResults) {
		return 0 // All succeeded
	} else if successCount > 0 {
		return 1 // Partial success
	}
	return 2 // All failed
}

// PipelineExecutionResult holds the result of a single pipeline execution
type PipelineExecutionResult struct {
	InstanceID         int
	Success            bool
	ParseSuccessRate   float64
	OverallSuccessRate float64
	ExecutionTime      time.Duration
	Result             *integration.PipelineResult
	Error              error
}

// runPipelineInstance runs a single pipeline instance and sends the result to the channel
func runPipelineInstance(instanceID int, config integration.PipelineConfig, basePath string, results chan<- *PipelineExecutionResult) {
	start := time.Now()

	// Create pipeline runner for this instance
	runner := integration.NewPipelineRunner(config, basePath)
	result, _, err := runner.RunWithResults()

	executionResult := &PipelineExecutionResult{
		InstanceID:    instanceID,
		Success:       result != nil && result.Success,
		ExecutionTime: time.Since(start),
		Result:        result,
		Error:         err,
	}

	// Calculate success rates if result is available
	if result != nil && result.TestResults.TotalStatements > 0 {
		executionResult.ParseSuccessRate = float64(result.TestResults.ParsedCount) / float64(result.TestResults.TotalStatements) * 100
		executionResult.OverallSuccessRate = float64(result.TestResults.ExecutedCount) / float64(result.TestResults.TotalStatements) * 100
	}

	results <- executionResult
}

// printConcurrentSummary prints a summary of all concurrent executions
func printConcurrentSummary(results []*PipelineExecutionResult, totalTime time.Duration) {
	fmt.Printf("\n=== CONCURRENT EXECUTION SUMMARY ===\n")
	fmt.Printf("Total instances: %d\n", len(results))
	fmt.Printf("Total wall-clock time: %v\n", totalTime.Round(time.Millisecond))

	successCount := 0
	var totalParseRate, totalOverallRate float64
	var totalExecutionTime time.Duration

	for _, result := range results {
		if result.Error == nil {
			totalParseRate += result.ParseSuccessRate
			totalOverallRate += result.OverallSuccessRate
			totalExecutionTime += result.ExecutionTime
			if result.Success {
				successCount++
			}
		}
	}

	validResults := len(results)
	if validResults > 0 {
		fmt.Printf("Average parse success rate: %.1f%%\n", totalParseRate/float64(validResults))
		fmt.Printf("Average overall success rate: %.1f%%\n", totalOverallRate/float64(validResults))
		fmt.Printf("Average execution time per instance: %v\n", (totalExecutionTime / time.Duration(validResults)).Round(time.Second))
	}

	fmt.Printf("\n=== INDIVIDUAL RESULTS ===\n")
	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("Instance %d: FAILED - %v\n", result.InstanceID, result.Error)
		} else {
			fmt.Printf("Instance %d: Parse %.1f%% | Overall %.1f%% | Time %v\n",
				result.InstanceID,
				result.ParseSuccessRate,
				result.OverallSuccessRate,
				result.ExecutionTime.Round(time.Second),
			)
		}
	}
}
