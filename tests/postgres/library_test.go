package postgres_test

import (
	"database/sql"
	"testing"

	"sql-parser/repo"
	"sql-parser/tools"
)

type LibraryDBTeardown struct {
	db        *sql.DB
	repo      repo.LibraryDatabase
	t         *testing.T
	terminate func()
}

func setupLibraryDB(t *testing.T) *LibraryDBTeardown {
	db, terminate, err := tools.GetDB(false) // Always use PostgreSQL for library tests
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	r := repo.NewLibraryManagementRepo(db)

	if err := r.CleanupDB(); err != nil {
		t.Fatalf("Failed to cleanup DB: %v", err)
	}
	if err := r.CreateTables(); err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}
	return &LibraryDBTeardown{db: db, repo: r, t: t, terminate: terminate}
}

func (d *LibraryDBTeardown) Close() {
	if err := d.repo.CleanupDB(); err != nil {
		d.t.Errorf("Failed to cleanup DB: %v", err)
	}
	d.db.Close()
	d.terminate()
}

func TestLibraryDatabaseOperations(t *testing.T) {
	dbT := setupLibraryDB(t)
	defer dbT.Close()

	// Insert sample data
	authorID, bookID, memberID, loanID, err := dbT.repo.InsertSampleData()
	if err != nil {
		t.Fatalf("InsertSampleData failed: %v", err)
	}
	if authorID == 0 || bookID == 0 || memberID == 0 || loanID == 0 {
		t.Errorf("Expected non-zero IDs, got authorID=%d, bookID=%d, memberID=%d, loanID=%d",
			authorID, bookID, memberID, loanID)
	}

	// Query and check results
	details, err := dbT.repo.QueryLibraryMemberDetails()
	if err != nil {
		t.Fatalf("QueryLibraryMemberDetails failed: %v", err)
	}
	if len(details) == 0 {
		t.Error("Expected at least one library member detail result")
	}

	// Check if we found our test data
	found := false
	for _, detail := range details {
		if detail.FullName == "Jane Smith" &&
			detail.Email == "jane.smith@example.com" &&
			detail.MembershipType == "PREMIUM" &&
			detail.Title.Valid && detail.Title.String == "1984" &&
			detail.AuthorName.Valid && detail.AuthorName.String == "George Orwell" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find inserted member, book, and author in results")
	}
}

func TestLibraryMemberDetails(t *testing.T) {
	dbT := setupLibraryDB(t)
	defer dbT.Close()

	// Insert test data
	_, _, _, _, err := dbT.repo.InsertSampleData()
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Query library member details
	details, err := dbT.repo.QueryLibraryMemberDetails()
	if err != nil {
		t.Fatalf("Failed to query library member details: %v", err)
	}

	// Verify the structure of the results
	if len(details) == 0 {
		t.Error("Expected at least one library member detail")
	}

	for _, detail := range details {
		// Check required fields
		if detail.FullName == "" {
			t.Error("Expected non-empty FullName")
		}
		if detail.Email == "" {
			t.Error("Expected non-empty Email")
		}
		if detail.MembershipType == "" {
			t.Error("Expected non-empty MembershipType")
		}

		// Check nullable fields
		if detail.Title.Valid {
			t.Logf("Book Title: %s", detail.Title.String)
		}
		if detail.ISBN.Valid {
			t.Logf("ISBN: %s", detail.ISBN.String)
		}
		if detail.AuthorName.Valid {
			t.Logf("Author Name: %s", detail.AuthorName.String)
		}
		if detail.LoanDate.Valid {
			t.Logf("Loan Date: %v", detail.LoanDate.Time)
		}
		if detail.DueDate.Valid {
			t.Logf("Due Date: %v", detail.DueDate.Time)
		}
		if detail.ReturnDate.Valid {
			t.Logf("Return Date: %v", detail.ReturnDate.Time)
		}
		if detail.FineAmount.Valid {
			t.Logf("Fine Amount: %.2f", detail.FineAmount.Float64)
		}
	}
}

func TestActiveLoans(t *testing.T) {
	dbT := setupLibraryDB(t)
	defer dbT.Close()

	// Insert test data
	_, _, _, _, err := dbT.repo.InsertSampleData()
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Query active loans
	loans, err := dbT.repo.QueryActiveLoans()
	if err != nil {
		t.Fatalf("Failed to query active loans: %v", err)
	}

	// Verify we have at least one active loan
	if len(loans) == 0 {
		t.Error("Expected at least one active loan")
	}

	// Check the structure of active loans
	for _, loan := range loans {
		if loan.MemberName == "" {
			t.Error("Expected non-empty MemberName")
		}
		if loan.BookTitle == "" {
			t.Error("Expected non-empty BookTitle")
		}
		if loan.AuthorName == "" {
			t.Error("Expected non-empty AuthorName")
		}
		if loan.LoanID == 0 {
			t.Error("Expected non-zero LoanID")
		}

		t.Logf("Active Loan: %s borrowed '%s' by %s",
			loan.MemberName, loan.BookTitle, loan.AuthorName)
	}
}
