package repo

// PostgresDDL contains all the table creation statements for the PostgreSQL database.
// It creates the following structure:
// - departments: Stores department information with auto-incrementing IDs
// - employees: Stores employee information with references to departments and managers
// - projects: Stores project information with status and date constraints
// - project_assignments: Links employees to projects with their roles
// Also creates necessary indexes and a view for employee details.
const PostgresDDL = `
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

// InsertDepartmentSQL inserts a new department and returns its ID.
// Parameters:
//   - $1: Name of the department
//   - $2: Location of the department
//
// Returns: The generated department ID (SERIAL)
const InsertDepartmentSQL = `
INSERT INTO departments (dept_name, location)
VALUES ($1, $2)
RETURNING dept_id`

// InsertEmployeeSQL inserts a new employee and returns their ID.
// Parameters:
//   - $1: Employee's first name
//   - $2: Employee's last name
//   - $3: Employee's email address
//   - $4: Employee's hire date
//   - $5: Employee's salary
//   - $6: ID of the department the employee belongs to
//
// Returns: The generated employee ID (SERIAL)
const InsertEmployeeSQL = `
INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING emp_id`

// InsertProjectSQL inserts a new project and returns its ID.
// Parameters:
//   - $1: Name of the project
//   - $2: Project start date
//   - $3: Project end date
//   - $4: Project budget
//   - $5: Project status (must be one of: 'ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED')
//
// Returns: The generated project ID (SERIAL)
const InsertProjectSQL = `
INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING project_id`

// InsertProjectAssignmentSQL assigns an employee to a project with a specific role.
// Parameters:
//   - $1: ID of the employee
//   - $2: ID of the project
//   - $3: Role of the employee in the project
//   - $4: Number of hours allocated to the project
const InsertProjectAssignmentSQL = `
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES ($1, $2, $3, $4)`

// QueryEmployeeDetailsSQL retrieves detailed information about employees including:
// - Basic employee information (ID, name, email)
// - Department information
// - Manager information
// - Project information
// The query joins multiple tables to provide a comprehensive view of each employee's details.
const QueryEmployeeDetailsSQL = `
SELECT e.emp_id, e.first_name, e.last_name, e.email, d.dept_name, 
       m.first_name as manager_first_name, m.last_name as manager_last_name,
       p.project_name
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
LEFT JOIN project_assignments pa ON e.emp_id = pa.emp_id
LEFT JOIN projects p ON pa.project_id = p.project_id
`

// CleanupDBSQL contains the SQL statements to clean up the database.
// The statements are ordered to respect dependencies:
// 1. Drop views
// 2. Drop indexes
// 3. Drop tables in reverse order of their dependencies
const CleanupDBSQL = `--sql
DROP VIEW IF EXISTS employee_details;
DROP INDEX IF EXISTS idx_project_status;
DROP INDEX IF EXISTS idx_dept_location;
DROP INDEX IF EXISTS idx_emp_name;
DROP INDEX IF EXISTS idx_emp_email;
DROP TABLE IF EXISTS project_assignments;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;
`
