CREATE TABLE authors (
  author_id STRING(36) NOT NULL,
  author_name STRING(100) NOT NULL,
  nationality STRING(50),
  birth_year INT64,
  created_at TIMESTAMP OPTIONS (allow_commit_timestamp = true)
) PRIMARY KEY (author_id);

CREATE TABLE books (
  book_id STRING(36) NOT NULL,
  title STRING(200) NOT NULL,
  isbn STRING(20),
  publication_year INT64,
  price NUMERIC,
  author_id STRING(36),
  category STRING(50),
  total_copies INT64 DEFAULT 1,
  CONSTRAINT fk_books_author FOREIGN KEY (author_id) REFERENCES authors(author_id)
) PRIMARY KEY (book_id);

CREATE UNIQUE INDEX idx_book_isbn ON books(isbn);

CREATE TABLE members (
  member_id STRING(36) NOT NULL,
  full_name STRING(100) NOT NULL,
  email STRING(150),
  join_date DATE NOT NULL,
  membership_type STRING(30) NOT NULL DEFAULT 'STANDARD',
  phone_number STRING(20),
  CONSTRAINT check_membership CHECK (membership_type IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR'))
) PRIMARY KEY (member_id);

CREATE TABLE loans (
  loan_id STRING(36) NOT NULL,
  member_id STRING(36),
  book_id STRING(36),
  loan_date DATE NOT NULL,
  due_date DATE NOT NULL,
  return_date DATE,
  fine_amount NUMERIC DEFAULT 0,
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
  l.loan_id AS loan_id,
  m.full_name AS member_name,
  b.title AS book_title,
  a.author_name AS author_name,
  l.loan_date AS loan_date,
  l.due_date AS due_date,
  l.fine_amount AS fine_amount
FROM loans AS l
JOIN members AS m ON l.member_id = m.member_id
JOIN books AS b ON l.book_id = b.book_id
JOIN authors AS a ON b.author_id = a.author_id
WHERE l.return_date IS NULL;

INSERT INTO authors (author_id, author_name, nationality, birth_year, created_at)
VALUES (GENERATE_UUID(), @p1, @p2, CAST(@p3 AS INT64), PENDING_COMMIT_TIMESTAMP());

INSERT INTO books (book_id, title, isbn, publication_year, price, author_id, category, total_copies)
VALUES (GENERATE_UUID(), @p1, @p2, CAST(@p3 AS INT64), CAST(@p4 AS NUMERIC), @p5, @p6, CAST(@p7 AS INT64));

INSERT INTO members (member_id, full_name, email, join_date, membership_type, phone_number)
VALUES (GENERATE_UUID(), @p1, @p2, CAST(@p3 AS DATE), @p4, @p5);

INSERT INTO loans (loan_id, member_id, book_id, loan_date, due_date, fine_amount)
VALUES (GENERATE_UUID(), @p1, @p2, CAST(@p3 AS DATE), CAST(@p4 AS DATE), CAST(@p5 AS NUMERIC));

SELECT m.member_id, m.full_name, m.email, m.membership_type,
       b.title, b.isbn, a.author_name,
       l.loan_date, l.due_date, l.return_date, l.fine_amount
FROM members AS m
LEFT JOIN loans AS l ON m.member_id = l.member_id
LEFT JOIN books AS b ON l.book_id = b.book_id
LEFT JOIN authors AS a ON b.author_id = a.author_id;

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