package parser

import (
	"postgres-example/tools"
	"strings"
	"testing"
)

func TestSQLParser(t *testing.T) {
	// Test SQL content with various statements and comments
	sqlContent := `
		-- This is a single line comment
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) -- inline comment
		);
		
		/* This is a
		   multi-line comment */
		INSERT INTO users (name) VALUES ('John Doe');
		INSERT INTO users (name) VALUES ('Jane Smith');
		
		-- Another comment
		SELECT * FROM users WHERE id = 1;
		
		/* Another multi-line comment
		   that spans multiple lines */
		SELECT name FROM users;
		
		CREATE INDEX idx_users_name ON users(name);
	`

	statements, err := tools.ParseSQLFromString(sqlContent)
	if err != nil {
		t.Fatalf("Failed to parse SQL: %v", err)
	}

	// Expected statements (comments should be removed)
	expected := []string{
		"CREATE TABLE users ( id SERIAL PRIMARY KEY, name VARCHAR(100) )",
		"INSERT INTO users (name) VALUES ('John Doe')",
		"INSERT INTO users (name) VALUES ('Jane Smith')",
		"SELECT * FROM users WHERE id = 1",
		"SELECT name FROM users",
		"CREATE INDEX idx_users_name ON users(name)",
	}

	if len(statements) != len(expected) {
		t.Fatalf("Expected %d statements, got %d", len(expected), len(statements))
	}

	for i, stmt := range statements {
		// Normalize whitespace for comparison
		normalizedStmt := strings.Join(strings.Fields(stmt), " ")
		normalizedExpected := strings.Join(strings.Fields(expected[i]), " ")

		if normalizedStmt != normalizedExpected {
			t.Errorf("Statement %d mismatch:\nGot:      %s\nExpected: %s", i, normalizedStmt, normalizedExpected)
		}
	}
}
