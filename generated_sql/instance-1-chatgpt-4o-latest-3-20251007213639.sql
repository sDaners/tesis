-- Create table: departments
CREATE TABLE departments (
  dept_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  dept_name STRING(50) NOT NULL,
  location STRING(100),
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (dept_id);

-- Create table: employees
CREATE TABLE employees (
  emp_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  first_name STRING(50) NOT NULL,
  last_name STRING(50) NOT NULL,
  email STRING(150),
  hire_date DATE NOT NULL,
  salary NUMERIC,
  dept_id STRING(36),
  manager_id STRING(36),
  phone_number STRING(20)
) PRIMARY KEY (emp_id);

-- Add foreign key constraints
ALTER TABLE employees ADD CONSTRAINT FK_employees_dept_id FOREIGN KEY (dept_id) REFERENCES departments (dept_id);
ALTER TABLE employees ADD CONSTRAINT FK_employees_manager_id FOREIGN KEY (manager_id) REFERENCES employees (emp_id);

-- Create index: unique email
CREATE UNIQUE INDEX idx_emp_email ON employees(email);

-- Create table: projects
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

-- Create table: project_assignments
CREATE TABLE project_assignments (
  emp_id STRING(36) NOT NULL,
  project_id STRING(36) NOT NULL,
  role STRING(50),
  hours_allocated INT64,
  CONSTRAINT FK_project_assignments_emp_id FOREIGN KEY (emp_id) REFERENCES employees (emp_id),
  CONSTRAINT FK_project_assignments_project_id FOREIGN KEY (project_id) REFERENCES projects (project_id)
) PRIMARY KEY (emp_id, project_id);

-- Indexes
CREATE INDEX idx_emp_name ON employees(last_name, first_name);
CREATE INDEX idx_dept_location ON departments(location);
CREATE INDEX idx_project_status ON projects(status);

-- View: employee_details
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

-- Inserts with generated UUIDs supplied from application
-- Replace @... with generated UUIDs at runtime

-- Insert into departments
INSERT INTO departments (dept_id, dept_name, location)
VALUES ('UUID-DEPT-001', 'Engineering', 'New York');

-- Insert into employees
INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id)
VALUES ('UUID-EMP-001', 'John', 'Doe', 'john.doe@example.com', DATE '2024-01-15', 75000, 'UUID-DEPT-001');

-- Insert into projects
INSERT INTO projects (project_id, project_name, start_date, end_date, budget, status)
VALUES ('UUID-PROJ-001', 'Apollo', DATE '2024-02-01', DATE '2024-12-31', 250000, 'ACTIVE');

-- Insert into project_assignments
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES ('UUID-EMP-001', 'UUID-PROJ-001', 'Developer', 160);

-- Select query
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

-- Drop view and indexes
DROP VIEW employee_details;
DROP INDEX idx_project_status;
DROP INDEX idx_dept_location;
DROP INDEX idx_emp_name;
DROP INDEX idx_emp_email;

-- Drop tables
DROP TABLE project_assignments;
DROP TABLE projects;
DROP TABLE employees;
DROP TABLE departments;