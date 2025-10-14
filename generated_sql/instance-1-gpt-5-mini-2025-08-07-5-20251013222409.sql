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
  CONSTRAINT fk_emp_dept FOREIGN KEY (dept_id) REFERENCES departments(dept_id),
  CONSTRAINT fk_emp_manager FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
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
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id;

-- Insert a department (application supplies @dept_name, @location)
INSERT INTO departments (dept_id, dept_name, location)
VALUES (GENERATE_UUID(), @dept_name, @location);

-- Retrieve the dept_id just inserted (application should run this immediately after insert)
SELECT dept_id FROM departments
WHERE dept_name = @dept_name AND location = @location
ORDER BY created_at DESC LIMIT 1;

-- Insert an employee (application supplies @first_name, @last_name, @email, @hire_date (DATE), @salary, @dept_id)
INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id)
VALUES (GENERATE_UUID(), @first_name, @last_name, @email, @hire_date, @salary, @dept_id);

-- Retrieve the emp_id just inserted
SELECT emp_id FROM employees
WHERE first_name = @first_name AND last_name = @last_name AND email = @email
LIMIT 1;

-- Insert a project (application supplies @project_name, @start_date (DATE), @end_date (DATE), @budget, @status)
INSERT INTO projects (project_id, project_name, start_date, end_date, budget, status)
VALUES (GENERATE_UUID(), @project_name, @start_date, @end_date, @budget, @status);

-- Retrieve the project_id just inserted
SELECT project_id FROM projects
WHERE project_name = @project_name AND start_date = @start_date
LIMIT 1;

-- Insert a project assignment.
-- IMPORTANT: Provide valid existing emp_id and project_id values (e.g., from the SELECTs above).
-- Do NOT rely on generated defaults for foreign key columns, otherwise FK violations will occur.
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES (@emp_id, @project_id, @role, @hours_allocated);

-- Complex SELECT joining employees, departments, managers, project assignments and projects
SELECT e.emp_id, e.first_name, e.last_name, e.email, d.dept_name,
       m.first_name AS manager_first_name, m.last_name AS manager_last_name,
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