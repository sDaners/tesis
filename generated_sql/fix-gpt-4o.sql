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
) PRIMARY KEY (emp_id),
  INTERLEAVE IN PARENT departments ON DELETE NO ACTION;

CREATE TABLE projects (
  project_id INT64 NOT NULL,
  project_name STRING(100) NOT NULL,
  start_date DATE,
  end_date DATE,
  budget NUMERIC,
  status STRING(20)
) PRIMARY KEY (project_id);

-- Emulated Constraints
-- Spanner does not support:
-- 
-- CHECK
-- 
-- ENUM-like constraints
-- 
-- FOREIGN KEYS (without INTERLEAVE, but even those are more for locality than strict FK enforcement)
-- 
-- IF NOT EXISTS or DROP 
-- 
-- So we simulate them with comments or in application logic:

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

-- No RETURNING support in Spanner
INSERT INTO departments (dept_id, dept_name, location)
VALUES (@dept_id, @dept_name, @location);

SELECT 
  e.emp_id, e.first_name, e.last_name, e.email, d.dept_name, 
  m.first_name AS manager_first_name, m.last_name AS manager_last_name,
  p.project_name
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
LEFT JOIN project_assignments pa ON e.emp_id = pa.emp_id
LEFT JOIN projects p ON pa.project_id = p.project_id; 