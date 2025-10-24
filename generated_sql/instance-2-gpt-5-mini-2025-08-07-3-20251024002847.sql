CREATE TABLE authors (
  author_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  author_name STRING(100) NOT NULL,
  nationality STRING(50),
  birth_year STRING(10),
  created_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (author_id);

CREATE TABLE books (
  book_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  title STRING(200) NOT NULL,
  isbn STRING(20),
  publication_year STRING(10),
  price STRING,
  author_id STRING(36),
  category STRING(50),
  total_copies INT64 NOT NULL DEFAULT (1),
  created_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP()),
  FOREIGN KEY (author_id) REFERENCES authors(author_id)
) PRIMARY KEY (book_id);

CREATE UNIQUE INDEX idx_book_isbn ON books(isbn);

CREATE TABLE members (
  member_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  full_name STRING(100) NOT NULL,
  email STRING(150),
  join_date STRING,
  membership_type STRING(30) NOT NULL DEFAULT ('STANDARD'),
  phone_number STRING(20),
  CONSTRAINT check_membership CHECK (membership_type IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR')),
  created_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (member_id);

CREATE TABLE loans (
  loan_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  member_id STRING(36),
  book_id STRING(36),
  loan_date STRING NOT NULL,
  due_date STRING NOT NULL,
  return_date STRING,
  fine_amount STRING NOT NULL DEFAULT ('0'),
  created_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP()),
  CONSTRAINT check_loan_dates CHECK (due_date > loan_date),
  FOREIGN KEY (member_id) REFERENCES members(member_id),
  FOREIGN KEY (book_id) REFERENCES books(book_id)
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
  a.author_name AS author_name,
  l.loan_date,
  l.due_date,
  l.fine_amount
FROM loans l
JOIN members m ON l.member_id = m.member_id
JOIN books b ON l.book_id = b.book_id
JOIN authors a ON b.author_id = a.author_id
WHERE l.return_date IS NULL;

-- Insert author (parameters expected: @p1 author_name, @p2 nationality, @p3 birth_year)
INSERT INTO authors (author_id, author_name, nationality, birth_year)
VALUES (GENERATE_UUID(), @p1, @p2, @p3);

-- Retrieve the most recently created author id (useful if client didn't supply parameters)
SELECT author_id FROM authors
WHERE author_name = @p1
ORDER BY created_at DESC
LIMIT 1;

-- Insert book (parameters: @p1 title, @p2 isbn, @p3 publication_year, @p4 price, @p5 author_id, @p6 category, @p7 total_copies)
INSERT INTO books (book_id, title, isbn, publication_year, price, author_id, category, total_copies)
VALUES (GENERATE_UUID(), @p1, @p2, @p3, @p4, @p5, @p6, @p7);

-- Retrieve book id by ISBN (isbn is unique)
SELECT book_id FROM books
WHERE isbn = @p2
LIMIT 1;

-- Insert member (parameters: @p1 full_name, @p2 email, @p3 join_date, @p4 membership_type, @p5 phone_number)
INSERT INTO members (member_id, full_name, email, join_date, membership_type, phone_number)
VALUES (GENERATE_UUID(), @p1, @p2, @p3, @p4, @p5);

-- Retrieve most recently created member id (by email if provided)
SELECT member_id FROM members
WHERE email = @p2
ORDER BY created_at DESC
LIMIT 1;

-- Insert loan (parameters: @p1 member_id, @p2 book_id, @p3 loan_date, @p4 due_date, @p5 fine_amount)
INSERT INTO loans (loan_id, member_id, book_id, loan_date, due_date, fine_amount)
VALUES (GENERATE_UUID(), @p1, @p2, @p3, @p4, @p5);

-- Retrieve most recently created loan id for a given member/book/date
SELECT loan_id FROM loans
WHERE member_id = @p1 AND book_id = @p2 AND loan_date = @p3
ORDER BY created_at DESC
LIMIT 1;

-- Select members with their loans and related book/author info
SELECT m.member_id, m.full_name, m.email, m.membership_type,
       b.title, b.isbn, a.author_name,
       l.loan_date, l.due_date, l.return_date, l.fine_amount
FROM members m
LEFT JOIN loans l ON m.member_id = l.member_id
LEFT JOIN books b ON l.book_id = b.book_id
LEFT JOIN authors a ON b.author_id = a.author_id;

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