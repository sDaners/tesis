CREATE SEQUENCE authors_author_id_seq OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE SEQUENCE books_book_id_seq OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE SEQUENCE members_member_id_seq OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE SEQUENCE loans_loan_id_seq OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE TABLE authors (
  author_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE authors_author_id_seq)),
  author_name STRING(100) NOT NULL,
  nationality STRING(50),
  birth_year INT64,
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (author_id);

CREATE TABLE books (
  book_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE books_book_id_seq)),
  title STRING(200) NOT NULL,
  isbn STRING(20),
  publication_year INT64,
  price NUMERIC,
  author_id INT64,
  category STRING(50),
  total_copies INT64 DEFAULT (1),
  CONSTRAINT fk_books_author FOREIGN KEY (author_id) REFERENCES authors (author_id)
) PRIMARY KEY (book_id);

CREATE UNIQUE INDEX idx_book_isbn ON books(isbn);

CREATE TABLE members (
  member_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE members_member_id_seq)),
  full_name STRING(100) NOT NULL,
  email STRING(150),
  join_date DATE NOT NULL,
  membership_type STRING(30) DEFAULT ('STANDARD'),
  phone_number STRING(20),
  CONSTRAINT check_membership CHECK (membership_type IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR'))
) PRIMARY KEY (member_id);

CREATE TABLE loans (
  loan_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE loans_loan_id_seq)),
  member_id INT64,
  book_id INT64,
  loan_date DATE NOT NULL,
  due_date DATE NOT NULL,
  return_date DATE,
  fine_amount NUMERIC DEFAULT (CAST(0 AS NUMERIC)),
  CONSTRAINT check_loan_dates CHECK (due_date > loan_date),
  CONSTRAINT fk_loans_member FOREIGN KEY (member_id) REFERENCES members (member_id),
  CONSTRAINT fk_loans_book FOREIGN KEY (book_id) REFERENCES books (book_id)
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
VALUES (@p1, @p2, SAFE_CAST(@p3 AS INT64))
THEN RETURN author_id;

INSERT INTO books (title, isbn, publication_year, price, author_id, category, total_copies)
VALUES (
  @p1,
  @p2,
  SAFE_CAST(@p3 AS INT64),
  COALESCE(SAFE_CAST(@p4 AS NUMERIC), CAST(0 AS NUMERIC)),
  SAFE_CAST(@p5 AS INT64),
  @p6,
  COALESCE(SAFE_CAST(@p7 AS INT64), 1)
)
THEN RETURN book_id;

INSERT INTO members (full_name, email, join_date, membership_type, phone_number)
VALUES (
  @p1,
  @p2,
  COALESCE(SAFE_CAST(@p3 AS DATE), CURRENT_DATE()),
  CASE WHEN @p4 IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR') THEN @p4 ELSE 'STANDARD' END,
  @p5
)
THEN RETURN member_id;

INSERT INTO loans (member_id, book_id, loan_date, due_date, fine_amount)
VALUES (
  SAFE_CAST(@p1 AS INT64),
  SAFE_CAST(@p2 AS INT64),
  COALESCE(SAFE_CAST(@p3 AS DATE), CURRENT_DATE()),
  CASE
    WHEN SAFE_CAST(@p4 AS DATE) IS NOT NULL THEN SAFE_CAST(@p4 AS DATE)
    ELSE DATE_ADD(COALESCE(SAFE_CAST(@p3 AS DATE), CURRENT_DATE()), INTERVAL 14 DAY)
  END,
  COALESCE(SAFE_CAST(@p5 AS NUMERIC), CAST(0 AS NUMERIC))
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
DROP SEQUENCE loans_loan_id_seq;
DROP SEQUENCE members_member_id_seq;
DROP SEQUENCE books_book_id_seq;
DROP SEQUENCE authors_author_id_seq;