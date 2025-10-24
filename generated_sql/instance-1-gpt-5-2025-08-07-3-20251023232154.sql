CREATE SEQUENCE seq_authors_author_id OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE SEQUENCE seq_books_book_id OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE SEQUENCE seq_members_member_id OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE SEQUENCE seq_loans_loan_id OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE TABLE authors (
  author_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE seq_authors_author_id)),
  author_name STRING(100) NOT NULL,
  nationality STRING(50),
  birth_year INT64,
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (author_id);

CREATE TABLE books (
  book_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE seq_books_book_id)),
  title STRING(200) NOT NULL,
  isbn STRING(20),
  publication_year INT64,
  price NUMERIC,
  author_id INT64,
  category STRING(50),
  total_copies INT64 DEFAULT (1),
  FOREIGN KEY (author_id) REFERENCES authors(author_id)
) PRIMARY KEY (book_id);

CREATE UNIQUE NULL_FILTERED INDEX idx_book_isbn ON books(isbn);

CREATE TABLE members (
  member_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE seq_members_member_id)),
  full_name STRING(100) NOT NULL,
  email STRING(150),
  join_date DATE NOT NULL,
  membership_type STRING(30) DEFAULT ('STANDARD'),
  phone_number STRING(20),
  CONSTRAINT check_membership CHECK (membership_type IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR'))
) PRIMARY KEY (member_id);

CREATE TABLE loans (
  loan_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE seq_loans_loan_id)),
  member_id INT64,
  book_id INT64,
  loan_date DATE NOT NULL,
  due_date DATE NOT NULL,
  return_date DATE,
  fine_amount NUMERIC DEFAULT (0),
  CONSTRAINT check_loan_dates CHECK (due_date > loan_date),
  FOREIGN KEY (member_id) REFERENCES members(member_id),
  FOREIGN KEY (book_id) REFERENCES books(book_id)
) PRIMARY KEY (loan_id);

CREATE INDEX idx_book_title ON books(title);
CREATE INDEX idx_author_name ON authors(author_name);
CREATE INDEX idx_member_email ON members(email);
CREATE INDEX idx_loan_status ON loans(return_date);

CREATE OR REPLACE VIEW active_loans SQL SECURITY INVOKER AS
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

INSERT INTO authors (author_name, nationality, birth_year)
VALUES (CAST(@p1 AS STRING), CAST(@p2 AS STRING), SAFE_CAST(@p3 AS INT64))
THEN RETURN author_id;

INSERT INTO books (title, isbn, publication_year, price, author_id, category, total_copies)
VALUES (
  CAST(@p1 AS STRING),
  CAST(@p2 AS STRING),
  SAFE_CAST(@p3 AS INT64),
  COALESCE(SAFE_CAST(@p4 AS NUMERIC), 0),
  SAFE_CAST(@p5 AS INT64),
  CAST(@p6 AS STRING),
  COALESCE(SAFE_CAST(@p7 AS INT64), 1)
)
THEN RETURN book_id;

INSERT INTO members (full_name, email, join_date, membership_type, phone_number)
VALUES (
  CAST(@p1 AS STRING),
  CAST(@p2 AS STRING),
  COALESCE(SAFE_CAST(@p3 AS DATE), CURRENT_DATE()),
  CASE
    WHEN CAST(@p4 AS STRING) IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR') THEN CAST(@p4 AS STRING)
    ELSE 'STANDARD'
  END,
  CAST(@p5 AS STRING)
)
THEN RETURN member_id;

INSERT INTO loans (member_id, book_id, loan_date, due_date, fine_amount)
VALUES (
  SAFE_CAST(@p1 AS INT64),
  SAFE_CAST(@p2 AS INT64),
  COALESCE(SAFE_CAST(@p3 AS DATE), CURRENT_DATE()),
  COALESCE(SAFE_CAST(@p4 AS DATE), DATE_ADD(CURRENT_DATE(), INTERVAL 14 DAY)),
  COALESCE(SAFE_CAST(@p5 AS NUMERIC), 0)
)
THEN RETURN loan_id;

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