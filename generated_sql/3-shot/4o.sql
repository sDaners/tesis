--  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

-- Total statements: 23
-- Successfully parsed: 22
-- Parse errors: 1
-- Executed: 10
-- Execution errors: 12
-- Parse success rate: 95.7%
-- Execution success rate (of parsed): 45.5%
-- Overall success rate: 43.5%

-- Create tables

CREATE TABLE departments (
  dept_id INT64 NOT NULL DEFAULT (GENERATE_UUID()),
  dept_name STRING(50) NOT NULL,
  location STRING(100),
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP()),
) PRIMARY KEY(dept_id);

CREATE TABLE employees (
  emp_id INT64 NOT NULL DEFAULT (GENERATE_UUID()),
  first_name STRING(50) NOT NULL,
  last_name STRING(50) NOT NULL,
  email STRING(150),
  hire_date DATE NOT NULL,
  salary NUMERIC,
  dept_id INT64,
  manager_id INT64,
  phone_number STRING(20),
) PRIMARY KEY(emp_id),
  INTERLEAVE IN PARENT departments ON DELETE NO ACTION;

ALTER TABLE employees
  ADD CONSTRAINT fk_employees_dept
  FOREIGN KEY (dept_id)
  REFERENCES departments (dept_id);

ALTER TABLE employees
  ADD CONSTRAINT fk_employees_manager
  FOREIGN KEY (manager_id)
  REFERENCES employees (emp_id);

CREATE UNIQUE INDEX idx_emp_email ON employees(email);

CREATE TABLE projects (
  project_id INT64 NOT NULL DEFAULT (GENERATE_UUID()),
  project_name STRING(100) NOT NULL,
  start_date DATE,
  end_date DATE,
  budget NUMERIC,
  status STRING(20) DEFAULT ('ACTIVE'),
  CONSTRAINT check_dates CHECK (end_date > start_date),
  CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY(project_id);

CREATE TABLE project_assignments (
  emp_id INT64 NOT NULL,
  project_id INT64 NOT NULL,
  role STRING(50),
  hours_allocated INT64,
) PRIMARY KEY(emp_id, project_id);

ALTER TABLE project_assignments
  ADD CONSTRAINT fk_pa_emp
  FOREIGN KEY (emp_id)
  REFERENCES employees (emp_id);

ALTER TABLE project_assignments
  ADD CONSTRAINT fk_pa_project
  FOREIGN KEY (project_id)
  REFERENCES projects (project_id);

-- Create indexes

CREATE INDEX idx_emp_name ON employees(last_name, first_name);
CREATE INDEX idx_dept_location ON departments(location);
CREATE INDEX idx_project_status ON projects(status);

-- Create view

CREATE VIEW employee_details AS
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

-- Inserts with RETURNING

-- Note: Spanner does not support RETURNING clause directly. Instead, use `SELECT LAST_INSERT_ID()` logic in application if applicable.
-- Here's the adjusted insert format, application must fetch the ID using the generated value or a UUID beforehand.

-- Example (application should generate UUID and use it):
-- INSERT INTO departments (dept_id, dept_name, location) VALUES ('generated-id', '...', '...');

-- Read-only queries

SELECT e.emp_id, e.first_name, e.last_name, e.email, d.dept_name, 
       m.first_name AS manager_first_name, m.last_name AS manager_last_name,
       p.project_name
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
LEFT JOIN project_assignments pa ON e.emp_id = pa.emp_id
LEFT JOIN projects p ON pa.project_id = p.project_id;

-- Drop statements

DROP VIEW employee_details;
DROP INDEX idx_project_status;
DROP INDEX idx_dept_location;
DROP INDEX idx_emp_name;
DROP INDEX idx_emp_email;
DROP TABLE project_assignments;
DROP TABLE projects;
DROP TABLE employees;
DROP TABLE departments;
