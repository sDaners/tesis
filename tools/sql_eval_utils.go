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
	case strings.Contains(lower, "expected token: ), but: ("):
		return "Syntax Error: PRIMARY/FOREIGN KEY Placement"
	case strings.Contains(lower, "expected token: (, but: <ident>") && strings.Contains(lower, "current_timestamp"):
		return "Syntax Error: CURRENT_TIMESTAMP Parentheses"
	case strings.Contains(lower, "expected token: (, but: <string>"):
		return "Syntax Error: String Literal Quotes"
	case strings.Contains(lower, "expected token"):
		return "Syntax Error: Expected Token"
	case strings.Contains(lower, "unexpected token"):
		return "Syntax Error: Unexpected Token"
	case strings.Contains(lower, "expecting"):
		return "Syntax Error: Missing Token"
	case strings.Contains(lower, "syntax error"):
		return "Syntax Error: General"
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
		"InvalidArgument":    "Invalid SQL syntax or unsupported features. FIX: Check for Spanner-specific syntax requirements (e.g., CURRENT_TIMESTAMP vs CURRENT_TIMESTAMP(), required clauses in views)",
		"NotFound":           "Referenced table, column, or object not found. FIX: Ensure all tables/columns exist before referencing them, or create them first in dependency order",
		"FailedPrecondition": "Constraint violations or prerequisite not met. FIX: Check for NOT NULL constraints, foreign key violations, or missing required data",
		"AlreadyExists":      "Object already exists (duplicate creation). FIX: Use CREATE OR REPLACE, or check if object exists before creating",
		"PermissionDenied":   "Insufficient permissions for operation. FIX: Verify user has required permissions for the database operation",
		"Unimplemented":      "Feature not implemented in Spanner. FIX: Use alternative Spanner-supported syntax or features",
		"Internal":           "Internal Spanner error. FIX: Retry the operation or contact support",
		"Unavailable":        "Service temporarily unavailable. FIX: Implement retry logic with exponential backoff",
		"DeadlineExceeded":   "Operation timeout. FIX: Optimize query performance or increase timeout settings",
		"ResourceExhausted":  "Resource limits exceeded. FIX: Reduce query complexity, add pagination, or increase quotas",
		"Cancelled":          "Operation was cancelled. FIX: Check for client-side cancellation or timeouts",
		"Unknown":            "Unknown error occurred. FIX: Check error details for more specific information",
	}
	if desc, ok := descriptions[code]; ok {
		return desc
	}
	return "Unknown error code"
}

// GetErrorCategoryDescription maps derived error categories to readable descriptions.
func GetErrorCategoryDescription(category string) string {
	descriptions := map[string]string{
		"Syntax Error: CURRENT_TIMESTAMP":           "Spanner requires DEFAULT values to be between parentheses. FIX: Use (CURRENT_TIMESTAMP()) for timestamp defaults",
		"Syntax Error: Missing Parentheses":         "SQL statement missing required opening parentheses. FIX: Add missing '(' where expected by parser",
		"Syntax Error: Missing Closing Parentheses": "SQL statement missing required closing parentheses. FIX: Add missing ')' to complete statement",
		"Syntax Error: General":                     "General SQL syntax errors not matching specific patterns. FIX: Check statement structure against Spanner SQL reference",
		"Type Mismatch: GENERATE_UUID on INT64":     "Spanner requires DEFAULT values to be between parentheses. FIX: Use (GENERATE_UUID()) for UUID columns",
		"Type Mismatch: General":                    "Data type mismatches between expected and provided types. FIX: Verify column types match inserted/compared values",
		"Unsupported Feature: Sequence Kind":        "The sequence was not properly defined. FIX: Sequence types should be avoided, use GENERATE_UUID() for primary keys",
		"Unsupported Feature: General":              "General Spanner unsupported features. FIX: Replace with Spanner-compatible alternatives",
		"Missing Clause: SQL SECURITY":              "VIEW definitions missing required SQL SECURITY clause. FIX: Add 'SQL SECURITY INVOKER' clause to view definition",
		"Missing Clause: General":                   "SQL statements missing required clauses. FIX: Add required clauses per Spanner SQL syntax",
		"Function Not Found: NEXTVAL":               "NEXTVAL() function not available in Spanner. FIX: Use GENERATE_UUID() for unique values or application-generated sequences",
		"Function Not Found: General":               "SQL functions not available in Spanner. FIX: Check Spanner function reference for supported alternatives",
		"Identity Column: Missing Sequence Kind":    "Identity columns require explicit sequence kind specification. FIX: Sequence types should be avoided, use GENERATE_UUID() for primary keys",
		"Table Not Found (InvalidArgument)":         "Table references that result in InvalidArgument rather than NotFound. FIX: There is likely a error creating the referenced table, so ignore this error",
		"Foreign Key: Syntax Error":                 "Foreign key constraint syntax errors. FIX: Use CONSTRAINT name FOREIGN KEY (col) REFERENCES table(col) syntax",
		"Default Value: Parsing Error":              "Default value expressions that cannot be parsed. FIX: Use simple literals or supported functions like CURRENT_TIMESTAMP",
		"View Definition: Error":                    "Errors in view definition syntax or structure. FIX: Ensure view uses SELECT statement and includes SQL SECURITY clause",
		"NotFound":                                  "Referenced objects (tables, columns, etc.) not found. FIX: There is likely a error creating the referenced table, so ignore this error",
		"FailedPrecondition":                        "Constraint violations or prerequisites not met. FIX: Ensure data meets NOT NULL, foreign key, and other constraints",
		"AlreadyExists":                             "Attempting to create objects that already exist. FIX: Use CREATE OR REPLACE or check existence first",
		"PermissionDenied":                          "Insufficient permissions for the operation. FIX: Grant necessary permissions or use appropriate service account",
		"Unimplemented":                             "Features not yet implemented in Spanner. FIX: Check Spanner roadmap or use alternative approaches",
		"InvalidArgument: Other":                    "InvalidArgument errors not matching specific patterns. FIX: Review error message details for specific syntax issues",
	}
	if desc, ok := descriptions[category]; ok {
		return desc
	}
	return "No description available for this error category"
}

// GetParseErrorDescription maps memefish parse error types to readable descriptions.
func GetParseErrorDescription(errorType string) string {
	descriptions := map[string]string{
		"Syntax Error: PRIMARY/FOREIGN KEY Placement": "PRIMARY KEY constraints MUST be placed OUTSIDE the column definition parentheses in Spanner. FIX: Move PRIMARY KEY clause to after the closing parenthesis of column definitions. CORRECT: ') PRIMARY KEY (column_name);' WRONG: 'PRIMARY KEY (column_name)' inside column list. This is a critical Spanner-specific syntax requirement.",
		"Syntax Error: CURRENT_TIMESTAMP Parentheses": "CURRENT_TIMESTAMP function call in DEFAULT clause needs proper parentheses. FIX: Change 'DEFAULT CURRENT_TIMESTAMP()' to 'DEFAULT (CURRENT_TIMESTAMP())'. In Spanner, DEFAULT values must be wrapped in parentheses.",
		"Syntax Error: String Literal Quotes":         "String literals in SQL statements are using incorrect quote types. FIX: Use single quotes for string literals instead of double quotes. Change \"ACTIVE\" to 'ACTIVE'. Spanner requires single quotes for string constants.",
		"Syntax Error: Missing Token":                 "SQL statements missing required tokens (parentheses, keywords, etc.). FIX: Add missing syntax elements as indicated by parser",
		"Syntax Error: General":                       "General SQL syntax errors not matching specific patterns. FIX: Review statement structure, check for typos and syntax compliance with Spanner SQL",
		"Syntax Error: Unexpected Token":              "Unexpected tokens found where different syntax was expected. FIX: Remove or relocate unexpected elements to correct positions",
		"Syntax Error: Expected Token":                "Missing expected tokens in SQL syntax. FIX: Add required keywords, punctuation, or identifiers where expected. After DEFAULT remember to wrap the value in parentheses",
		"Invalid Syntax":                              "SQL syntax that doesn't conform to Spanner SQL grammar. FIX: Rewrite using valid Spanner SQL syntax patterns",
		"Unsupported Feature":                         "SQL features that are not supported by Spanner. FIX: Replace with Spanner-compatible alternatives (e.g., use ARRAY instead of arrays)",
		"Unknown Element":                             "Unknown SQL elements or identifiers. FIX: Check spelling of keywords, functions, and identifiers against Spanner documentation",
		"Parse Error: Other":                          "Other parsing errors not categorized above. FIX: Review error message for specific guidance",
	}
	if desc, ok := descriptions[errorType]; ok {
		return desc
	}
	return "No description available for this parse error type"
}

// GetAIRecommendations generates AI-specific recommendations based on error patterns
func GetAIRecommendations(fr models.TestFileResult) []string {
	var recommendations []string

	// Parse error recommendations
	if len(fr.ParseErrors) > 0 {
		recommendations = append(recommendations, "PARSE ERROR PATTERNS DETECTED:")

		if fr.ParseErrorCodes["Syntax Error: General"] > 0 {
			recommendations = append(recommendations,
				"• Multiple syntax errors found. Common issues:",
				"  - CURRENT_TIMESTAMP() should be (CURRENT_TIMESTAMP()). Default values must be between parentheses",
				"  - Views require 'SQL SECURITY INVOKER' clause after the table name",
				"  - RETURNING should be replaced with 'THEN RETURN'",
			)
		}

		if fr.ParseErrorCodes["Syntax Error: Expected Token"] > 0 || fr.ParseErrorCodes["Syntax Error: Missing Token"] > 0 {
			recommendations = append(recommendations,
				"• Missing or unexpected tokens detected:",
				"  - Check parentheses, commas, and keyword placement",
				"  - Ensure proper statement termination with semicolons",
				"  - Verify correct positioning of PRIMARY KEY constraints")
		}

		if fr.ParseErrorCodes["Syntax Error: PRIMARY/FOREIGN KEY Placement"] > 0 {
			recommendations = append(recommendations,
				"• PRIMARY KEY or FOREIGN KEY placement issues - CRITICAL SPANNER SYNTAX:",
				"  - PRIMARY KEY must be OUTSIDE column definitions: ') PRIMARY KEY (column_name);'",
				"  - WRONG: 'PRIMARY KEY (column_name)' inside the column list",
				"  - CORRECT: Close column list with ), then add PRIMARY KEY (column_name);",
				"  - FOREIGN KEY constraints go INSIDE column definitions, before closing )")
		}

		if fr.ParseErrorCodes["Syntax Error: CURRENT_TIMESTAMP Parentheses"] > 0 {
			recommendations = append(recommendations,
				"• CURRENT_TIMESTAMP parentheses issues:",
				"  - Change 'DEFAULT CURRENT_TIMESTAMP()' to 'DEFAULT (CURRENT_TIMESTAMP())'",
				"  - Spanner requires DEFAULT values to be wrapped in parentheses",
				"  - This is a very common Spanner-specific syntax requirement")
		}

		if fr.ParseErrorCodes["Syntax Error: String Literal Quotes"] > 0 {
			recommendations = append(recommendations,
				"• String literal quote issues:",
				"  - Use single quotes for strings: 'ACTIVE' not \"ACTIVE\"",
				"  - Change all double-quoted strings to single quotes in CHECK constraints",
				"  - Spanner requires single quotes for string literals")
		}
	}

	// Execution error recommendations
	if len(fr.ExecutionErrors) > 0 {
		hasParseErrors := len(fr.ParseErrors) > 0
		executionRecommendationsAdded := false

		// Skip NotFound errors if there are parse errors (they're likely caused by failed table creation)
		if fr.ErrorCodes["NotFound"] > 0 && !hasParseErrors {
			if !executionRecommendationsAdded {
				recommendations = append(recommendations, "EXECUTION ERROR PATTERNS DETECTED:")
				executionRecommendationsAdded = true
			}
			recommendations = append(recommendations,
				"• Table/column not found errors:",
				"  - The table may not be found because it's creation failed earlier, in that case ignore this error",
				"  - Create tables in dependency order (referenced tables first)",
				"  - Verify table and column names match exactly",
				"  - Check for typos in table/column references")
		}

		if fr.ErrorCodes["FailedPrecondition"] > 0 {
			if !executionRecommendationsAdded {
				recommendations = append(recommendations, "EXECUTION ERROR PATTERNS DETECTED:")
				executionRecommendationsAdded = true
			}
			recommendations = append(recommendations,
				"• Constraint violation errors:",
				"  - Ensure primary keys are auto generated. Example: `key STRING(36) DEFAULT (GENERATE_UUID())`",
				"  - Ensure NOT NULL columns have values in INSERT statements",
				"  - Verify foreign key relationships exist before inserting",
				"  - Check data types match column definitions")
		}

		if fr.ErrorCodes["InvalidArgument"] > 0 {
			if !executionRecommendationsAdded {
				recommendations = append(recommendations, "EXECUTION ERROR PATTERNS DETECTED:")
				executionRecommendationsAdded = true
			}
			recommendations = append(recommendations,
				"• Invalid argument errors (often syntax-related):",
				"  - Use Spanner-specific SQL syntax and functions",
				"  - Replace unsupported features with Spanner alternatives",
				"  - Check function signatures and parameter types")
		}

		// Skip Table Not Found (InvalidArgument) errors if there are parse errors
		if fr.ErrorCategories["Table Not Found (InvalidArgument)"] > 0 && !hasParseErrors {
			if !executionRecommendationsAdded {
				recommendations = append(recommendations, "EXECUTION ERROR PATTERNS DETECTED:")
				executionRecommendationsAdded = true
			}
			recommendations = append(recommendations,
				"• Table references causing InvalidArgument:",
				"  - The table may not be found because it's creation failed earlier, in that case ignore this error",
				"  - This usually indicates table creation failed earlier",
				"  - Fix table creation statements first, then retry queries")
		}
	}

	// Spanner-specific best practices
	recommendations = append(recommendations, "SPANNER SQL BEST PRACTICES FOR AI AGENTS:")
	recommendations = append(recommendations,
		"• CRITICAL: PRIMARY KEY must be OUTSIDE column definitions: ') PRIMARY KEY (column_name);'",
		"• Use GENERATE_UUID() for primary keys instead of auto-increment",
		"• Create tables before referencing them in foreign keys or queries",
		"• Use STRING(36) with generated UUIDs for primary keys",
		"• Include SQL SECURITY INVOKER in all view definitions",
		"• Use ARRAY<TYPE> for array columns, not array syntax from other databases",
	)

	return recommendations
}
