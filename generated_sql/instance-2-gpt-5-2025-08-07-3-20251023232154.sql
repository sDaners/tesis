CREATE TABLE authors (
  author_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  author_name STRING(100) NOT NULL,
  nationality STRING(50),
  birth_year INT64,
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (author_id);

CREATE TABLE members (
  member_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  full_name STRING(100) NOT NULL,
  email STRING(150),
  join_date DATE NOT NULL,
  membership_type STRING(30) DEFAULT ('STANDARD'),
  phone_number STRING(20),
  CONSTRAINT check_membership CHECK (membership_type IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR'))
) PRIMARY KEY (member_id);

CREATE TABLE books (
  book_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  title STRING(200) NOT NULL,
  isbn STRING(20),
  publication_year INT64,
  price NUMERIC,
  author_id STRING(36),
  category STRING(50),
  total_copies INT64 DEFAULT (1),
  CONSTRAINT fk_books_author FOREIGN KEY (author_id) REFERENCES authors (author_id)
) PRIMARY KEY (book_id);

CREATE UNIQUE INDEX idx_book_isbn ON books(isbn);

CREATE TABLE loans (
  loan_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  member_id STRING(36),
  book_id STRING(36),
  loan_date DATE NOT NULL,
  due_date DATE NOT NULL,
  return_date DATE,
  fine_amount NUMERIC DEFAULT (0),
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
VALUES (@author_name, @nationality, SAFE_CAST(@birth_year AS INT64))
THEN RETURN author_id;

INSERT INTO books (title, isbn, publication_year, price, author_id, category, total_copies)
VALUES (
  @title,
  @isbn,
  SAFE_CAST(@publication_year AS INT64),
  SAFE_CAST(@price AS NUMERIC),
  CASE
    WHEN REGEXP_CONTAINS(CAST(@author_id AS STRING), r'^[0-9a-fA-F-]{36}$') THEN CAST(@author_id AS STRING)
    ELSE NULL
  END,
  @category,
  SAFE_CAST(@total_copies AS INT64)
)
THEN RETURN book_id;

INSERT INTO members (full_name, email, join_date, membership_type, phone_number)
VALUES (
  @full_name,
  @email,
  COALESCE(SAFE_CAST(@join_date AS DATE), CURRENT_DATE()),
  CASE
    WHEN @membership_type IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR') THEN @membership_type
    ELSE 'STANDARD'
  END,
  @phone_number
)
THEN RETURN member_id;

INSERT INTO loans (member_id, book_id, loan_date, due_date, fine_amount)
VALUES (
  CASE
    WHEN REGEXP_CONTAINS(CAST(@member_id AS STRING), r'^[0-9a-fA-F-]{36}$') THEN CAST(@member_id AS STRING)
    ELSE NULL
  END,
  CASE
    WHEN REGEXP_CONTAINS(CAST(@book_id AS STRING), r'^[0-9a-fA-F-]{36}$') THEN CAST(@book_id AS STRING)
    ELSE NULL
  END,
  COALESCE(SAFE_CAST(@loan_date AS DATE), CURRENT_DATE()),
  COALESCE(SAFE_CAST(@due_date AS DATE), CURRENT_DATE()),
  COALESCE(SAFE_CAST(@fine_amount AS NUMERIC), 0)
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