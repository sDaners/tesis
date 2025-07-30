package tools

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudspannerecosystem/memefish"
)

// ExtractStatementsFromFile extracts SQL statements from a file using memefish
func ExtractStatementsFromFile(filename string) ([]string, error) {
	// Read the entire file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Clean comments before parsing
	cleanedContent := cleanComments(string(content))

	// Use memefish to split the content into raw statements
	rawStatements, err := memefish.SplitRawStatements(filename, cleanedContent)
	if err != nil {
		return nil, fmt.Errorf("failed to split statements: %w", err)
	}

	// Convert raw statements to strings and filter out empty ones
	var statements []string
	for _, rawStmt := range rawStatements {
		stmt := strings.TrimSpace(rawStmt.Statement)
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements, nil
}

func cleanComments(content string) string {
	lines := strings.Split(content, "\n")
	var cleanedLines []string

	for _, line := range lines {
		commentPos := strings.Index(line, "--")
		if commentPos >= 0 {
			// Remove everything from "--" to end of line
			line = line[:commentPos]
		}
		cleanedLines = append(cleanedLines, line)
	}

	return strings.Join(cleanedLines, "\n")
}
