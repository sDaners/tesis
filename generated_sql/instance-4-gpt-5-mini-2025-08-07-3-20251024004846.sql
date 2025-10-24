CREATE TABLE authors (
  author_id INT64 NOT NULL GENERATED ALWAYS AS IDENTITY,
  author_name STRING(100) NOT NULL,
  nationality STRING(50),
  birth_year STRING,
  created_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (author_id);

CREATE TABLE books (
  book_id INT64 NOT NULL GENERATED ALWAYS AS IDENTITY,
  title STRING(200) NOT NULL,
  isbn STRING(20),
  publication_year STRING,
  price NUMERIC,
  author_id INT64,
  category STRING(50),
  total_copies INT64 DEFAULT (1),
  CONSTRAINT fk_books_author FOREIGN KEY (author_id) REFERENCES authors(author_id)
) PRIMARY KEY (book_id);

CREATE UNIQUE INDEX idx_book_isbn ON books(isbn);

CREATE TABLE members (
  member_id INT64 NOT NULL GENERATED ALWAYS AS IDENTITY,
  full_name STRING(100) NOT NULL,
  email STRING(150),
  join_date DATE NOT NULL,
  membership_type STRING(30) DEFAULT ('STANDARD'),
  phone_number STRING(20),
  CONSTRAINT check_membership CHECK (membership_type IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR'))
) PRIMARY KEY (member_id);

CREATE TABLE loans (
  loan_id INT64 NOT NULL GENERATED ALWAYS AS IDENTITY,
  member_id INT64,
  book_id INT64,
  loan_date DATE NOT NULL,
  due_date DATE NOT NULL,
  return_date DATE,
  fine_amount NUMERIC DEFAULT (0),
  CONSTRAINT check_loan_dates CHECK (due_date > loan_date),
  CONSTRAINT fk_loans_member FOREIGN KEY (member_id) REFERENCES members(member_id),
  CONSTRAINT fk_loans_book FOREIGN KEY (book_id) REFERENCES books(book_id)
) PRIMARY KEY (loan_id);

CREATE INDEX idx_book_title ON books(title);
CREATE INDEX idx_author_name ON authors(author_name);
CREATE INDEX idx_member_email ON members(email);
CREATE INDEX idx_loan_status ON loans(return_date);

CREATE VIEW active_loans SQL SECURITY INVOKER AS
SELECT 
  l.loan_id,
  m.full_name AS member_name,
  b.title AS book_title,
  a.author_name,
  l.loan_date,
  l.due_date,
  l.fine_amount
FROM loans l
JOIN members m ON l.member_id = m.member_id
JOIN books b ON l.book_id = b.book_id
JOIN authors a ON b.author_id = a.author_id
WHERE l.return_date IS NULL;

-- Insert an author (use parameters in client code where possible)
-- Example with parameters (client): @author_name, @nationality, @birth_year
INSERT INTO authors (author_name, nationality, birth_year)
VALUES (@author_name, @nationality, @birth_year);

-- If client code cannot supply parameters, the following fallback selects the most recently created author:
SELECT author_id
FROM authors
ORDER BY created_at DESC
LIMIT 1;

-- Insert a book (use parameters in client code where possible)
INSERT INTO books (title, isbn, publication_year, price, author_id, category, total_copies)
VALUES (@title, @isbn, @publication_year, @price, @author_id, @category, @total_copies);

-- Fallback to get last inserted book (or by unique isbn if available)
SELECT book_id
FROM books
WHERE isbn IS NOT NULL AND isbn != ''
ORDER BY book_id DESC
LIMIT 1;

-- Insert a member (use parameters in client code where possible)
INSERT INTO members (full_name, email, join_date, membership_type, phone_number)
VALUES (@full_name, @email, @join_date, @membership_type, @phone_number);

-- Fallback to get last inserted member
SELECT member_id
FROM members
ORDER BY member_id DESC
LIMIT 1;

-- Insert a loan (use parameters in client code where possible)
INSERT INTO loans (member_id, book_id, loan_date, due_date, fine_amount)
VALUES (@member_id, @book_id, @loan_date, @due_date, @fine_amount);

-- Fallback to get last inserted loan
SELECT loan_id
FROM loans
ORDER BY loan_id DESC
LIMIT 1;

-- Example query: list members with their loans (works without parameters)
SELECT m.member_id, m.full_name, m.email, m.membership_type,
       b.title, b.isbn, a.author_name,
       l.loan_date, l.due_date, l.return_date, l.fine_amount
FROM members m
LEFT JOIN loans l ON m.member_id = l.member_id
LEFT JOIN books b ON l.book_id = b.book_id
LEFT JOIN authors a ON b.author_id = a.author_id;

-- Clean up statements
DROP VIEW active_loans;
DROP INDEX idx_loan_status;
DROP INDEX idx_member_email;
DROP INDEX idx_author_name;
DROP INDEX idx_book_title;
DROP INDEX idx_book_isbn;
DROP TABLE loans;
DROP TABLE books;
DROP TABLE members;
DROP TABLE authors;