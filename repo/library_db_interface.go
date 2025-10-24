package repo

import "sql-parser/models"

// LibraryDatabase defines the interface for library management database operations.
type LibraryDatabase interface {
	CreateTables() error
	InsertSampleData() (int, int, int, int, error) // Returns authorID, bookID, memberID, loanID
	QueryLibraryMemberDetails() ([]models.LibraryMemberDetails, error)
	QueryActiveLoans() ([]models.ActiveLoan, error)
	CleanupDB() error
}
