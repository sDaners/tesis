package tools

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sql-parser/models"

	"github.com/google/uuid"
)

// AllureLabel represents an Allure label
type AllureLabel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AllureParameter represents an Allure parameter
type AllureParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AllureAttachment represents an Allure attachment
type AllureAttachment struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

// AllureStatusDetails represents status details for failed tests
type AllureStatusDetails struct {
	Known   bool   `json:"known"`
	Muted   bool   `json:"muted"`
	Flaky   bool   `json:"flaky"`
	Message string `json:"message,omitempty"`
	Trace   string `json:"trace,omitempty"`
}

// AllureStep represents a test step
type AllureStep struct {
	Name          string               `json:"name"`
	Status        string               `json:"status"`
	StatusDetails *AllureStatusDetails `json:"statusDetails,omitempty"`
	Stage         string               `json:"stage"`
	Description   string               `json:"description,omitempty"`
	Start         int64                `json:"start"`
	Stop          int64                `json:"stop"`
	UUID          string               `json:"uuid"`
	Attachments   []AllureAttachment   `json:"attachments,omitempty"`
	Parameters    []AllureParameter    `json:"parameters,omitempty"`
}

// AllureResult represents a test result in Allure format
type AllureResult struct {
	UUID          string               `json:"uuid"`
	HistoryId     string               `json:"historyId"`
	TestCaseId    string               `json:"testCaseId"`
	RerunOf       string               `json:"rerunOf,omitempty"`
	FullName      string               `json:"fullName"`
	Name          string               `json:"name"`
	Status        string               `json:"status"`
	StatusDetails *AllureStatusDetails `json:"statusDetails,omitempty"`
	Stage         string               `json:"stage"`
	Description   string               `json:"description,omitempty"`
	Start         int64                `json:"start"`
	Stop          int64                `json:"stop"`
	Labels        []AllureLabel        `json:"labels"`
	Parameters    []AllureParameter    `json:"parameters,omitempty"`
	Attachments   []AllureAttachment   `json:"attachments,omitempty"`
	Steps         []AllureStep         `json:"steps,omitempty"`
}

// AllureReporter handles the generation of Allure reports from SQL test results
type AllureReporter struct {
	outputDir string
}

// NewAllureReporter creates a new Allure reporter instance
func NewAllureReporter(outputDir string) *AllureReporter {
	if outputDir == "" {
		outputDir = "allure-results"
	}
	return &AllureReporter{outputDir: outputDir}
}

// GenerateAllureReport creates Allure reports from TestFileResult data
func (r *AllureReporter) GenerateAllureReport(results []models.TestFileResult) error {
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create allure output directory: %w", err)
	}

	// Process each file's results
	for _, fileResult := range results {
		if err := r.generateFileTestResult(fileResult); err != nil {
			return fmt.Errorf("failed to generate allure results for file %s: %w", fileResult.Filename, err)
		}
	}

	return nil
}

// GenerateAtomicAllureReport creates Allure reports from AtomicTestSummary data
func (r *AllureReporter) GenerateAtomicAllureReport(summary models.AtomicTestSummary) error {
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Generate test results for each file
	for _, fileResult := range summary.FileResults {
		if err := r.generateAtomicFileTestResults(fileResult); err != nil {
			return fmt.Errorf("generating file test results for %s: %w", fileResult.Filename, err)
		}
	}

	return nil
}

// generateFileTestResult creates an Allure test result for a single SQL file
func (r *AllureReporter) generateFileTestResult(fileResult models.TestFileResult) error {
	testUUID := uuid.New().String()
	now := time.Now()

	// Create test case ID from file name
	testCaseId := fmt.Sprintf("%x", md5.Sum([]byte(fileResult.Filename)))

	// Determine test status
	status := "passed"
	var statusDetails *AllureStatusDetails

	if len(fileResult.ParseErrors) > 0 {
		status = "broken"
		statusDetails = &AllureStatusDetails{
			Message: fmt.Sprintf("Parse errors: %d", len(fileResult.ParseErrors)),
			Trace:   strings.Join(fileResult.ParseErrors, "\n"),
		}
	} else if len(fileResult.ExecutionErrors) > 0 {
		status = "failed"
		statusDetails = &AllureStatusDetails{
			Message: fmt.Sprintf("Execution errors: %d", len(fileResult.ExecutionErrors)),
			Trace:   strings.Join(fileResult.ExecutionErrors, "\n"),
		}
	}

	// Create Allure result
	result := AllureResult{
		UUID:          testUUID,
		HistoryId:     testCaseId,
		TestCaseId:    testCaseId,
		FullName:      fmt.Sprintf("sql.compatibility.%s", strings.ReplaceAll(fileResult.Filename, ".", "_")),
		Name:          fmt.Sprintf("SQL Compatibility Test: %s", fileResult.Filename),
		Status:        status,
		StatusDetails: statusDetails,
		Stage:         "finished",
		Description:   fmt.Sprintf("SQL compatibility test for file: %s", fileResult.Filename),
		Start:         now.UnixMilli(),
		Stop:          now.Add(fileResult.ExecutionTime).UnixMilli(),
		Labels: []AllureLabel{
			{Name: "suite", Value: "SQL Compatibility Tests"},
			{Name: "feature", Value: "SQL File Processing"},
			{Name: "story", Value: filepath.Base(fileResult.Filename)},
			{Name: "severity", Value: getSeverityFromResults(fileResult)},
			{Name: "tag", Value: "sql-compatibility"},
			{Name: "framework", Value: "go"},
			{Name: "language", Value: "go"},
		},
		Parameters: []AllureParameter{
			{Name: "filename", Value: fileResult.Filename},
			{Name: "total_statements", Value: fmt.Sprintf("%d", fileResult.TotalStatements)},
			{Name: "parsed_count", Value: fmt.Sprintf("%d", fileResult.ParsedCount)},
			{Name: "executed_count", Value: fmt.Sprintf("%d", fileResult.ExecutedCount)},
			{Name: "execution_time", Value: fileResult.ExecutionTime.String()},
		},
	}

	// Add steps for detailed breakdown
	r.addFileSteps(&result, fileResult)

	// Save to JSON file
	filename := fmt.Sprintf("%s-result.json", testUUID)
	return r.saveAllureResult(filename, result)
}

// generateAtomicFileTestResults creates Allure test results for atomic SQL tests
func (r *AllureReporter) generateAtomicFileTestResults(fileResult models.AtomicFileResult) error {
	// Create individual test results for each statement
	for _, stmtResult := range fileResult.StatementResults {
		testUUID := uuid.New().String()
		now := time.Now()

		// Create test case ID from file name and statement number
		testCaseId := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s_stmt_%d", fileResult.Filename, stmtResult.StatementNum))))

		// Determine test status
		status := "passed"
		var statusDetails *AllureStatusDetails

		if !stmtResult.ParseSuccess {
			status = "broken"
			statusDetails = &AllureStatusDetails{
				Message: "Parse Error",
				Trace:   stmtResult.ParseError,
			}
		} else if !stmtResult.ExecSuccess {
			status = "failed"
			statusDetails = &AllureStatusDetails{
				Message: "Execution Error",
				Trace:   stmtResult.ExecError,
			}
		}

		// Truncate long statements for display
		stmt := stmtResult.Statement
		if len(stmt) > 100 {
			stmt = stmt[:97] + "..."
		}

		// Create Allure result
		result := AllureResult{
			UUID:          testUUID,
			HistoryId:     testCaseId,
			TestCaseId:    testCaseId,
			FullName:      fmt.Sprintf("sql.atomic.%s.stmt_%d", strings.ReplaceAll(fileResult.Filename, ".", "_"), stmtResult.StatementNum),
			Name:          fmt.Sprintf("%s - Statement %d (%s)", fileResult.Filename, stmtResult.StatementNum, stmtResult.StatementType),
			Status:        status,
			StatusDetails: statusDetails,
			Stage:         "finished",
			Description:   fmt.Sprintf("Statement: %s", stmt),
			Start:         now.UnixMilli(),
			Stop:          now.Add(stmtResult.ExecutionTime).UnixMilli(),
			Labels: []AllureLabel{
				{Name: "suite", Value: "Atomic SQL Tests"},
				{Name: "feature", Value: stmtResult.StatementType + " Statements"},
				{Name: "story", Value: filepath.Base(fileResult.Filename)},
				{Name: "severity", Value: getSeverityFromAtomicResult(stmtResult)},
				{Name: "tag", Value: "atomic-sql"},
				{Name: "tag", Value: strings.ToLower(stmtResult.StatementType)},
				{Name: "framework", Value: "go"},
				{Name: "language", Value: "go"},
			},
			Parameters: []AllureParameter{
				{Name: "filename", Value: fileResult.Filename},
				{Name: "statement_number", Value: fmt.Sprintf("%d", stmtResult.StatementNum)},
				{Name: "statement_type", Value: stmtResult.StatementType},
				{Name: "execution_time", Value: stmtResult.ExecutionTime.String()},
			},
		}

		// Add step for statement execution
		stepUUID := uuid.New().String()
		step := AllureStep{
			Name:          fmt.Sprintf("Execute %s Statement", stmtResult.StatementType),
			Status:        status,
			StatusDetails: statusDetails,
			Stage:         "finished",
			Description:   fmt.Sprintf("Execute SQL statement: %s", stmt),
			Start:         now.UnixMilli(),
			Stop:          now.Add(stmtResult.ExecutionTime).UnixMilli(),
			UUID:          stepUUID,
		}

		// Add statement as attachment
		if len(stmtResult.Statement) > 0 {
			attachmentName := fmt.Sprintf("statement_%d.sql", stmtResult.StatementNum)
			attachmentPath := filepath.Join(r.outputDir, attachmentName)

			if err := os.WriteFile(attachmentPath, []byte(stmtResult.Statement), 0644); err == nil {
				step.Attachments = append(step.Attachments, AllureAttachment{
					Name:   "SQL Statement",
					Source: attachmentName,
					Type:   "text/plain",
				})
			}
		}

		result.Steps = append(result.Steps, step)

		// Save to JSON file
		filename := fmt.Sprintf("%s-result.json", testUUID)
		if err := r.saveAllureResult(filename, result); err != nil {
			return err
		}
	}

	return nil
}

// addFileSteps adds detailed steps to a file test result
func (r *AllureReporter) addFileSteps(result *AllureResult, fileResult models.TestFileResult) {
	now := time.Now()

	// Parsing step
	parseStepUUID := uuid.New().String()
	parseStatus := "passed"
	var parseStatusDetails *AllureStatusDetails

	if len(fileResult.ParseErrors) > 0 {
		parseStatus = "failed"
		parseStatusDetails = &AllureStatusDetails{
			Message: fmt.Sprintf("Parse errors: %d", len(fileResult.ParseErrors)),
		}

		// Add parse errors as attachment
		if len(fileResult.ParseErrors) > 0 {
			attachmentName := fmt.Sprintf("parse_errors_%s.txt", parseStepUUID)
			attachmentPath := filepath.Join(r.outputDir, attachmentName)
			parseErrorsText := strings.Join(fileResult.ParseErrors, "\n")

			if err := os.WriteFile(attachmentPath, []byte(parseErrorsText), 0644); err == nil {
				parseStep := AllureStep{
					Name:          "Parse SQL Statements",
					Status:        parseStatus,
					StatusDetails: parseStatusDetails,
					Stage:         "finished",
					Start:         now.UnixMilli(),
					Stop:          now.Add(time.Second).UnixMilli(),
					UUID:          parseStepUUID,
					Attachments: []AllureAttachment{
						{
							Name:   "Parse Errors",
							Source: attachmentName,
							Type:   "text/plain",
						},
					},
				}
				result.Steps = append(result.Steps, parseStep)
			}
		}
	} else {
		parseStep := AllureStep{
			Name:   "Parse SQL Statements",
			Status: parseStatus,
			Stage:  "finished",
			Start:  now.UnixMilli(),
			Stop:   now.Add(time.Second).UnixMilli(),
			UUID:   parseStepUUID,
		}
		result.Steps = append(result.Steps, parseStep)
	}

	// Execution step
	execStepUUID := uuid.New().String()
	execStatus := "passed"
	var execStatusDetails *AllureStatusDetails

	if len(fileResult.ExecutionErrors) > 0 {
		execStatus = "failed"
		execStatusDetails = &AllureStatusDetails{
			Message: fmt.Sprintf("Execution errors: %d", len(fileResult.ExecutionErrors)),
		}

		// Add execution errors as attachment
		if len(fileResult.ExecutionErrors) > 0 {
			attachmentName := fmt.Sprintf("execution_errors_%s.txt", execStepUUID)
			attachmentPath := filepath.Join(r.outputDir, attachmentName)
			execErrorsText := strings.Join(fileResult.ExecutionErrors, "\n")

			if err := os.WriteFile(attachmentPath, []byte(execErrorsText), 0644); err == nil {
				execStep := AllureStep{
					Name:          "Execute SQL Statements",
					Status:        execStatus,
					StatusDetails: execStatusDetails,
					Stage:         "finished",
					Start:         now.Add(time.Second).UnixMilli(),
					Stop:          now.Add(2 * time.Second).UnixMilli(),
					UUID:          execStepUUID,
					Attachments: []AllureAttachment{
						{
							Name:   "Execution Errors",
							Source: attachmentName,
							Type:   "text/plain",
						},
					},
				}
				result.Steps = append(result.Steps, execStep)
			}
		}
	} else {
		execStep := AllureStep{
			Name:   "Execute SQL Statements",
			Status: execStatus,
			Stage:  "finished",
			Start:  now.Add(time.Second).UnixMilli(),
			Stop:   now.Add(2 * time.Second).UnixMilli(),
			UUID:   execStepUUID,
		}
		result.Steps = append(result.Steps, execStep)
	}
}

// saveAllureResult saves an Allure result to a JSON file
func (r *AllureReporter) saveAllureResult(filename string, result AllureResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal allure result: %w", err)
	}

	filePath := filepath.Join(r.outputDir, filename)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write allure result file: %w", err)
	}

	return nil
}

// getSeverityFromResults determines test severity based on results
func getSeverityFromResults(fileResult models.TestFileResult) string {
	if len(fileResult.ParseErrors) > 0 {
		return "critical"
	} else if len(fileResult.ExecutionErrors) > 0 {
		return "major"
	} else if fileResult.ExecutedCount < fileResult.TotalStatements {
		return "minor"
	}
	return "normal"
}

// getSeverityFromAtomicResult determines test severity for atomic results
func getSeverityFromAtomicResult(stmtResult models.AtomicStatementResult) string {
	if !stmtResult.ParseSuccess {
		return "critical"
	} else if !stmtResult.ExecSuccess {
		return "major"
	}
	return "normal"
}
