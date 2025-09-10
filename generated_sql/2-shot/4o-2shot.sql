--  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
-- 2 shot
-- Total statements: 27
-- Successfully parsed: 27
-- Parse errors: 0
-- Executed: 23
-- Execution errors: 4
-- Parse success rate: 100.0%
-- Execution success rate (of parsed): 85.2%
-- Overall success rate: 85.2%

--  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
-- 1 shot
-- Total statements: 23
-- Successfully parsed: 19
-- Parse errors: 4
-- Executed: 11
-- Execution errors: 8
-- Parse success rate: 82.6%
-- Execution success rate (of parsed): 57.9%
-- Overall success rate: 47.8%

-- https://chatgpt.com/share/68afc20e-8acc-8003-8fe5-e54074d9ad33

-- TABLE: departments
CREATE TABLE departments (
  dept_id INT64 NOT NULL,
  dept_name STRING(50) NOT NULL,
  location STRING(100),
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (dept_id);

-- TABLE: employees
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

-- TABLE: projects
CREATE TABLE projects (
  project_id INT64 NOT NULL,
  project_name STRING(100) NOT NULL,
  start_date DATE,
  end_date DATE,
  budget NUMERIC,
  status STRING(20) DEFAULT ('ACTIVE')
) PRIMARY KEY (project_id);

-- TABLE: project_assignments
CREATE TABLE project_assignments (
  emp_id INT64 NOT NULL,
  project_id INT64 NOT NULL,
  role STRING(50),
  hours_allocated INT64
) PRIMARY KEY (emp_id, project_id);

-- FOREIGN KEYS (added separately because Spanner does not support inline constraints)
ALTER TABLE employees ADD CONSTRAINT fk_employees_dept
  FOREIGN KEY (dept_id) REFERENCES departments (dept_id);

ALTER TABLE employees ADD CONSTRAINT fk_employees_manager
  FOREIGN KEY (manager_id) REFERENCES employees (emp_id);

ALTER TABLE project_assignments ADD CONSTRAINT fk_pa_emp
  FOREIGN KEY (emp_id) REFERENCES employees (emp_id);

ALTER TABLE project_assignments ADD CONSTRAINT fk_pa_project
  FOREIGN KEY (project_id) REFERENCES projects (project_id);

-- INDEXES
CREATE UNIQUE INDEX idx_emp_email ON employees(email);
CREATE INDEX idx_emp_name ON employees(last_name, first_name);
CREATE INDEX idx_dept_location ON departments(location);
CREATE INDEX idx_project_status ON projects(status);

-- VIEW: employee_details
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

-- INSERTS (Spanner does not support RETURNING — do a SELECT after insert)
-- departments insert
-- Parameters: @dept_id, @dept_name, @location
INSERT INTO departments (dept_id, dept_name, location, created_at)
VALUES (@dept_id, @dept_name, @location, CURRENT_TIMESTAMP());

-- employees insert
-- Parameters: @emp_id, @first_name, @last_name, @email, @hire_date, @salary, @dept_id
INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id)
VALUES (@emp_id, @first_name, @last_name, @email, @hire_date, @salary, @dept_id);

-- projects insert
-- Parameters: @project_id, @project_name, @start_date, @end_date, @budget, @status
INSERT INTO projects (project_id, project_name, start_date, end_date, budget, status)
VALUES (@project_id, @project_name, @start_date, @end_date, @budget, @status);

-- project_assignments insert
-- Parameters: @emp_id, @project_id, @role, @hours_allocated
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES (@emp_id, @project_id, @role, @hours_allocated);

-- SELECT with joins
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

-- CLEANUP
DROP VIEW employee_details;
DROP INDEX idx_project_status;
DROP INDEX idx_dept_location;
DROP INDEX idx_emp_name;
DROP INDEX idx_emp_email;
DROP TABLE project_assignments;
DROP TABLE projects;
DROP TABLE employees;
DROP TABLE departments;