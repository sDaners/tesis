package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sql-parser/models"
	"sql-parser/tools"
)

// Pipeline orchestrates the complete workflow from prompt to tested SQL
type Pipeline struct {
	client        *OpenAIClient
	sessionMgr    *SessionManager
	promptReader  *PromptReader
	basePath      string
	maxIterations int
}

// NewPipeline creates a new pipeline instance
func NewPipeline(basePath string, maxIterations int) (*Pipeline, error) {
	client, err := NewOpenAIClientFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	sessionMgr := NewSessionManager(client)
	promptReader := NewPromptReader(basePath)

	return &Pipeline{
		client:        client,
		sessionMgr:    sessionMgr,
		promptReader:  promptReader,
		basePath:      basePath,
		maxIterations: maxIterations,
	}, nil
}

// RunSingleShot runs a single-shot prompt without iteration
func (p *Pipeline) RunSingleShot() (*PipelineResult, error) {
	start := time.Now()

	// Read the initial prompt
	initialPrompt, err := p.promptReader.FormatInitialPrompt()
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt: %w", err)
	}

	// Send to OpenAI
	response, err := p.client.SendSingleMessage(initialPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to send message to OpenAI: %w", err)
	}

	// Extract SQL from response
	generatedSQL := p.promptReader.ExtractSQLFromResponse(response)

	// Test the generated SQL
	testResult, err := p.testSQLString(generatedSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to test SQL: %w", err)
	}

	success := len(testResult.ParseErrors) == 0 && len(testResult.ExecutionErrors) == 0

	result := &PipelineResult{
		SessionID:     "single-shot",
		InitialPrompt: initialPrompt,
		GeneratedSQL:  generatedSQL,
		TestResults:   testResult,
		Iterations:    1,
		Success:       success,
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
	initialPrompt, err := p.promptReader.FormatInitialPrompt()
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt: %w", err)
	}

	var generatedSQL string
	var testResult models.TestFileResult
	var allMessages []ConversationMessage
	totalTokens := 0

	// First iteration - send initial prompt
	response, err := p.sessionMgr.SendMessage(session.ID, initialPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to send initial message: %w", err)
	}

	generatedSQL = p.promptReader.ExtractSQLFromResponse(response)

	for iteration := 1; iteration <= p.maxIterations; iteration++ {
		// Test the current SQL
		testResult, err = p.testSQLString(generatedSQL)
		if err != nil {
			return nil, fmt.Errorf("failed to test SQL on iteration %d: %w", iteration, err)
		}

		// Check if we have success
		success := len(testResult.ParseErrors) == 0 && len(testResult.ExecutionErrors) == 0

		if success {
			// Get final conversation history
			allMessages, _ = p.sessionMgr.GetConversationHistory(session.ID)

			result := &PipelineResult{
				SessionID:     session.ID,
				InitialPrompt: initialPrompt,
				GeneratedSQL:  generatedSQL,
				TestResults:   testResult,
				Iterations:    iteration,
				Success:       true,
				Messages:      allMessages,
				TotalTime:     time.Since(start),
				TokensUsed:    totalTokens,
			}

			return result, nil
		}

		// If not successful and we have more iterations, send test results as feedback
		if iteration < p.maxIterations {
			testResultsString := p.formatTestResultsForPrompt(testResult)
			feedbackPrompt := p.promptReader.FormatTestResultsPrompt(generatedSQL, testResultsString)

			response, err = p.sessionMgr.SendMessage(session.ID, feedbackPrompt)
			if err != nil {
				return nil, fmt.Errorf("failed to send feedback on iteration %d: %w", iteration, err)
			}

			generatedSQL = p.promptReader.ExtractSQLFromResponse(response)
		}
	}

	// Get final conversation history
	allMessages, _ = p.sessionMgr.GetConversationHistory(session.ID)

	// Max iterations reached without success
	result := &PipelineResult{
		SessionID:     session.ID,
		InitialPrompt: initialPrompt,
		GeneratedSQL:  generatedSQL,
		TestResults:   testResult,
		Iterations:    p.maxIterations,
		Success:       false,
		Messages:      allMessages,
		TotalTime:     time.Since(start),
		TokensUsed:    totalTokens,
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
		// This would need to be implemented to actually test against Spanner
		// For now, we'll use a mock or skip execution testing
		// TODO: Integrate with actual Spanner testing
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

	results.WriteString(fmt.Sprintf("Total statements: %d\n", fr.TotalStatements))
	results.WriteString(fmt.Sprintf("Successfully parsed: %d\n", fr.ParsedCount))
	results.WriteString(fmt.Sprintf("Parse errors: %d\n", len(fr.ParseErrors)))
	results.WriteString(fmt.Sprintf("Successfully executed: %d\n", fr.ExecutedCount))
	results.WriteString(fmt.Sprintf("Execution errors: %d\n", len(fr.ExecutionErrors)))

	if len(fr.ParseErrors) > 0 {
		results.WriteString("\nParse Errors:\n")
		for _, e := range fr.ParseErrorDetails {
			results.WriteString(fmt.Sprintf("- %s\n  Statement: %s\n", e.Error, e.Statement))
		}
	}

	if len(fr.ExecutionErrors) > 0 {
		results.WriteString("\nExecution Errors:\n")
		for i, e := range fr.ExecutionErrors {
			results.WriteString(fmt.Sprintf("%d. %s\n", i+1, e))
		}
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
