CREATE TABLE departments (
  dept_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  dept_name STRING(50) NOT NULL,
  location STRING(100),
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
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
  FOREIGN KEY (dept_id) REFERENCES departments (dept_id),
  FOREIGN KEY (manager_id) REFERENCES employees (emp_id)
) PRIMARY KEY (emp_id);

CREATE UNIQUE INDEX idx_emp_email ON employees(email);

CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name STRING(100) NOT NULL,
  start_date DATE,
  end_date DATE,
  budget NUMERIC,
  status STRING(20) DEFAULT ('ACTIVE'),
  CONSTRAINT check_dates CHECK (end_date > start_date),
  CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id);

CREATE TABLE project_assignments (
  emp_id STRING(36) NOT NULL,
  project_id STRING(36) NOT NULL,
  role STRING(50),
  hours_allocated INT64,
  FOREIGN KEY (emp_id) REFERENCES employees (emp_id),
  FOREIGN KEY (project_id) REFERENCES projects (project_id)
) PRIMARY KEY (emp_id, project_id);

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

-- Note: Spanner does not support parameterized INSERTs using $1, $2, etc.
-- Replace with actual values or use client-side parameter binding in application code.

-- Example inserts with literal values (replace with actual values in application):
INSERT INTO departments (dept_id, dept_name, location)
VALUES ('<UUID>', 'Engineering', 'New York');

INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id)
VALUES ('<UUID>', 'John', 'Doe', 'john.doe@example.com', '2023-01-01', 75000, '<DEPT_UUID>');

INSERT INTO projects (project_id, project_name, start_date, end_date, budget, status)
VALUES ('<UUID>', 'Project Apollo', '2023-01-01', '2023-12-31', 1000000, 'ACTIVE');

INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES ('<EMP_UUID>', '<PROJECT_UUID>', 'Developer', 160);

SELECT 
  e.emp_id, 
  e.first_name, 
  e.last_name, 
  e.email, 
  d.dept_name, 
  m.first_name AS manager_first_name, 
  m.last_name AS manager_last_name,
  p.project_name
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
LEFT JOIN project_assignments pa ON e.emp_id = pa.emp_id
LEFT JOIN projects p ON pa.project_id = p.project_id;

DROP VIEW employee_details;
DROP INDEX idx_project_status;
DROP INDEX idx_dept_location;
DROP INDEX idx_emp_name;
DROP INDEX idx_emp_email;
DROP TABLE project_assignments;
DROP TABLE projects;
DROP TABLE employees;
DROP TABLE departments;