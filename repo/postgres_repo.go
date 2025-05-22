package repo

import (
	"context"
	"database/sql"
	"time"

	"postgres-example/models"

	"github.com/georgysavva/scany/v2/sqlscan"
)

const ddl = `
CREATE TABLE IF NOT EXISTS departments (
    dept_id SERIAL PRIMARY KEY,
    dept_name VARCHAR(50) NOT NULL,
    location VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employees (
    emp_id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(150),
    hire_date DATE NOT NULL,
    salary NUMERIC,
    dept_id INTEGER REFERENCES departments(dept_id),
    manager_id INTEGER REFERENCES employees(emp_id),
    phone_number VARCHAR(20)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_emp_email ON employees(email);

CREATE TABLE IF NOT EXISTS projects (
    project_id SERIAL PRIMARY KEY,
    project_name VARCHAR(100) NOT NULL,
    start_date DATE,
    end_date DATE,
    budget NUMERIC,
    status VARCHAR(20) DEFAULT 'ACTIVE',
    CONSTRAINT check_dates CHECK (end_date > start_date),
    CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
);

CREATE TABLE IF NOT EXISTS project_assignments (
    emp_id INTEGER,
    project_id INTEGER,
    role VARCHAR(50),
    hours_allocated INTEGER,
    PRIMARY KEY (emp_id, project_id),
    FOREIGN KEY (emp_id) REFERENCES employees(emp_id),
    FOREIGN KEY (project_id) REFERENCES projects(project_id)
);

CREATE INDEX IF NOT EXISTS idx_emp_name ON employees(last_name, first_name);
CREATE INDEX IF NOT EXISTS idx_dept_location ON departments(location);
CREATE INDEX IF NOT EXISTS idx_project_status ON projects(status);

CREATE OR REPLACE VIEW employee_details AS
SELECT 
    e.emp_id,
    e.first_name,
    e.last_name,
    e.email,
    d.dept_name,
    m.first_name as manager_first_name,
    m.last_name as manager_last_name
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id;
`

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{DB: db}
}

func (r *PostgresRepo) CreateTables() error {
	_, err := r.DB.Exec(ddl)
	return err
}

func (r *PostgresRepo) InsertSampleData() (int, int, int, error) {
	var deptID int
	err := r.DB.QueryRow(`
		INSERT INTO departments (dept_name, location)
		VALUES ($1, $2)
		RETURNING dept_id`, "Engineering", "New York").Scan(&deptID)
	if err != nil {
		return 0, 0, 0, err
	}

	var empID int
	err = r.DB.QueryRow(`
		INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING emp_id`,
		"John", "Doe", "john.doe@example.com", time.Now(), 75000.00, deptID).Scan(&empID)
	if err != nil {
		return 0, 0, 0, err
	}

	var projectID int
	err = r.DB.QueryRow(`
		INSERT INTO projects (project_name, start_date, end_date, budget, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING project_id`,
		"Database Migration", time.Now(), time.Now().AddDate(0, 3, 0), 50000.00, "ACTIVE").Scan(&projectID)
	if err != nil {
		return 0, 0, 0, err
	}

	_, err = r.DB.Exec(`
		INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
		VALUES ($1, $2, $3, $4)`,
		empID, projectID, "Lead Developer", 40)
	if err != nil {
		return 0, 0, 0, err
	}

	return deptID, empID, projectID, nil
}

func (r *PostgresRepo) QueryEmployeeDetails() ([]models.EmployeeDetails, error) {
	query := `
		SELECT e.emp_id, e.first_name, e.last_name, e.email, d.dept_name, 
		       m.first_name as manager_first_name, m.last_name as manager_last_name,
		   p.project_name
		FROM employees e
		LEFT JOIN departments d ON e.dept_id = d.dept_id
		LEFT JOIN employees m ON e.manager_id = m.emp_id
		LEFT JOIN project_assignments pa ON e.emp_id = pa.emp_id
		LEFT JOIN projects p ON pa.project_id = p.project_id
	`

	var details []models.EmployeeDetails
	err := sqlscan.Select(context.Background(), r.DB, &details, query)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (r *PostgresRepo) CleanupDB() error {
	_, err := r.DB.Exec(`
		DROP VIEW IF EXISTS employee_details;
		DROP INDEX IF EXISTS idx_project_status;
		DROP INDEX IF EXISTS idx_dept_location;
		DROP INDEX IF EXISTS idx_emp_name;
		DROP INDEX IF EXISTS idx_emp_email;
		DROP TABLE IF EXISTS project_assignments;
		DROP TABLE IF EXISTS projects;
		DROP TABLE IF EXISTS employees;
		DROP TABLE IF EXISTS departments;
	`)
	return err
}
