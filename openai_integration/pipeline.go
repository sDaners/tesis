package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"sql-parser/models"
	"sql-parser/repo"
	"sql-parser/tools"
)

// Global mutex to protect concurrent writes to pipeline_results.json
var pipelineResultsMutex sync.Mutex

// Pipeline orchestrates the complete workflow from prompt to tested SQL
type Pipeline struct {
	client        *OpenAIClient
	sessionMgr    *SessionManager
	promptReader  *PromptReader
	basePath      string
	maxIterations int
	verbose       bool
	debugPrompt   bool
	debugFile     string
	shortPrompts  bool
	uniqueID      string
	model         string
}

// NewPipeline creates a new pipeline instance with the default model
func NewPipeline(basePath string, maxIterations int, verbose bool) (*Pipeline, error) {
	return NewPipelineWithModel(basePath, maxIterations, "", verbose)
}

// NewPipelineWithModel creates a new pipeline instance with a specific model
func NewPipelineWithModel(basePath string, maxIterations int, model string, verbose bool) (*Pipeline, error) {
	var client *OpenAIClient

	// Create client with custom model
	config := OpenAIConfig{
		APIKey:      os.Getenv("OPENAI_API_KEY"),
		Model:       model,
		Temperature: 0.7,
		MaxTokens:   8096,
		BaseURL:     DefaultOpenAIURL,
		Verbose:     verbose,
	}
	client = NewOpenAIClient(config)

	sessionMgr := NewSessionManager(client)
	promptReader := NewPromptReader(basePath)

	return &Pipeline{
		client:        client,
		sessionMgr:    sessionMgr,
		promptReader:  promptReader,
		basePath:      basePath,
		maxIterations: maxIterations,
		verbose:       verbose,
		debugPrompt:   false,
		debugFile:     "",
		shortPrompts:  false,
		uniqueID:      "",
		model:         model,
	}, nil
}

// SetDebugPrompt enables debug mode and creates a debug file for saving prompts
func (p *Pipeline) SetDebugPrompt(debug bool) {
	p.debugPrompt = debug
	if debug {
		timestamp := time.Now().Format("20060102_150405")
		p.debugFile = fmt.Sprintf("debug_prompts_%s.txt", timestamp)

		// Initialize the debug file with header
		header := fmt.Sprintf("=== PROMPT DEBUG LOG ===\nTimestamp: %s\nMode: Debug\n\n", time.Now().Format("2006-01-02 15:04:05"))
		if err := os.WriteFile(p.debugFile, []byte(header), 0644); err != nil {
			fmt.Printf("Warning: Failed to create debug file %s: %v\n", p.debugFile, err)
			p.debugPrompt = false
		} else {
			fmt.Printf("Debug prompts will be saved to: %s\n", p.debugFile)
		}
	}
}

// SetShortPrompts enables/disables shorter prompt generation for iterative feedback
func (p *Pipeline) SetShortPrompts(short bool) {
	p.shortPrompts = short
}

// SetUniqueID sets a unique identifier for this pipeline instance (for concurrent execution)
func (p *Pipeline) SetUniqueID(uniqueID string) {
	p.uniqueID = uniqueID
}

// savePromptToDebugFile appends a prompt to the debug file
func (p *Pipeline) savePromptToDebugFile(promptType, content string) {
	if !p.debugPrompt || p.debugFile == "" {
		return
	}

	timestamp := time.Now().Format("15:04:05.000")
	debugEntry := fmt.Sprintf("[%s] === %s ===\n%s\n\n", timestamp, promptType, content)

	// Append to file
	file, err := os.OpenFile(p.debugFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Warning: Failed to open debug file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(debugEntry); err != nil {
		fmt.Printf("Warning: Failed to write to debug file: %v\n", err)
	}
}

// printIterationResult prints the results of a single iteration in real-time
func (p *Pipeline) printIterationResult(iteration int, testResult models.TestFileResult) {
	parseRate := 0.0
	execRate := 0.0
	overall := 0.0

	if testResult.TotalStatements > 0 {
		parseRate = float64(testResult.ParsedCount) / float64(testResult.TotalStatements) * 100
		if testResult.ParsedCount > 0 {
			execRate = float64(testResult.ExecutedCount) / float64(testResult.ParsedCount) * 100
		}
		overall = float64(testResult.ExecutedCount) / float64(testResult.TotalStatements) * 100
	}

	fmt.Printf("Iteration %d: Parse %.1f%% | Execution %.1f%% | Overall %.1f%%\n",
		iteration, parseRate, execRate, overall)
}

// RunSingleShot runs a single-shot prompt without iteration
func (p *Pipeline) RunSingleShot() (*PipelineResult, error) {
	start := time.Now()

	// Read the initial prompt
	initialPrompt, err := p.promptReader.ReadPromptFile()
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt: %w", err)
	}

	// Save prompt to debug file if enabled
	p.savePromptToDebugFile("INITIAL PROMPT (Single Shot)", initialPrompt)

	// Send to OpenAI
	fmt.Printf("  └─ Sending prompt to AI...\n")
	aiStart := time.Now()
	response, err := p.client.SendSingleMessage(initialPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to send message to OpenAI: %w", err)
	}
	fmt.Printf("  └─ [%.3fs] AI response received\n", time.Since(aiStart).Seconds())

	// Save AI response to debug file if enabled
	p.savePromptToDebugFile("AI RESPONSE (Single Shot)", response)

	// Extract SQL from response
	generatedSQL := p.promptReader.ExtractSQLFromResponse(response)

	// Test the generated SQL
	testStart := time.Now()
	testResult, err := p.testSQLString(generatedSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to test SQL: %w", err)
	}
	fmt.Printf("  └─ [%.3fs] SQL testing completed\n", time.Since(testStart).Seconds())

	success := len(testResult.ParseErrors) == 0 && len(testResult.ExecutionErrors) == 0

	// Create single iteration result
	iterationResult := IterationResult{
		Iteration:    1,
		TestResults:  testResult,
		Success:      success,
		GeneratedSQL: generatedSQL,
	}

	// Print iteration result in real-time
	p.printIterationResult(1, testResult)

	result := &PipelineResult{
		SessionID:        "single-shot",
		InitialPrompt:    initialPrompt,
		GeneratedSQL:     generatedSQL,
		TestResults:      testResult,
		Iterations:       1,
		IterationResults: []IterationResult{iterationResult},
		Success:          success,
		Messages: []ConversationMessage{
			{Role: "user", Content: initialPrompt},
			{Role: "assistant", Content: response},
		},
		TotalTime:  time.Since(start),
		TokensUsed: 0, // Would need to track from OpenAI response
	}

	return result, nil
}

// RunIterative runs the pipeline with iterative feedback until success or max iterations
func (p *Pipeline) RunIterative() (*PipelineResult, error) {
	start := time.Now()

	// Create a new session
	session, err := p.sessionMgr.CreateSession(DefaultModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Read initial prompt
	initialPrompt, err := p.promptReader.ReadPromptFile()
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt: %w", err)
	}

	// Save initial prompt to debug file if enabled
	p.savePromptToDebugFile("INITIAL PROMPT (Iterative)", initialPrompt)

	var generatedSQL string
	var testResult models.TestFileResult
	var allMessages []ConversationMessage
	var iterationResults []IterationResult
	totalTokens := 0

	// First iteration - send initial prompt
	fmt.Printf("  └─ Sending initial prompt to AI...\n")
	aiInitialStart := time.Now()
	response, err := p.sessionMgr.SendMessage(session.ID, initialPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to send initial message: %w", err)
	}

	// Save initial AI response to debug file if enabled
	p.savePromptToDebugFile("AI RESPONSE (Initial - Iterative)", response)

	generatedSQL = p.promptReader.ExtractSQLFromResponse(response)

	fmt.Printf("  └─ [%.3fs] Initial AI response received\n", time.Since(aiInitialStart).Seconds())

	for iteration := 1; iteration <= p.maxIterations; iteration++ {
		iterationStart := time.Now()

		// Test the current SQL
		testResult, err = p.testSQLString(generatedSQL)
		if err != nil {
			return nil, fmt.Errorf("failed to test SQL on iteration %d: %w", iteration, err)
		}

		// Check if we have success
		success := len(testResult.ParseErrors) == 0 && len(testResult.ExecutionErrors) == 0

		// Store iteration result
		iterationResult := IterationResult{
			Iteration:    iteration,
			TestResults:  testResult,
			Success:      success,
			GeneratedSQL: generatedSQL,
		}
		iterationResults = append(iterationResults, iterationResult)

		// Print iteration result in real-time
		p.printIterationResult(iteration, testResult)
		fmt.Printf("  └─ [%.3fs] Iteration %d completed\n", time.Since(iterationStart).Seconds(), iteration)

		if success {
			// Get final conversation history
			allMessages, _ = p.sessionMgr.GetConversationHistory(session.ID)

			result := &PipelineResult{
				SessionID:        session.ID,
				InitialPrompt:    initialPrompt,
				GeneratedSQL:     generatedSQL,
				TestResults:      testResult,
				Iterations:       iteration,
				IterationResults: iterationResults,
				Success:          true,
				Messages:         allMessages,
				TotalTime:        time.Since(start),
				TokensUsed:       totalTokens,
			}

			return result, nil
		}

		// If not successful and we have more iterations, send test results as feedback
		if iteration < p.maxIterations {
			aiStart := time.Now()
			testResultsString := p.formatTestResultsForPrompt(testResult)

			// Save feedback prompt to debug file if enabled
			p.savePromptToDebugFile(fmt.Sprintf("PROMPT (Iteration %d)", iteration+1), testResultsString)

			response, err = p.sessionMgr.SendMessage(session.ID, testResultsString)
			if err != nil {
				return nil, fmt.Errorf("failed to send feedback on iteration %d: %w", iteration, err)
			}

			generatedSQL = p.promptReader.ExtractSQLFromResponse(response)

			fmt.Printf("  └─ [%.3fs] AI response received for iteration %d\n", time.Since(aiStart).Seconds(), iteration+1)
		}
	}

	// Get final conversation history
	allMessages, _ = p.sessionMgr.GetConversationHistory(session.ID)

	// Max iterations reached without success
	result := &PipelineResult{
		SessionID:        session.ID,
		InitialPrompt:    initialPrompt,
		GeneratedSQL:     generatedSQL,
		TestResults:      testResult,
		Iterations:       p.maxIterations,
		IterationResults: iterationResults,
		Success:          false,
		Messages:         allMessages,
		TotalTime:        time.Since(start),
		TokensUsed:       totalTokens,
	}

	return result, nil
}

// testSQLString tests a SQL string using the existing testing infrastructure
func (p *Pipeline) testSQLString(sqlContent string) (models.TestFileResult, error) {
	// We'll use the string-based evaluation we already created
	result, err := p.evaluateSQLString(sqlContent, "generated")
	if err != nil {
		return models.TestFileResult{}, err
	}
	return result.FileResult, nil
}

// evaluateSQLString replicates the evaluation logic from cmd/sql-eval-string
func (p *Pipeline) evaluateSQLString(content string, filename string) (*EvaluationResult, error) {
	start := time.Now()

	fr := models.TestFileResult{
		Filename:        filename,
		ParseErrorCodes: make(map[string]int),
		ErrorCodes:      make(map[string]int),
		ErrorCategories: make(map[string]int),
	}

	statements, err := tools.ExtractStatementsFromString(content)
	if err != nil {
		return nil, fmt.Errorf("extract statements: %w", err)
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
		db, terminate, err := tools.GetDBWithIdentifier(true, p.uniqueID)
		if err != nil {
			return nil, fmt.Errorf("connect DB: %w", err)
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

	return &EvaluationResult{FileResult: fr}, nil
}

// formatTestResultsForPrompt formats test results into a string for the prompt
func (p *Pipeline) formatTestResultsForPrompt(fr models.TestFileResult) string {
	var results strings.Builder

	results.WriteString("The generated sql code has gone through some testing, here are the results:\n\n")

	if !p.shortPrompts {
		// Full format - include all summaries and statistics
		results.WriteString(fmt.Sprintf("Total statements: %d\n", fr.TotalStatements))
		results.WriteString(fmt.Sprintf("Successfully parsed: %d\n", fr.ParsedCount))
		results.WriteString(fmt.Sprintf("Parse errors: %d\n", len(fr.ParseErrors)))
		results.WriteString(fmt.Sprintf("Successfully executed: %d\n", fr.ExecutedCount))
		results.WriteString(fmt.Sprintf("Execution errors: %d\n", len(fr.ExecutionErrors)))

		if fr.TotalStatements > 0 {
			parseRate := float64(fr.ParsedCount) / float64(fr.TotalStatements) * 100
			execRate := 0.0
			if fr.ParsedCount > 0 {
				execRate = float64(fr.ExecutedCount) / float64(fr.ParsedCount) * 100
			}
			overall := float64(fr.ExecutedCount) / float64(fr.TotalStatements) * 100
			results.WriteString(fmt.Sprintf("Parse success rate: %.1f%%\n", parseRate))
			results.WriteString(fmt.Sprintf("Execution success rate (of parsed): %.1f%%\n", execRate))
			results.WriteString(fmt.Sprintf("Overall success rate: %.1f%%\n", overall))
		}

		if len(fr.ParseErrorCodes) > 0 {
			results.WriteString("\n")
			results.WriteString("Parse Error Summary:\n")
			for typ, cnt := range fr.ParseErrorCodes {
				desc := tools.GetParseErrorDescription(typ)
				results.WriteString(fmt.Sprintf("- %s: %d (%s)\n", typ, cnt, desc))
			}
		}

		if len(fr.ErrorCodes) > 0 {
			results.WriteString("\n")
			results.WriteString("Execution Error Code Summary:\n")
			for code, cnt := range fr.ErrorCodes {
				desc := tools.GetErrorCodeDescription(code)
				results.WriteString(fmt.Sprintf("- %s: %d (%s)\n", code, cnt, desc))
			}
		}

		if len(fr.ErrorCategories) > 0 {
			results.WriteString("\n")
			results.WriteString("Execution Error Categories:\n")
			for cat, cnt := range fr.ErrorCategories {
				desc := tools.GetErrorCategoryDescription(cat)
				results.WriteString(fmt.Sprintf("- %s: %d (%s)\n", cat, cnt, desc))
			}
		}
	}

	// Error details - for short prompts, only show the first line of each statement
	if len(fr.ParseErrors) > 0 {
		results.WriteString("\n")
		results.WriteString("Parse Errors:\n")
		for _, e := range fr.ParseErrorDetails {
			stmt := e.Statement
			if p.shortPrompts {
				// For short prompts, only show the first line
				lines := strings.Split(stmt, "\n")
				stmt = lines[0]
				if len(lines) > 1 {
					stmt += "..."
				}
			} else {
				// For full prompts, use the original truncation logic
				if len(stmt) > 200 {
					stmt = stmt[:100] + "..." + stmt[len(stmt)-100:]
				}
			}
			results.WriteString(fmt.Sprintf("- %s\n  Statement: %s\n", e.Error, stmt))
		}
	}

	if len(fr.ExecutionErrors) > 0 {
		results.WriteString("\n")
		results.WriteString("Execution Errors:\n")
		for i, e := range fr.ExecutionErrors {
			results.WriteString(fmt.Sprintf("%d. %s\n", i+1, e))
		}
	}

	// Add AI-specific recommendations
	recommendations := tools.GetAIRecommendations(fr)
	if len(recommendations) > 0 {
		results.WriteString("\n")
		results.WriteString("=== AI AGENT RECOMMENDATIONS ===\n")
		results.WriteString("ALWAYS: The response to this message should be a the entire sql code with fixes applied to it.\n")
		for _, rec := range recommendations {
			results.WriteString(rec + "\n")
		}
		results.WriteString("\n")
		results.WriteString("TIP: Focus on fixing parse errors first, as they prevent execution.\n")
		results.WriteString("TIP: Refer to Spanner SQL documentation for supported syntax and functions.\n")
	}

	return results.String()
}

// SaveResultToFile saves a pipeline result to a file
func (p *Pipeline) SaveResultToFile(result *PipelineResult, filename string) error {
	filePath := filepath.Join(p.basePath, "generated_sql", filename)

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the SQL content to file
	if err := os.WriteFile(filePath, []byte(result.GeneratedSQL), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// EvaluationResult wraps the file result for internal use
type EvaluationResult struct {
	FileResult models.TestFileResult
}

// LoadAccumulatedResults loads accumulated results from file, or creates a new structure if file doesn't exist
func LoadAccumulatedResults(filePath string) (*AccumulatedResults, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File doesn't exist, create new structure
		return &AccumulatedResults{
			Executions: make(map[string]*ExecutionMetrics),
			UpdatedAt:  time.Now(),
			Count:      0,
		}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read accumulated results file: %w", err)
	}

	var results AccumulatedResults
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal accumulated results: %w", err)
	}

	return &results, nil
}

// SaveAccumulatedResults saves accumulated results to file
func SaveAccumulatedResults(results *AccumulatedResults, filePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Update metadata
	results.UpdatedAt = time.Now()
	results.Count = len(results.Executions)

	// Marshal to JSON with pretty formatting for readability
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal accumulated results: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write accumulated results file: %w", err)
	}

	return nil
}

// AddResultToAccumulated adds a pipeline result to the accumulated results (lightweight version)
func (p *Pipeline) AddResultToAccumulated(result *PipelineResult, mode string) error {
	// Lock to ensure thread-safe access to pipeline_results.json
	pipelineResultsMutex.Lock()
	defer pipelineResultsMutex.Unlock()

	filePath := filepath.Join(p.basePath, "pipeline_results.json")

	// Load existing results
	accumulated, err := LoadAccumulatedResults(filePath)
	if err != nil {
		return fmt.Errorf("failed to load accumulated results: %w", err)
	}

	// Create lightweight metrics from the full result
	metrics := p.createExecutionMetrics(result, mode)

	// Add the new metrics using conversation ID as key
	accumulated.Executions[result.SessionID] = metrics

	// Save back to file
	if err := SaveAccumulatedResults(accumulated, filePath); err != nil {
		return fmt.Errorf("failed to save accumulated results: %w", err)
	}

	fmt.Printf("Results saved to accumulated file: %s (total executions: %d)\n", filePath, len(accumulated.Executions))
	return nil
}

// createExecutionMetrics extracts metrics from a PipelineResult including all iterations
func (p *Pipeline) createExecutionMetrics(result *PipelineResult, mode string) *ExecutionMetrics {
	var iterationResults []IterationMetrics

	// For single mode, we have only one iteration
	if mode == "single" {
		iteration := p.createIterationMetrics(1, result.TestResults)
		iterationResults = append(iterationResults, iteration)
	} else {
		// For iterative mode, process each iteration result
		for _, iterResult := range result.IterationResults {
			iteration := p.createIterationMetrics(iterResult.Iteration, iterResult.TestResults)
			iterationResults = append(iterationResults, iteration)
		}
	}

	return &ExecutionMetrics{
		ConversationID:   result.SessionID,
		Mode:             mode,
		Model:            p.model,
		Success:          result.Success,
		TotalIterations:  result.Iterations,
		ShortPrompts:     p.shortPrompts,
		IterationResults: iterationResults,
		Timestamp:        time.Now(),
	}
}

// createIterationMetrics creates metrics for a single iteration
func (p *Pipeline) createIterationMetrics(iterationNum int, testResults models.TestFileResult) IterationMetrics {
	parseRate := 0.0
	execRate := 0.0
	overall := 0.0

	if testResults.TotalStatements > 0 {
		parseRate = float64(testResults.ParsedCount) / float64(testResults.TotalStatements) * 100
		if testResults.ParsedCount > 0 {
			execRate = float64(testResults.ExecutedCount) / float64(testResults.ParsedCount) * 100
		}
		overall = float64(testResults.ExecutedCount) / float64(testResults.TotalStatements) * 100
	}

	success := len(testResults.ParseErrors) == 0 && len(testResults.ExecutionErrors) == 0

	return IterationMetrics{
		IterationNumber:      iterationNum,
		TotalStatements:      testResults.TotalStatements,
		SuccessfullyParsed:   testResults.ParsedCount,
		ParseErrors:          len(testResults.ParseErrors),
		Executed:             testResults.ExecutedCount,
		ExecutionErrors:      len(testResults.ExecutionErrors),
		ParseSuccessRate:     parseRate,
		ExecutionSuccessRate: execRate,
		OverallSuccessRate:   overall,
		Success:              success,
	}
}
