package tools

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// SQLParser represents a parser for SQL files
type SQLParser struct {
	filename string
}

// NewSQLParser creates a new SQL parser instance
func NewSQLParser(filename string) *SQLParser {
	return &SQLParser{
		filename: filename,
	}
}

// ParseSQL reads an SQL file and returns a list of SQL statements
// separated by their first word (CREATE, INSERT, SELECT)
// Comments are ignored
func (p *SQLParser) ParseSQL() ([]string, error) {
	file, err := os.Open(p.filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", p.filename, err)
	}
	defer file.Close()

	return p.parseFromReader(file)
}

// parseFromReader parses SQL from an io.Reader
func (p *SQLParser) parseFromReader(reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	var lines []string

	// First pass: read all lines and remove comments
	for scanner.Scan() {
		line := scanner.Text()
		cleanLine := p.removeComments(line)
		if strings.TrimSpace(cleanLine) != "" {
			lines = append(lines, cleanLine)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Join all lines into a single string
	fullContent := strings.Join(lines, " ")

	// Split into statements and filter by target keywords
	statements := p.extractStatements(fullContent)

	return statements, nil
}

// removeComments removes SQL comments from a line
// Handles both single-line comments (--) and multi-line comments (/* */)
func (p *SQLParser) removeComments(line string) string {
	// Handle single-line comments (--)
	if idx := strings.Index(line, "--"); idx != -1 {
		line = line[:idx]
	}

	// Handle multi-line comments (/* */)
	for {
		startIdx := strings.Index(line, "/*")
		if startIdx == -1 {
			break
		}

		endIdx := strings.Index(line[startIdx:], "*/")
		if endIdx == -1 {
			// Comment continues to end of line
			line = line[:startIdx]
			break
		}

		// Remove the comment block
		line = line[:startIdx] + line[startIdx+endIdx+2:]
	}

	return line
}

// extractStatements extracts SQL statements that start with CREATE, INSERT, or SELECT
func (p *SQLParser) extractStatements(content string) []string {
	var statements []string

	// Regular expression to find statements starting with CREATE, INSERT, or SELECT
	// This regex looks for these keywords at word boundaries and captures everything until
	// we find another statement or the end of the content
	re := regexp.MustCompile(`(?i)\b(CREATE|INSERT|SELECT)\b[^;]*(?:;|$)`)

	matches := re.FindAllString(content, -1)

	for _, match := range matches {
		// Clean up the statement
		statement := strings.TrimSpace(match)
		if statement != "" {
			// Remove trailing semicolon if present
			statement = strings.TrimSuffix(statement, ";")
			statement = strings.TrimSpace(statement)
			if statement != "" {
				statements = append(statements, statement)
			}
		}
	}

	// If no statements found with semicolons, try a different approach
	// Split by keywords and clean up
	if len(statements) == 0 {
		statements = p.extractStatementsByKeyword(content)
	}

	return statements
}

// extractStatementsByKeyword extracts statements by splitting on keywords
func (p *SQLParser) extractStatementsByKeyword(content string) []string {
	var statements []string

	// Normalize whitespace
	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
	content = strings.TrimSpace(content)

	// Split by keywords while preserving the keywords
	keywords := []string{"CREATE", "INSERT", "SELECT"}

	// Case-insensitive split
	parts := []string{content}

	for _, keyword := range keywords {
		var newParts []string
		for _, part := range parts {
			// Case-insensitive split that preserves the delimiter
			re := regexp.MustCompile(`(?i)\b` + keyword + `\b`)
			splits := re.Split(part, -1)
			matches := re.FindAllString(part, -1)

			if len(matches) > 0 {
				// Add first part if not empty
				if strings.TrimSpace(splits[0]) != "" {
					newParts = append(newParts, strings.TrimSpace(splits[0]))
				}

				// Add keyword + remaining parts
				for i, match := range matches {
					if i+1 < len(splits) {
						combined := match + " " + splits[i+1]
						newParts = append(newParts, strings.TrimSpace(combined))
					} else {
						newParts = append(newParts, strings.TrimSpace(match))
					}
				}
			} else {
				if strings.TrimSpace(part) != "" {
					newParts = append(newParts, strings.TrimSpace(part))
				}
			}
		}
		parts = newParts
	}

	// Filter parts that start with our target keywords
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			upperPart := strings.ToUpper(part)
			if strings.HasPrefix(upperPart, "CREATE ") ||
				strings.HasPrefix(upperPart, "INSERT ") ||
				strings.HasPrefix(upperPart, "SELECT ") {
				statements = append(statements, part)
			}
		}
	}

	return statements
}

// ParseSQLFromString parses SQL statements from a string
func ParseSQLFromString(content string) ([]string, error) {
	parser := &SQLParser{}
	reader := strings.NewReader(content)
	return parser.parseFromReader(reader)
}
