package tools

import (
	"strings"

	"github.com/cloudspannerecosystem/memefish"

	"sql-parser/models"
)

// ParseStatementsWithMemefish parses each statement using memefish and returns parse results.
func ParseStatementsWithMemefish(statements []string, filename string) []models.ParseResult {
	var results []models.ParseResult
	for _, stmt := range statements {
		pr := models.ParseResult{Statement: stmt}
		parsedStmt, err := memefish.ParseStatement(filename, stmt)
		if err != nil {
			pr.Parsed = false
			pr.Error = err
		} else {
			pr.Parsed = true
			pr.Type = GetStatementType(parsedStmt, stmt)
		}
		results = append(results, pr)
	}
	return results
}

// GetStatementType determines a basic statement type from the original statement text.
func GetStatementType(_ interface{}, originalStmt string) string {
	upper := strings.ToUpper(strings.TrimSpace(originalStmt))
	switch {
	case strings.HasPrefix(upper, "CREATE"):
		return "CREATE"
	case strings.HasPrefix(upper, "INSERT"):
		return "INSERT"
	case strings.HasPrefix(upper, "SELECT"):
		return "SELECT"
	case strings.HasPrefix(upper, "DROP"):
		return "DROP"
	case strings.HasPrefix(upper, "ALTER"):
		return "ALTER"
	case strings.HasPrefix(upper, "UPDATE"):
		return "UPDATE"
	case strings.HasPrefix(upper, "DELETE"):
		return "DELETE"
	default:
		return "OTHER"
	}
}

// CategorizeMemefishError categorizes memefish parsing errors for reporting.
func CategorizeMemefishError(errMsg string) string {
	lower := strings.ToLower(errMsg)
	switch {
	case strings.Contains(lower, "syntax error"):
		if strings.Contains(lower, "expecting") {
			return "Syntax Error: Missing Token"
		}
		return "Syntax Error: General"
	case strings.Contains(lower, "unexpected token"):
		return "Syntax Error: Unexpected Token"
	case strings.Contains(lower, "expecting"):
		return "Syntax Error: Expected Token"
	case strings.Contains(lower, "invalid"):
		return "Invalid Syntax"
	case strings.Contains(lower, "not supported"):
		return "Unsupported Feature"
	case strings.Contains(lower, "unknown"):
		return "Unknown Element"
	default:
		return "Parse Error: Other"
	}
}

// ExtractSpannerErrorCode extracts the gRPC/Spanner error code from an error string.
func ExtractSpannerErrorCode(errMsg string) string {
	if strings.Contains(errMsg, "rpc error: code = ") {
		start := strings.Index(errMsg, "rpc error: code = ") + len("rpc error: code = ")
		end := strings.Index(errMsg[start:], " desc = ")
		if end != -1 {
			return errMsg[start : start+end]
		}
	}
	if strings.Contains(errMsg, "spanner: code = ") {
		start := strings.Index(errMsg, "spanner: code = ") + len("spanner: code = ")
		if start < len(errMsg) && errMsg[start] == '"' {
			start++
			end := strings.Index(errMsg[start:], "\"")
			if end != -1 {
				return errMsg[start : start+end]
			}
		}
	}
	return ""
}

// CategorizeInvalidArgumentError provides finer categorization for InvalidArgument errors.
func CategorizeInvalidArgumentError(errMsg string) string {
	lower := strings.ToLower(errMsg)
	if strings.Contains(lower, "syntax error") {
		if strings.Contains(lower, "current_timestamp") {
			return "Syntax Error: CURRENT_TIMESTAMP"
		} else if strings.Contains(lower, "expecting '('") {
			return "Syntax Error: Missing Parentheses"
		} else if strings.Contains(lower, "expecting ')'") {
			return "Syntax Error: Missing Closing Parentheses"
		}
		return "Syntax Error: General"
	}
	if strings.Contains(lower, "expected type") && strings.Contains(lower, "found") {
		if strings.Contains(lower, "generate_uuid") {
			return "Type Mismatch: GENERATE_UUID on INT64"
		}
		return "Type Mismatch: General"
	}
	if strings.Contains(lower, "unsupported") {
		if strings.Contains(lower, "sequence kind") {
			return "Unsupported Feature: Sequence Kind"
		}
		return "Unsupported Feature: General"
	}
	if strings.Contains(lower, "missing") {
		if strings.Contains(lower, "sql security") {
			return "Missing Clause: SQL SECURITY"
		}
		return "Missing Clause: General"
	}
	if strings.Contains(lower, "function not found") {
		if strings.Contains(lower, "nextval") {
			return "Function Not Found: NEXTVAL"
		}
		return "Function Not Found: General"
	}
	if strings.Contains(lower, "sequence kind") && strings.Contains(lower, "not specified") {
		return "Identity Column: Missing Sequence Kind"
	}
	if strings.Contains(lower, "table not found") {
		return "Table Not Found (InvalidArgument)"
	}
	if strings.Contains(lower, "foreign key") {
		return "Foreign Key: Syntax Error"
	}
	if strings.Contains(lower, "default value") {
		return "Default Value: Parsing Error"
	}
	if strings.Contains(lower, "constraint") || strings.Contains(lower, "check") {
		return "Constraint: Unsupported"
	}
	if strings.Contains(lower, "definition of view") {
		return "View Definition: Error"
	}
	return "InvalidArgument: Other"
}

// GetErrorCodeDescription maps Spanner error codes to readable descriptions.
func GetErrorCodeDescription(code string) string {
	descriptions := map[string]string{
		"InvalidArgument":    "Invalid SQL syntax or unsupported features",
		"NotFound":           "Referenced table, column, or object not found",
		"FailedPrecondition": "Constraint violations or prerequisite not met",
		"AlreadyExists":      "Object already exists (duplicate creation)",
		"PermissionDenied":   "Insufficient permissions for operation",
		"Unimplemented":      "Feature not implemented in Spanner",
		"Internal":           "Internal Spanner error",
		"Unavailable":        "Service temporarily unavailable",
		"DeadlineExceeded":   "Operation timeout",
		"ResourceExhausted":  "Resource limits exceeded",
		"Cancelled":          "Operation was cancelled",
		"Unknown":            "Unknown error occurred",
	}
	if desc, ok := descriptions[code]; ok {
		return desc
	}
	return "Unknown error code"
}

// GetErrorCategoryDescription maps derived error categories to readable descriptions.
func GetErrorCategoryDescription(category string) string {
	descriptions := map[string]string{
		"Syntax Error: CURRENT_TIMESTAMP":           "CURRENT_TIMESTAMP() function call syntax not supported in Spanner",
		"Syntax Error: Missing Parentheses":         "SQL statement missing required opening parentheses",
		"Syntax Error: Missing Closing Parentheses": "SQL statement missing required closing parentheses",
		"Syntax Error: General":                     "General SQL syntax errors not matching specific patterns",
		"Type Mismatch: GENERATE_UUID on INT64":     "GENERATE_UUID() function used on INT64 columns instead of STRING",
		"Type Mismatch: General":                    "Data type mismatches between expected and provided types",
		"Unsupported Feature: Sequence Kind":        "Identity column sequence kind not specified or unsupported",
		"Unsupported Feature: General":              "General Spanner unsupported features",
		"Missing Clause: SQL SECURITY":              "VIEW definitions missing required SQL SECURITY clause",
		"Missing Clause: General":                   "SQL statements missing required clauses",
		"Function Not Found: NEXTVAL":               "NEXTVAL() function not available in Spanner",
		"Function Not Found: General":               "SQL functions not available in Spanner",
		"Identity Column: Missing Sequence Kind":    "Identity columns require explicit sequence kind specification",
		"Table Not Found (InvalidArgument)":         "Table references that result in InvalidArgument rather than NotFound",
		"Foreign Key: Syntax Error":                 "Foreign key constraint syntax errors",
		"Default Value: Parsing Error":              "Default value expressions that cannot be parsed",
		"Constraint: Unsupported":                   "CHECK constraints and other constraint types not supported",
		"View Definition: Error":                    "Errors in view definition syntax or structure",
		"NotFound":                                  "Referenced objects (tables, columns, etc.) not found",
		"FailedPrecondition":                        "Constraint violations or prerequisites not met",
		"AlreadyExists":                             "Attempting to create objects that already exist",
		"PermissionDenied":                          "Insufficient permissions for the operation",
		"Unimplemented":                             "Features not yet implemented in Spanner",
		"InvalidArgument: Other":                    "InvalidArgument errors not matching specific patterns",
	}
	if desc, ok := descriptions[category]; ok {
		return desc
	}
	return "No description available for this error category"
}

// GetParseErrorDescription maps memefish parse error types to readable descriptions.
func GetParseErrorDescription(errorType string) string {
	descriptions := map[string]string{
		"Syntax Error: Missing Token":    "SQL statements missing required tokens (parentheses, keywords, etc.)",
		"Syntax Error: General":          "General SQL syntax errors not matching specific patterns",
		"Syntax Error: Unexpected Token": "Unexpected tokens found where different syntax was expected",
		"Syntax Error: Expected Token":   "Missing expected tokens in SQL syntax",
		"Invalid Syntax":                 "SQL syntax that doesn't conform to Spanner SQL grammar",
		"Unsupported Feature":            "SQL features that are not supported by Spanner",
		"Unknown Element":                "Unknown SQL elements or identifiers",
		"Parse Error: Other":             "Other parsing errors not categorized above",
	}
	if desc, ok := descriptions[errorType]; ok {
		return desc
	}
	return "No description available for this parse error type"
}
