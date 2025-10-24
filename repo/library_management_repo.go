package repo

import (
	"context"
	"database/sql"
	"time"

	"sql-parser/models"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type LibraryManagementRepo struct {
	DB *sql.DB
}

func NewLibraryManagementRepo(db *sql.DB) *LibraryManagementRepo {
	return &LibraryManagementRepo{DB: db}
}

func (r *LibraryManagementRepo) CreateTables() error {
	_, err := r.DB.Exec(LibraryManagementDDL)
	return err
}

func (r *LibraryManagementRepo) InsertSampleData() (int, int, int, int, error) {
	// Insert author
	var authorID int
	err := r.DB.QueryRow(InsertAuthorSQL, "George Orwell", "British", 1903).Scan(&authorID)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Insert book
	var bookID int
	err = r.DB.QueryRow(InsertBookSQL,
		"1984", "978-0-452-28423-4", 1949, 12.99, authorID, "Fiction", 3).Scan(&bookID)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Insert member
	var memberID int
	err = r.DB.QueryRow(InsertMemberSQL,
		"Jane Smith", "jane.smith@example.com", time.Now(), "PREMIUM", "555-0123").Scan(&memberID)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Insert loan
	var loanID int
	loanDate := time.Now()
	dueDate := loanDate.AddDate(0, 0, 14) // 14 days from now
	err = r.DB.QueryRow(InsertLoanSQL,
		memberID, bookID, loanDate, dueDate, 0.00).Scan(&loanID)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return authorID, bookID, memberID, loanID, nil
}

func (r *LibraryManagementRepo) QueryLibraryMemberDetails() ([]models.LibraryMemberDetails, error) {
	query := QueryLibraryMemberDetailsSQL

	var details []models.LibraryMemberDetails
	err := sqlscan.Select(context.Background(), r.DB, &details, query)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (r *LibraryManagementRepo) QueryActiveLoans() ([]models.ActiveLoan, error) {
	query := QueryActiveLoansSQL

	var loans []models.ActiveLoan
	err := sqlscan.Select(context.Background(), r.DB, &loans, query)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *LibraryManagementRepo) CleanupDB() error {
	_, err := r.DB.Exec(CleanupLibraryDBSQL)
	return err
}
