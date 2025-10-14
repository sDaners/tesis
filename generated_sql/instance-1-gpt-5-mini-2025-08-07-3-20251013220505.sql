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
  CONSTRAINT check_dates CHECK (end_date > start_date),
  CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
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
FROM employees AS e
LEFT JOIN departments AS d ON e.dept_id = d.dept_id
LEFT JOIN employees AS m ON e.manager_id = m.emp_id;

-- Insert a new department (dept_name is NOT NULL). This will only insert when @dept_name is provided.
INSERT INTO departments (dept_id, dept_name, location, created_at)
SELECT GENERATE_UUID(), CAST(@dept_name AS STRING), CAST(@location AS STRING), (CURRENT_TIMESTAMP())
WHERE @dept_name IS NOT NULL;

-- Insert a new employee (first_name, last_name, hire_date are NOT NULL). This will only insert when required params are provided.
INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id)
SELECT GENERATE_UUID(),
       CAST(@first_name AS STRING),
       CAST(@last_name AS STRING),
       CAST(@email AS STRING),
       CAST(@hire_date AS DATE),
       CAST(@salary AS NUMERIC),
       CAST(@dept_id AS STRING)
WHERE @first_name IS NOT NULL AND @last_name IS NOT NULL AND @hire_date IS NOT NULL;

-- Insert a new project (project_name is NOT NULL). This will only insert when project_name is provided.
INSERT INTO projects (project_id, project_name, start_date, end_date, budget, status)
SELECT GENERATE_UUID(),
       CAST(@project_name AS STRING),
       CAST(@start_date AS DATE),
       CAST(@end_date AS DATE),
       CAST(@budget AS NUMERIC),
       CAST(@status AS STRING)
WHERE @project_name IS NOT NULL;

-- Insert a project assignment (emp_id and project_id are part of the primary key and required).
-- This will only insert when both @emp_id and @project_id are provided (prevents NULL primary-key inserts).
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
SELECT CAST(@emp_id AS STRING),
       CAST(@project_id AS STRING),
       CAST(@role AS STRING),
       CAST(@hours_allocated AS INT64)
WHERE @emp_id IS NOT NULL AND @project_id IS NOT NULL;

SELECT e.emp_id, e.first_name, e.last_name, e.email, d.dept_name,
       m.first_name AS manager_first_name, m.last_name AS manager_last_name,
       p.project_name
FROM employees AS e
LEFT JOIN departments AS d ON e.dept_id = d.dept_id
LEFT JOIN employees AS m ON e.manager_id = m.emp_id
LEFT JOIN project_assignments AS pa ON e.emp_id = pa.emp_id
LEFT JOIN projects AS p ON pa.project_id = p.project_id;

DROP VIEW employee_details;
DROP INDEX idx_project_status;
DROP INDEX idx_dept_location;
DROP INDEX idx_emp_name;
DROP INDEX idx_emp_email;
DROP TABLE project_assignments;
DROP TABLE projects;
DROP TABLE employees;
DROP TABLE departments;