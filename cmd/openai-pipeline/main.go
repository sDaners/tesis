package main

import (
	"flag"
	"fmt"
	"os"
	integration "sql-parser/openai_integration"
)

func main() {
	os.Exit(run())
}

func run() int {
	var (
		mode            = flag.String("mode", "iterative", "Mode: 'single' or 'iterative'")
		maxIterations   = flag.Int("iterations", 2, "Maximum iterations for iterative mode")
		outputFile      = flag.String("output", "", "Output file for generated SQL (optional)")
		verbose         = flag.Bool("verbose", false, "Verbose output")
		debugPrompt     = flag.Bool("debug-prompt", false, "Save prompts to file for debugging")
		saveAccumulated = flag.Bool("save-results", true, "Save results to accumulated JSON file for graphing")
		shortPrompts    = flag.Bool("short-prompts", false, "Generate shorter iterative prompts by removing summaries and truncating error details")
		ragEnabled      = flag.Bool("rag", false, "Enable RAG mode: combine prompt.txt with spanner_sql_generation_guidelines.txt")
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
		fmt.Fprintf(os.Stderr, "  --short-prompts generates shorter iterative prompts by removing summaries\n")
		fmt.Fprintf(os.Stderr, "\nData Collection:\n")
		fmt.Fprintf(os.Stderr, "  --save-results=false disables saving to pipeline_results.json for graphing\n")
	}

	flag.Parse()

	// Get base path (parent directory)
	basePath, err := integration.GetBasePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting base path: %v\n", err)
		return 2
	}

	// Create pipeline configuration
	config := integration.PipelineConfig{
		Mode:            *mode,
		MaxIterations:   *maxIterations,
		OutputFile:      *outputFile,
		Verbose:         *verbose,
		DebugPrompt:     *debugPrompt,
		SaveAccumulated: *saveAccumulated,
		ShortPrompts:    *shortPrompts,
		RAGEnabled:      *ragEnabled,
		UniqueID:        "", // Single instance doesn't need unique ID
	}

	// Create and run pipeline
	runner := integration.NewPipelineRunner(config, basePath)
	result, exitCode, err := runner.RunWithResults()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return exitCode
	}

	// Print results
	if result != nil {
		integration.PrintResults(result, *verbose)
	}

	return exitCode
}
