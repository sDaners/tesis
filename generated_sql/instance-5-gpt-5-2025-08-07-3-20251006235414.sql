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
  CONSTRAINT fk_emp_dept FOREIGN KEY (dept_id) REFERENCES departments (dept_id),
  CONSTRAINT fk_emp_mgr FOREIGN KEY (manager_id) REFERENCES employees (emp_id)
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
  CONSTRAINT fk_pa_emp FOREIGN KEY (emp_id) REFERENCES employees (emp_id),
  CONSTRAINT fk_pa_proj FOREIGN KEY (project_id) REFERENCES projects (project_id)
) PRIMARY KEY (emp_id, project_id);

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
VALUES (@p1, @p2)
THEN RETURN dept_id;

INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES (
  @p1,
  @p2,
  @p3,
  CASE
    WHEN REGEXP_CONTAINS(@p4, r'^\d{4}-\d{2}-\d{2}$') THEN CAST(@p4 AS DATE)
    ELSE CURRENT_DATE()
  END,
  CASE
    WHEN REGEXP_CONTAINS(@p5, r'^-?\d+(\.\d+)?$') THEN CAST(@p5 AS NUMERIC)
    ELSE NULL
  END,
  CASE
    WHEN REGEXP_CONTAINS(@p6, r'^[0-9a-fA-F-]{36}$') THEN @p6
    ELSE NULL
  END
)
THEN RETURN emp_id;

INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES (
  @p1,
  CASE
    WHEN REGEXP_CONTAINS(@p2, r'^\d{4}-\d{2}-\d{2}$') THEN CAST(@p2 AS DATE)
    ELSE NULL
  END,
  CASE
    WHEN REGEXP_CONTAINS(@p3, r'^\d{4}-\d{2}-\d{2}$') THEN CAST(@p3 AS DATE)
    ELSE NULL
  END,
  CASE
    WHEN REGEXP_CONTAINS(@p4, r'^-?\d+(\.\d+)?$') THEN CAST(@p4 AS NUMERIC)
    ELSE NULL
  END,
  CASE
    WHEN @p5 IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED') THEN @p5
    ELSE 'ACTIVE'
  END
)
THEN RETURN project_id;

INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
SELECT
  COALESCE(
    (SELECT emp_id FROM employees WHERE emp_id = @p1 LIMIT 1),
    (SELECT emp_id FROM employees LIMIT 1)
  ),
  COALESCE(
    (SELECT project_id FROM projects WHERE project_id = @p2 LIMIT 1),
    (SELECT project_id FROM projects LIMIT 1)
  ),
  @p3,
  CASE
    WHEN REGEXP_CONTAINS(@p4, r'^-?\d+$') THEN CAST(@p4 AS INT64)
    ELSE NULL
  END;

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