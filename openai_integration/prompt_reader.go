package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PromptReader handles reading and formatting prompts
type PromptReader struct {
	basePath string
}

// NewPromptReader creates a new prompt reader
func NewPromptReader(basePath string) *PromptReader {
	return &PromptReader{
		basePath: basePath,
	}
}

// ReadPromptFile reads the contents of prompt.txt
func (pr *PromptReader) ReadPromptFile() (string, error) {
	promptPath := filepath.Join(pr.basePath, "prompt.txt")

	content, err := os.ReadFile(promptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", promptPath, err)
	}

	return string(content), nil
}

// FormatInitialPrompt creates the initial prompt to send to OpenAI
func (pr *PromptReader) FormatInitialPrompt() (string, error) {
	promptContent, err := pr.ReadPromptFile()
	if err != nil {
		return "", err
	}

	// The prompt.txt already contains the instructions and the PostgreSQL code to translate
	// We can send it as-is since it already has the proper format
	return promptContent, nil
}

// FormatTestResultsPrompt creates a prompt that includes test results for iteration
func (pr *PromptReader) FormatTestResultsPrompt(generatedSQL string, testResults string) string {
	var prompt strings.Builder

	prompt.WriteString("The SQL code you generated has been tested. Here are the results:\n\n")
	prompt.WriteString("=== TEST RESULTS ===\n")
	prompt.WriteString(testResults)
	prompt.WriteString("\n=== GENERATED SQL CODE ===\n")
	prompt.WriteString(generatedSQL)
	prompt.WriteString("\n\n=== INSTRUCTIONS ===\n")
	prompt.WriteString("Please analyze the test results and fix any errors in the SQL code. ")
	prompt.WriteString("Focus on parse errors first as they prevent execution. ")
	prompt.WriteString("Return the complete corrected SQL code.\n")
	prompt.WriteString("ALWAYS: The response to this message should be the entire sql code with fixes applied to it.")

	return prompt.String()
}

// ExtractSQLFromResponse attempts to extract SQL code from an AI response
func (pr *PromptReader) ExtractSQLFromResponse(response string) string {
	// Try to find SQL code blocks first
	if strings.Contains(response, "```sql") {
		start := strings.Index(response, "```sql")
		if start != -1 {
			start += 6 // Skip "```sql"
			end := strings.Index(response[start:], "```")
			if end != -1 {
				return strings.TrimSpace(response[start : start+end])
			}
		}
	}

	// Try to find generic code blocks
	if strings.Contains(response, "```") {
		start := strings.Index(response, "```")
		if start != -1 {
			// Skip the first ```
			start = strings.Index(response[start+3:], "\n")
			if start != -1 {
				start += 4 // Skip "```\n"
				end := strings.Index(response[start:], "```")
				if end != -1 {
					return strings.TrimSpace(response[start : start+end])
				}
			}
		}
	}

	// If no code blocks found, look for SQL keywords and return the relevant portion
	response = strings.TrimSpace(response)
	lines := strings.Split(response, "\n")

	var sqlLines []string
	inSQL := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check if line starts with SQL keywords
		if strings.HasPrefix(strings.ToUpper(line), "CREATE ") ||
			strings.HasPrefix(strings.ToUpper(line), "INSERT ") ||
			strings.HasPrefix(strings.ToUpper(line), "SELECT ") ||
			strings.HasPrefix(strings.ToUpper(line), "DROP ") ||
			strings.HasPrefix(strings.ToUpper(line), "ALTER ") {
			inSQL = true
		}

		// If we're in SQL and hit explanatory text, stop
		if inSQL && (strings.Contains(strings.ToLower(line), "explanation") ||
			strings.Contains(strings.ToLower(line), "note:") ||
			strings.Contains(strings.ToLower(line), "here's") ||
			strings.HasPrefix(line, "The ") ||
			strings.HasPrefix(line, "This ")) {
			break
		}

		if inSQL {
			sqlLines = append(sqlLines, line)
		}
	}

	if len(sqlLines) > 0 {
		return strings.Join(sqlLines, "\n")
	}

	// If all else fails, return the entire response
	return response
}
