package repo

// LibraryManagementDDL contains all the table creation statements for the library management database.
// It creates the following structure:
// - authors: Stores author information with auto-incrementing IDs
// - books: Stores book information with references to authors
// - members: Stores member information with membership types
// - loans: Links members to books with loan dates and fine tracking
// Also creates necessary indexes and a view for active loans.
const LibraryManagementDDL = `
CREATE TABLE IF NOT EXISTS authors (
    author_id SERIAL PRIMARY KEY,
    author_name VARCHAR(100) NOT NULL,
    nationality VARCHAR(50),
    birth_year INTEGER,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS books (
    book_id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    isbn VARCHAR(20),
    publication_year INTEGER,
    price NUMERIC,
    author_id INTEGER REFERENCES authors(author_id),
    category VARCHAR(50),
    total_copies INTEGER DEFAULT 1
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_book_isbn ON books(isbn);

CREATE TABLE IF NOT EXISTS members (
    member_id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(150),
    join_date DATE NOT NULL,
    membership_type VARCHAR(30) DEFAULT 'STANDARD',
    phone_number VARCHAR(20),
    CONSTRAINT check_membership CHECK (membership_type IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR'))
);

CREATE TABLE IF NOT EXISTS loans (
    loan_id SERIAL PRIMARY KEY,
    member_id INTEGER,
    book_id INTEGER,
    loan_date DATE NOT NULL,
    due_date DATE NOT NULL,
    return_date DATE,
    fine_amount NUMERIC DEFAULT 0.0,
    CONSTRAINT check_loan_dates CHECK (due_date > loan_date),
    FOREIGN KEY (member_id) REFERENCES members(member_id),
    FOREIGN KEY (book_id) REFERENCES books(book_id)
);

CREATE INDEX IF NOT EXISTS idx_book_title ON books(title);
CREATE INDEX IF NOT EXISTS idx_author_name ON authors(author_name);
CREATE INDEX IF NOT EXISTS idx_member_email ON members(email);
CREATE INDEX IF NOT EXISTS idx_loan_status ON loans(return_date);

CREATE OR REPLACE VIEW active_loans AS
SELECT 
    l.loan_id,
    m.full_name as member_name,
    b.title as book_title,
    a.author_name,
    l.loan_date,
    l.due_date,
    l.fine_amount
FROM loans l
JOIN members m ON l.member_id = m.member_id
JOIN books b ON l.book_id = b.book_id
JOIN authors a ON b.author_id = a.author_id
WHERE l.return_date IS NULL;
`

// InsertAuthorSQL inserts a new author and returns their ID.
// Parameters:
//   - $1: Author's name
//   - $2: Author's nationality
//   - $3: Author's birth year
//
// Returns: The generated author ID (SERIAL)
const InsertAuthorSQL = `
INSERT INTO authors (author_name, nationality, birth_year)
VALUES ($1, $2, $3)
RETURNING author_id`

// InsertBookSQL inserts a new book and returns its ID.
// Parameters:
//   - $1: Book title
//   - $2: Book ISBN
//   - $3: Publication year
//   - $4: Book price
//   - $5: Author ID
//   - $6: Book category
//   - $7: Total copies available
//
// Returns: The generated book ID (SERIAL)
const InsertBookSQL = `
INSERT INTO books (title, isbn, publication_year, price, author_id, category, total_copies)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING book_id`

// InsertMemberSQL inserts a new member and returns their ID.
// Parameters:
//   - $1: Member's full name
//   - $2: Member's email address
//   - $3: Member's join date
//   - $4: Membership type (must be one of: 'STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR')
//   - $5: Member's phone number
//
// Returns: The generated member ID (SERIAL)
const InsertMemberSQL = `
INSERT INTO members (full_name, email, join_date, membership_type, phone_number)
VALUES ($1, $2, $3, $4, $5)
RETURNING member_id`

// InsertLoanSQL creates a new loan and returns its ID.
// Parameters:
//   - $1: Member ID
//   - $2: Book ID
//   - $3: Loan date
//   - $4: Due date
//   - $5: Fine amount
//
// Returns: The generated loan ID (SERIAL)
const InsertLoanSQL = `
INSERT INTO loans (member_id, book_id, loan_date, due_date, fine_amount)
VALUES ($1, $2, $3, $4, $5)
RETURNING loan_id`

// QueryLibraryMemberDetailsSQL retrieves detailed information about library members including:
// - Basic member information (ID, name, email, membership type)
// - Book information for any loans
// - Author information
// - Loan details (dates, fines)
// The query joins multiple tables to provide a comprehensive view of each member's details.
const QueryLibraryMemberDetailsSQL = `
SELECT m.member_id, m.full_name, m.email, m.membership_type,
       b.title, b.isbn, a.author_name,
       l.loan_date, l.due_date, l.return_date, l.fine_amount
FROM members m
LEFT JOIN loans l ON m.member_id = l.member_id
LEFT JOIN books b ON l.book_id = b.book_id
LEFT JOIN authors a ON b.author_id = a.author_id
`

// QueryActiveLoansSQL retrieves all active loans (loans that haven't been returned).
const QueryActiveLoansSQL = `
SELECT * FROM active_loans
`

// CleanupLibraryDBSQL contains the SQL statements to clean up the library management database.
// The statements are ordered to respect dependencies:
// 1. Drop views
// 2. Drop indexes
// 3. Drop tables in reverse order of their dependencies
const CleanupLibraryDBSQL = `--sql
DROP VIEW IF EXISTS active_loans;
DROP INDEX IF EXISTS idx_loan_status;
DROP INDEX IF EXISTS idx_member_email;
DROP INDEX IF EXISTS idx_author_name;
DROP INDEX IF EXISTS idx_book_title;
DROP INDEX IF EXISTS idx_book_isbn;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS members;
DROP TABLE IF EXISTS authors;
`
