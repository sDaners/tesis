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

-- Example inserts with concrete values to avoid missing-parameter errors and satisfy NOT NULL / FK constraints.

-- Insert a department
INSERT INTO departments (dept_id, dept_name, location, created_at)
VALUES (GENERATE_UUID(), 'Engineering', 'New York', (CURRENT_TIMESTAMP()));

-- Insert a manager employee (no manager_id)
INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id)
VALUES (
  GENERATE_UUID(),
  'Bob',
  'Manager',
  'bob.manager@example.com',
  DATE '2018-06-01',
  120000,
  (SELECT dept_id FROM departments WHERE dept_name = 'Engineering' AND location = 'New York' LIMIT 1)
);

-- Insert a regular employee with manager set to the previously inserted manager
INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id, manager_id)
VALUES (
  GENERATE_UUID(),
  'Alice',
  'Smith',
  'alice.smith@example.com',
  DATE '2020-01-15',
  90000,
  (SELECT dept_id FROM departments WHERE dept_name = 'Engineering' AND location = 'New York' LIMIT 1),
  (SELECT emp_id FROM employees WHERE email = 'bob.manager@example.com' LIMIT 1)
);

-- Insert a project (valid dates satisfying CHECK)
INSERT INTO projects (project_id, project_name, start_date, end_date, budget, status)
VALUES (
  GENERATE_UUID(),
  'Platform Migration',
  DATE '2024-01-01',
  DATE '2024-12-31',
  500000,
  'ACTIVE'
);

-- Assign the employee to the project (ensure emp_id and project_id are selected from existing rows)
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES (
  (SELECT emp_id FROM employees WHERE email = 'alice.smith@example.com' LIMIT 1),
  (SELECT project_id FROM projects WHERE project_name = 'Platform Migration' LIMIT 1),
  'Developer',
  1200
);

-- Query joining employees, departments, managers, assignments and projects
SELECT e.emp_id,
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

-- Cleanup: drop view, indexes, and tables in dependency-aware order
DROP VIEW employee_details;
DROP INDEX idx_project_status;
DROP INDEX idx_dept_location;
DROP INDEX idx_emp_name;
DROP INDEX idx_emp_email;
DROP TABLE project_assignments;
DROP TABLE projects;
DROP TABLE employees;
DROP TABLE departments;