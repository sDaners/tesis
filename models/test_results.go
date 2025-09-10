package models

import (
	"time"
)

// ParseError holds detailed information about a parse error
type ParseError struct {
	Statement string
	Error     string
}

type ExecutionError struct {
	Statement   string
	Code        string
	Description string
}

// TestFileResult holds the results for a single SQL file test
type TestFileResult struct {
	Filename        string
	TotalStatements int
	// Parsing results
	ParsedCount       int
	ParseErrors       []string       // Deprecated: use ParseErrorDetails instead
	ParseErrorDetails []ParseError   // Detailed parse errors with statements
	ParseErrorCodes   map[string]int // error_type -> count
	// Statement type counts (from parsing)
	CreateStatements int
	InsertStatements int
	SelectStatements int
	DropStatements   int
	// Execution results
	ExecutedCount   int
	FailedCount     int
	ErrorRate       float64
	ExecutionTime   time.Duration
	ExecutionErrors []string
	ErrorCodes      map[string]int // error_code -> count
	ErrorCategories map[string]int // detailed_category -> count
}

// ParseResult holds the result of parsing a single statement
type ParseResult struct {
	Statement string
	Parsed    bool
	Error     error
	Type      string // CREATE, INSERT, SELECT, DROP, etc.
}

// AtomicStatementResult holds the results for a single SQL statement test
type AtomicStatementResult struct {
	Filename      string
	StatementNum  int
	Statement     string
	StatementType string
	ParseSuccess  bool
	ParseError    string
	ExecSuccess   bool
	ExecError     string
	ExecutionTime time.Duration
}

// AtomicFileResult aggregates results for all statements in a file
type AtomicFileResult struct {
	Filename         string
	TotalStatements  int
	ParsedCount      int
	ExecutedCount    int
	StatementResults []AtomicStatementResult
}

// AtomicTestSummary holds overall atomic test results
type AtomicTestSummary struct {
	TotalFiles      int
	TotalStatements int
	ParsedCount     int
	ExecutedCount   int
	FileResults     []AtomicFileResult
	ExecutionTime   time.Duration
}
