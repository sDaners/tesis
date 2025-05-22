package repo

import (
	"context"
	"database/sql"
	"time"

	"postgres-example/models"

	"github.com/georgysavva/scany/v2/sqlscan"
)

const spannerDDL = `
CREATE TABLE IF NOT EXISTS departments (
    dept_id INT64 NOT NULL,
    dept_name STRING(50) NOT NULL,
    location STRING(100),
    created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP()),
) PRIMARY KEY (dept_id);

CREATE TABLE IF NOT EXISTS employees (
    emp_id INT64 NOT NULL,
    first_name STRING(50) NOT NULL,
    last_name STRING(50) NOT NULL,
    email STRING(150),
    hire_date DATE NOT NULL,
    salary NUMERIC,
    dept_id INT64,
    manager_id INT64,
    phone_number STRING(20),
    CONSTRAINT fk_dept FOREIGN KEY (dept_id) REFERENCES departments(dept_id),
    CONSTRAINT fk_manager FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
) PRIMARY KEY (emp_id);

CREATE UNIQUE INDEX idx_emp_email ON employees(email);

CREATE TABLE IF NOT EXISTS projects (
    project_id INT64 NOT NULL,
    project_name STRING(100) NOT NULL,
    start_date DATE,
    end_date DATE,
    budget NUMERIC,
    status STRING(20) DEFAULT ('ACTIVE'),
    CONSTRAINT check_dates CHECK (end_date > start_date),
    CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id);

CREATE TABLE IF NOT EXISTS project_assignments (
    emp_id INT64 NOT NULL,
    project_id INT64 NOT NULL,
    role STRING(50),
    hours_allocated INT64,
    CONSTRAINT fk_emp FOREIGN KEY (emp_id) REFERENCES employees(emp_id),
    CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects(project_id)
) PRIMARY KEY (emp_id, project_id);

CREATE INDEX idx_emp_name ON employees(last_name, first_name);
CREATE INDEX idx_dept_location ON departments(location);
CREATE INDEX idx_project_status ON projects(status);

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

type SpannerRepo struct {
	DB *sql.DB
}

func NewSpannerRepo(db *sql.DB) *SpannerRepo {
	return &SpannerRepo{DB: db}
}

func (r *SpannerRepo) CreateTables() error {
	_, err := r.DB.Exec(spannerDDL)
	return err
}

func (r *SpannerRepo) InsertSampleData() (int, int, int, error) {
	var deptID int64
	err := r.DB.QueryRow(`
		INSERT INTO departments (dept_id, dept_name, location)
		VALUES (GENERATE_UUID(), @dept_name, @location)
		THEN RETURN dept_id`,
		sql.Named("dept_name", "Engineering"),
		sql.Named("location", "New York")).Scan(&deptID)
	if err != nil {
		return 0, 0, 0, err
	}

	var empID int64
	err = r.DB.QueryRow(`
		INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id)
		VALUES (GENERATE_UUID(), @first_name, @last_name, @email, @hire_date, @salary, @dept_id)
		THEN RETURN emp_id`,
		sql.Named("first_name", "John"),
		sql.Named("last_name", "Doe"),
		sql.Named("email", "john.doe@example.com"),
		sql.Named("hire_date", time.Now()),
		sql.Named("salary", 75000.00),
		sql.Named("dept_id", deptID)).Scan(&empID)
	if err != nil {
		return 0, 0, 0, err
	}

	var projectID int64
	err = r.DB.QueryRow(`
		INSERT INTO projects (project_id, project_name, start_date, end_date, budget, status)
		VALUES (GENERATE_UUID(), @project_name, @start_date, @end_date, @budget, @status)
		THEN RETURN project_id`,
		sql.Named("project_name", "Database Migration"),
		sql.Named("start_date", time.Now()),
		sql.Named("end_date", time.Now().AddDate(0, 3, 0)),
		sql.Named("budget", 50000.00),
		sql.Named("status", "ACTIVE")).Scan(&projectID)
	if err != nil {
		return 0, 0, 0, err
	}

	_, err = r.DB.Exec(`
		INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
		VALUES (@emp_id, @project_id, @role, @hours)`,
		sql.Named("emp_id", empID),
		sql.Named("project_id", projectID),
		sql.Named("role", "Lead Developer"),
		sql.Named("hours", 40))
	if err != nil {
		return 0, 0, 0, err
	}

	return int(deptID), int(empID), int(projectID), nil
}

func (r *SpannerRepo) QueryEmployeeDetails() ([]models.EmployeeDetails, error) {
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

func (r *SpannerRepo) CleanupDB() error {
	// Drop objects in reverse order of dependencies
	statements := []string{
		"DROP VIEW IF EXISTS employee_details",
		"DROP INDEX IF EXISTS idx_project_status",
		"DROP INDEX IF EXISTS idx_dept_location",
		"DROP INDEX IF EXISTS idx_emp_name",
		"DROP INDEX IF EXISTS idx_emp_email",
		"DROP TABLE IF EXISTS project_assignments",
		"DROP TABLE IF EXISTS projects",
		"DROP TABLE IF EXISTS employees",
		"DROP TABLE IF EXISTS departments",
	}

	for _, stmt := range statements {
		_, err := r.DB.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
