package valid_repo

// SpannerDDL contains all the table creation statements for the Spanner database.
// It creates the following structure:
// - departments: Stores department information
// - employees: Stores employee information with references to departments and managers
// - projects: Stores project information with status and date constraints
// - project_assignments: Links employees to projects with their roles
// Also creates necessary indexes and a view for employee details.
var SpannerDDL = []string{
	`ALTER DATABASE db SET OPTIONS (default_sequence_kind = 'bit_reversed_positive')`,

	`CREATE TABLE IF NOT EXISTS departments (
    dept_id INT64 NOT NULL AUTO_INCREMENT,
    dept_name STRING(50) NOT NULL,
    location STRING(100),
    created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
    ) PRIMARY KEY (dept_id)`,

	`CREATE TABLE IF NOT EXISTS employees (
    emp_id INT64 NOT NULL AUTO_INCREMENT,
    first_name STRING(50) NOT NULL,
    last_name STRING(50) NOT NULL,
    email STRING(150),
    hire_date TIMESTAMP NOT NULL,
    salary FLOAT64,
    dept_id INT64,
    manager_id INT64,
    phone_number STRING(20),
    CONSTRAINT fk_dept FOREIGN KEY (dept_id) REFERENCES departments(dept_id),
    CONSTRAINT fk_manager FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
    ) PRIMARY KEY (emp_id)`,

	`CREATE UNIQUE INDEX idx_emp_email ON employees(email)`,

	`CREATE TABLE IF NOT EXISTS projects (
    project_id INT64 NOT NULL AUTO_INCREMENT,
    project_name STRING(100) NOT NULL,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    budget FLOAT64,
    status STRING(20) DEFAULT ('ACTIVE'),
    CONSTRAINT check_dates CHECK (end_date > start_date),
    CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
    ) PRIMARY KEY (project_id)`,

	`CREATE TABLE IF NOT EXISTS project_assignments (
    emp_id INT64 NOT NULL,
    project_id INT64 NOT NULL,
    role STRING(50),
    hours_allocated INT64,
    CONSTRAINT fk_emp FOREIGN KEY (emp_id) REFERENCES employees(emp_id),
    CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects(project_id)
    ) PRIMARY KEY (emp_id, project_id)`,

	`CREATE INDEX idx_emp_name ON employees(last_name, first_name)`,
	`CREATE INDEX idx_dept_location ON departments(location)`,
	`CREATE INDEX idx_project_status ON projects(status)`,

	`CREATE OR REPLACE VIEW employee_details
    SQL SECURITY INVOKER
    AS SELECT 
        e.emp_id,
        e.first_name,
        e.last_name,
        e.email,
        d.dept_name,
        m.first_name as manager_first_name,
        m.last_name as manager_last_name
    FROM employees e
    LEFT JOIN departments d ON e.dept_id = d.dept_id
    LEFT JOIN employees m ON e.manager_id = m.emp_id
    `,
}

// SpannerInsertDepartmentSQL inserts a new department and returns its ID.
// Parameters:
//   - @dept_name: Name of the department
//   - @location: Location of the department
//
// Returns: The generated department ID
const SpannerInsertDepartmentSQL = `
INSERT INTO departments (dept_name, location)
VALUES (@dept_name, @location)
THEN RETURN dept_id`

// SpannerInsertEmployeeSQL inserts a new employee and returns their ID.
// Parameters:
//   - @first_name: Employee's first name
//   - @last_name: Employee's last name
//   - @email: Employee's email address
//   - @hire_date: Employee's hire date
//   - @salary: Employee's salary
//   - @dept_id: ID of the department the employee belongs to
//
// Returns: The generated employee ID
const SpannerInsertEmployeeSQL = `
INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES (@first_name, @last_name, @email, @hire_date, @salary, @dept_id)
THEN RETURN emp_id`

// SpannerInsertProjectSQL inserts a new project and returns its ID.
// Parameters:
//   - @project_name: Name of the project
//   - @start_date: Project start date
//   - @end_date: Project end date
//   - @budget: Project budget
//   - @status: Project status (must be one of: 'ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED')
//
// Returns: The generated project ID
const SpannerInsertProjectSQL = `
INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES (@project_name, @start_date, @end_date, @budget, @status)
THEN RETURN project_id`

// SpannerInsertProjectAssignmentSQL assigns an employee to a project with a specific role.
// Parameters:
//   - @emp_id: ID of the employee
//   - @project_id: ID of the project
//   - @role: Role of the employee in the project
//   - @hours: Number of hours allocated to the project
const SpannerInsertProjectAssignmentSQL = `
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES (@emp_id, @project_id, @role, @hours)`

// SpannerQueryEmployeeDetailsSQL retrieves detailed information about employees including:
// - Basic employee information (ID, name, email)
// - Department information
// - Manager information
// - Project information
// The query joins multiple tables to provide a comprehensive view of each employee's details.
const SpannerQueryEmployeeDetailsSQL = `
SELECT e.emp_id, e.first_name, e.last_name, e.email, d.dept_name, 
       m.first_name as manager_first_name, m.last_name as manager_last_name,
       p.project_name
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
LEFT JOIN project_assignments pa ON e.emp_id = pa.emp_id
LEFT JOIN projects p ON pa.project_id = p.project_id
`

// SpannerCleanupStatements contains the SQL statements to clean up the database.
// The statements are ordered to respect dependencies:
// 1. Drop views
// 2. Drop indexes
// 3. Drop tables in reverse order of their dependencies
var SpannerCleanupStatements = []string{
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
