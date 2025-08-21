--  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

-- Total statements: 17
-- Successfully parsed: 17
-- Parse errors: 0
-- Executed: 16
-- Execution errors: 1
-- Parse success rate: 100.0%
-- Execution success rate (of parsed): 94.1%
-- Overall success rate: 94.1%

-- CHAT: https://chatgpt.com/c/684afdfe-2298-8003-82d7-c6bcb4b215f4

-- Table Creation (GoogleSQL doesn't support IF NOT EXISTS or SERIAL, 
-- so use INT64 with AUTO_INCREMENT if needed)

CREATE TABLE departments (
  dept_id INT64 NOT NULL,
  dept_name STRING(50) NOT NULL,
  location STRING(100),
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (dept_id);

CREATE TABLE employees (
  emp_id INT64 NOT NULL,
  first_name STRING(50) NOT NULL,
  last_name STRING(50) NOT NULL,
  email STRING(150),
  hire_date DATE NOT NULL,
  salary NUMERIC,
  dept_id INT64,
  manager_id INT64,
  phone_number STRING(20)
) PRIMARY KEY (emp_id);

CREATE TABLE projects (
  project_id INT64 NOT NULL,
  project_name STRING(100) NOT NULL,
  start_date DATE,
  end_date DATE,
  budget NUMERIC,
  status STRING(20)
) PRIMARY KEY (project_id);

CREATE TABLE project_assignments (
  emp_id INT64 NOT NULL,
  project_id INT64 NOT NULL,
  role STRING(50),
  hours_allocated INT64
) PRIMARY KEY (emp_id, project_id);

-- Foreign Keys (must be declared separately in Spanner)
ALTER TABLE employees ADD CONSTRAINT fk_employees_dept
  FOREIGN KEY (dept_id) REFERENCES departments (dept_id);

ALTER TABLE employees ADD CONSTRAINT fk_employees_manager
  FOREIGN KEY (manager_id) REFERENCES employees (emp_id);

ALTER TABLE project_assignments ADD CONSTRAINT fk_assignments_emp
  FOREIGN KEY (emp_id) REFERENCES employees (emp_id);

ALTER TABLE project_assignments ADD CONSTRAINT fk_assignments_project
  FOREIGN KEY (project_id) REFERENCES projects (project_id);

-- Check Constraints (Spanner only supports CHECK for IN and basic expressions)
ALTER TABLE projects ADD CONSTRAINT check_status 
  CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'));

ALTER TABLE projects ADD CONSTRAINT check_dates
  CHECK (end_date > start_date);

-- Indexes (IF NOT EXISTS is not supported in Spanner)
CREATE UNIQUE INDEX idx_emp_email ON employees(email);
CREATE INDEX idx_emp_name ON employees(last_name, first_name);
CREATE INDEX idx_dept_location ON departments(location);
CREATE INDEX idx_project_status ON projects(status);

CREATE OR REPLACE VIEW employee_details SQL SECURITY INVOKER AS
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

-- Insert Statements (Spanner does not support RETURNING)
-- Use this pattern in your application logic instead:
INSERT INTO departments (dept_id, dept_name, location)
VALUES (@dept_id, @dept_name, @location);

-- or with generated ID in your application, then reuse that ID.

-- Complex SELECT
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

-- Drop Statements (IF EXISTS not supported)
-- Spanner doesn't support DROP VIEW or DROP INDEX directly via DDL. 
-- You must delete objects manually from the schema management tool (like gcloud or Terraform) 
-- or use schema migration tools like Liquibase or Flyway.

-- Notes:
-- Replace SERIAL with INT64 and manage auto-increment IDs at application level 
-- or use GENERATED IDs if supported by your ORM.

-- - RETURNING is not supported in Spanner.

-- - DROP IF EXISTS, IF NOT EXISTS are not supported in GoogleSQL (Spanner).

-- - Default CURRENT_TIMESTAMP must use CURRENT_TIMESTAMP() function in Spanner.

-- - Check constraints are limited in expressiveness.

-- Let me know if you'd like me to generate DML inserts with GENERATE_UUID() 
-- for primary keys instead of integers. 