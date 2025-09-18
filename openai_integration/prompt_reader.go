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

	return response
}
