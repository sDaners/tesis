package models

import (
	"database/sql"
	"time"
)

// Author represents an author in the library system
type Author struct {
	AuthorID    int       `db:"author_id"`
	AuthorName  string    `db:"author_name"`
	Nationality string    `db:"nationality"`
	BirthYear   int       `db:"birth_year"`
	CreatedAt   time.Time `db:"created_at"`
}

// Book represents a book in the library system
type Book struct {
	BookID          int     `db:"book_id"`
	Title           string  `db:"title"`
	ISBN            string  `db:"isbn"`
	PublicationYear int     `db:"publication_year"`
	Price           float64 `db:"price"`
	AuthorID        int     `db:"author_id"`
	Category        string  `db:"category"`
	TotalCopies     int     `db:"total_copies"`
}

// Member represents a library member
type Member struct {
	MemberID       int       `db:"member_id"`
	FullName       string    `db:"full_name"`
	Email          string    `db:"email"`
	JoinDate       time.Time `db:"join_date"`
	MembershipType string    `db:"membership_type"`
	PhoneNumber    string    `db:"phone_number"`
}

// Loan represents a book loan
type Loan struct {
	LoanID     int          `db:"loan_id"`
	MemberID   int          `db:"member_id"`
	BookID     int          `db:"book_id"`
	LoanDate   time.Time    `db:"loan_date"`
	DueDate    time.Time    `db:"due_date"`
	ReturnDate sql.NullTime `db:"return_date"`
	FineAmount float64      `db:"fine_amount"`
}

// LibraryMemberDetails represents detailed information about library members and their loans
type LibraryMemberDetails struct {
	MemberID       int             `db:"member_id"`
	FullName       string          `db:"full_name"`
	Email          string          `db:"email"`
	MembershipType string          `db:"membership_type"`
	Title          sql.NullString  `db:"title"`
	ISBN           sql.NullString  `db:"isbn"`
	AuthorName     sql.NullString  `db:"author_name"`
	LoanDate       sql.NullTime    `db:"loan_date"`
	DueDate        sql.NullTime    `db:"due_date"`
	ReturnDate     sql.NullTime    `db:"return_date"`
	FineAmount     sql.NullFloat64 `db:"fine_amount"`
}

// ActiveLoan represents an active loan from the view
type ActiveLoan struct {
	LoanID     int       `db:"loan_id"`
	MemberName string    `db:"member_name"`
	BookTitle  string    `db:"book_title"`
	AuthorName string    `db:"author_name"`
	LoanDate   time.Time `db:"loan_date"`
	DueDate    time.Time `db:"due_date"`
	FineAmount float64   `db:"fine_amount"`
}
