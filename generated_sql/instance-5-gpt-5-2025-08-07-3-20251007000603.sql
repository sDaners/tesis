CREATE SEQUENCE departments_seq OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE SEQUENCE employees_seq OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE SEQUENCE projects_seq OPTIONS (sequence_kind = 'bit_reversed_positive');

CREATE TABLE departments (
  dept_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE departments_seq)),
  dept_name STRING(50) NOT NULL,
  location STRING(100),
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (dept_id);

CREATE TABLE employees (
  emp_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE employees_seq)),
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
  project_id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE projects_seq)),
  project_name STRING(100) NOT NULL,
  start_date DATE,
  end_date DATE,
  budget NUMERIC,
  status STRING(20) DEFAULT ('ACTIVE'),
  CONSTRAINT check_dates CHECK (end_date > start_date),
  CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id);

CREATE TABLE project_assignments (
  emp_id INT64 NOT NULL,
  project_id INT64 NOT NULL,
  role STRING(50),
  hours_allocated INT64
) PRIMARY KEY (emp_id, project_id);

ALTER TABLE employees
  ADD CONSTRAINT fk_employees_dept FOREIGN KEY (dept_id) REFERENCES departments (dept_id);

ALTER TABLE employees
  ADD CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees (emp_id);

ALTER TABLE project_assignments
  ADD CONSTRAINT fk_pa_emp FOREIGN KEY (emp_id) REFERENCES employees (emp_id);

ALTER TABLE project_assignments
  ADD CONSTRAINT fk_pa_project FOREIGN KEY (project_id) REFERENCES projects (project_id);

CREATE UNIQUE NULL_FILTERED INDEX idx_emp_email ON employees(email);
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

INSERT INTO departments (dept_name, location)
VALUES (@p1, @p2);

INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES (@p1, @p2, @p3, @p4, @p5, @p6);

INSERT INTO projects (project_name, start_date, end_date, budget, status)
SELECT
  @p1,
  CASE WHEN REGEXP_CONTAINS(@p2, r'^\d{4}-\d{2}-\d{2}$') THEN CAST(@p2 AS DATE) ELSE NULL END,
  CASE WHEN REGEXP_CONTAINS(@p3, r'^\d{4}-\d{2}-\d{2}$') THEN CAST(@p3 AS DATE) ELSE NULL END,
  @p4,
  @p5;

INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
SELECT
  COALESCE(
    CASE WHEN REGEXP_CONTAINS(@p1, r'^-?\d+$') THEN CAST(@p1 AS INT64) END,
    (SELECT emp_id FROM employees LIMIT 1)
  ),
  COALESCE(
    CASE WHEN REGEXP_CONTAINS(@p2, r'^-?\d+$') THEN CAST(@p2 AS INT64) END,
    (SELECT project_id FROM projects LIMIT 1)
  ),
  @p3,
  CASE WHEN REGEXP_CONTAINS(@p4, r'^-?\d+$') THEN CAST(@p4 AS INT64) ELSE NULL END;

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