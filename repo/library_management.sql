CREATE TABLE authors (
    author_id SERIAL PRIMARY KEY,
    author_name VARCHAR(100) NOT NULL,
    nationality VARCHAR(50),
    birth_year INTEGER,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE books (
    book_id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    isbn VARCHAR(20),
    publication_year INTEGER,
    price NUMERIC,
    author_id INTEGER REFERENCES authors(author_id),
    category VARCHAR(50),
    total_copies INTEGER DEFAULT 1
);

CREATE UNIQUE INDEX idx_book_isbn ON books(isbn);

CREATE TABLE members (
    member_id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(150),
    join_date DATE NOT NULL,
    membership_type VARCHAR(30) DEFAULT 'STANDARD',
    phone_number VARCHAR(20),
    CONSTRAINT check_membership CHECK (membership_type IN ('STANDARD', 'PREMIUM', 'STUDENT', 'SENIOR'))
);

CREATE TABLE loans (
    loan_id SERIAL PRIMARY KEY,
    member_id INTEGER,
    book_id INTEGER,
    loan_date DATE NOT NULL,
    due_date DATE NOT NULL,
    return_date DATE,
    fine_amount NUMERIC DEFAULT 0,
    CONSTRAINT check_loan_dates CHECK (due_date > loan_date),
    FOREIGN KEY (member_id) REFERENCES members(member_id),
    FOREIGN KEY (book_id) REFERENCES books(book_id)
);

CREATE INDEX idx_book_title ON books(title);
CREATE INDEX idx_author_name ON authors(author_name);
CREATE INDEX idx_member_email ON members(email);
CREATE INDEX idx_loan_status ON loans(return_date);

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

INSERT INTO authors (author_name, nationality, birth_year)
VALUES ($1, $2, $3)
RETURNING author_id;

INSERT INTO books (title, isbn, publication_year, price, author_id, category, total_copies)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING book_id;

INSERT INTO members (full_name, email, join_date, membership_type, phone_number)
VALUES ($1, $2, $3, $4, $5)
RETURNING member_id;

INSERT INTO loans (member_id, book_id, loan_date, due_date, fine_amount)
VALUES ($1, $2, $3, $4, $5)
RETURNING loan_id;

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

