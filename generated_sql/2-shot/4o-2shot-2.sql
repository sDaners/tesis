-- Second shot
-- Total statements: 19
-- Successfully parsed: 18
-- Parse errors: 1
-- Executed: 15
-- Execution errors: 3
-- Parse success rate: 94.7%
-- Execution success rate (of parsed): 83.3%
-- Overall success rate: 78.9%

--  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
-- First shot
-- Total statements: 19
-- Successfully parsed: 15
-- Parse errors: 4
-- Executed: 5
-- Execution errors: 10
-- Parse success rate: 78.9%
-- Execution success rate (of parsed): 33.3%
-- Overall success rate: 26.3%

-- Create departments table first (no dependencies)
CREATE TABLE departments (
  dept_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  dept_name STRING(50) NOT NULL,
  location STRING(100),
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (dept_id);

-- Create employees table (depends on departments and itself)
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
  CONSTRAINT fk_dept FOREIGN KEY (dept_id) REFERENCES departments (dept_id),
  CONSTRAINT fk_manager FOREIGN KEY (manager_id) REFERENCES employees (emp_id)
) PRIMARY KEY (emp_id);

-- Create projects table (independent)
CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name STRING(100) NOT NULL,
  start_date DATE,
  end_date DATE,
  budget NUMERIC,
  status STRING(20) DEFAULT ('ACTIVE')
) PRIMARY KEY (project_id);

-- Create project_assignments table (depends on employees and projects)
CREATE TABLE project_assignments (
  emp_id STRING(36) NOT NULL,
  project_id STRING(36) NOT NULL,
  role STRING(50),
  hours_allocated INT64,
  PRIMARY KEY (emp_id, project_id),
  CONSTRAINT fk_assignment_emp FOREIGN KEY (emp_id) REFERENCES employees (emp_id),
  CONSTRAINT fk_assignment_project FOREIGN KEY (project_id) REFERENCES projects (project_id)
);

-- Indexes
CREATE UNIQUE INDEX idx_emp_email ON employees(email);
CREATE INDEX idx_emp_name ON employees(last_name, first_name);
CREATE INDEX idx_dept_location ON departments(location);
CREATE INDEX idx_project_status ON projects(status);

-- View
CREATE OR REPLACE VIEW employee_details
SQL SECURITY INVOKER
AS
SELECT 
  e.emp_id,
  e.first_name,
  e.last_name,
  e.email,
  d.dept_name,
  m.first_name AS manager_first_name,
  m.last_name AS manager_last_name
FROM employees AS e
LEFT JOIN departments AS d ON e.dept_id = d.dept_id
LEFT JOIN employees AS m ON e.manager_id = m.emp_id;

-- Inserts (replace @params with client-side binding as needed)
INSERT INTO departments (dept_name, location)
VALUES ('Engineering', 'New York');

INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES ('John', 'Doe', 'john.doe@example.com', '2024-01-01', 85000.00, 'REPLACE_WITH_DEPT_ID');

INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES ('Project Phoenix', '2024-01-10', '2024-12-31', 1000000.00, 'ACTIVE');

INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES ('REPLACE_WITH_EMP_ID', 'REPLACE_WITH_PROJECT_ID', 'Lead Engineer', 100);

-- Select with joins
SELECT 
  e.emp_id,
  e.first_name,
  e.last_name,
  e.email,
  d.dept_name,
  m.first_name AS manager_first_name,
  m.last_name AS manager_last_name,
  p.project_name
FROM employees AS e
LEFT JOIN departments AS d ON e.dept_id = d.dept_id
LEFT JOIN employees AS m ON e.manager_id = m.emp_id
LEFT JOIN project_assignments AS pa ON e.emp_id = pa.emp_id
LEFT JOIN projects AS p ON pa.project_id = p.project_id;

-- Drop statements
DROP VIEW employee_details;
DROP TABLE project_assignments;
DROP TABLE projects;
DROP TABLE employees;
DROP TABLE departments;
