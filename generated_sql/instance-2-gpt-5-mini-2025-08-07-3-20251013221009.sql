CREATE TABLE departments (
    dept_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
    dept_name STRING(50) NOT NULL,
    location STRING(100),
    created_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (dept_id);

CREATE TABLE employees (
    emp_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
    first_name STRING(50) NOT NULL,
    last_name STRING(50) NOT NULL,
    email STRING(150),
    hire_date DATE NOT NULL,
    salary NUMERIC,
    dept_id STRING(36),
    manager_id STRING(36),
    phone_number STRING(20),
    CONSTRAINT fk_employees_dept FOREIGN KEY (dept_id) REFERENCES departments(dept_id),
    CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
) PRIMARY KEY (emp_id);

CREATE TABLE projects (
    project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
    project_name STRING(100) NOT NULL,
    start_date DATE,
    end_date DATE,
    budget NUMERIC,
    status STRING(20) NOT NULL DEFAULT ('ACTIVE'),
    CONSTRAINT chk_projects_dates CHECK (start_date IS NULL OR end_date IS NULL OR end_date > start_date),
    CONSTRAINT chk_projects_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id);

CREATE TABLE project_assignments (
    emp_id STRING(36) NOT NULL,
    project_id STRING(36) NOT NULL,
    role STRING(50),
    hours_allocated INT64,
    CONSTRAINT fk_pa_emp FOREIGN KEY (emp_id) REFERENCES employees(emp_id),
    CONSTRAINT fk_pa_project FOREIGN KEY (project_id) REFERENCES projects(project_id)
) PRIMARY KEY (emp_id, project_id);

CREATE UNIQUE INDEX idx_emp_email ON employees(email);
CREATE INDEX idx_emp_name ON employees(last_name, first_name);
CREATE INDEX idx_dept_location ON departments(location);
CREATE INDEX idx_project_status ON projects(status);

CREATE VIEW employee_details SQL SECURITY INVOKER AS
SELECT 
    e.emp_id,
    e.first_name,
    e.last_name,
    e.email,
    d.dept_name,
    m.first_name AS manager_first_name,
    m.last_name AS manager_last_name
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id;

-- Insert a department (use parameters)
INSERT INTO departments (dept_name, location)
VALUES (@p1, @p2);

-- Retrieve the dept_id for the inserted department (Spanner has no RETURNING)
SELECT dept_id
FROM departments
WHERE dept_name = @p1 AND location = @p2
ORDER BY created_at DESC
LIMIT 1;

-- Insert an employee
INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES (@p1, @p2, @p3, @p4, @p5, @p6);

-- Retrieve the emp_id for the inserted employee
SELECT emp_id
FROM employees
WHERE first_name = @p1 AND last_name = @p2 AND email = @p3 AND hire_date = @p4
ORDER BY emp_id DESC
LIMIT 1;

-- Insert a project
INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES (@p1, @p2, @p3, @p4, @p5);

-- Retrieve the project_id for the inserted project
SELECT project_id
FROM projects
WHERE project_name = @p1 AND start_date = @p2
ORDER BY project_id DESC
LIMIT 1;

-- Assign an employee to a project
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES (@p1, @p2, @p3, @p4);

-- Example join query across tables
SELECT e.emp_id, e.first_name, e.last_name, e.email, d.dept_name, 
       m.first_name AS manager_first_name, m.last_name AS manager_last_name,
       p.project_name
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
LEFT JOIN project_assignments pa ON e.emp_id = pa.emp_id
LEFT JOIN projects p ON pa.project_id = p.project_id;

DROP VIEW IF EXISTS employee_details;
DROP INDEX IF EXISTS idx_project_status;
DROP INDEX IF EXISTS idx_dept_location;
DROP INDEX IF EXISTS idx_emp_name;
DROP INDEX IF EXISTS idx_emp_email;
DROP TABLE IF EXISTS project_assignments;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;