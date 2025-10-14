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
  hire_date DATE,
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
  CONSTRAINT check_dates CHECK (start_date IS NULL OR end_date IS NULL OR end_date > start_date),
  CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id);

CREATE TABLE project_assignments (
  emp_id STRING(36) NOT NULL,
  project_id STRING(36) NOT NULL,
  role STRING(50),
  hours_allocated INT64,
  CONSTRAINT fk_pa_emp FOREIGN KEY (emp_id) REFERENCES employees(emp_id),
  CONSTRAINT fk_pa_proj FOREIGN KEY (project_id) REFERENCES projects(project_id)
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

-- Insert department (parameters: @p1 = dept_name, @p2 = location)
INSERT INTO departments (dept_name, location)
VALUES (@p1, @p2);

-- Insert employee (parameters: @p1=first_name, @p2=last_name, @p3=email, @p4=hire_date_str, @p5=salary_str, @p6=dept_id, @p7=manager_id, @p8=phone_number)
INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id, manager_id, phone_number)
SELECT
  @p1,
  @p2,
  @p3,
  CASE
    WHEN REGEXP_CONTAINS(@p4, r'^\d{4}-\d{2}-\d{2}$') THEN PARSE_DATE('%Y-%m-%d', @p4)
    ELSE NULL
  END,
  CASE
    WHEN REGEXP_CONTAINS(@p5, r'^[+-]?\d+(\.\d+)?$') THEN CAST(@p5 AS NUMERIC)
    ELSE NULL
  END,
  d.dept_id,
  @p7,
  @p8
FROM UNNEST([1]) AS dummy
LEFT JOIN departments AS d ON d.dept_id = @p6;

-- Insert project (parameters: @p1=project_name, @p2=start_date_str, @p3=end_date_str, @p4=budget_str, @p5=status_str)
INSERT INTO projects (project_name, start_date, end_date, budget, status)
SELECT
  @p1,
  CASE WHEN REGEXP_CONTAINS(@p2, r'^\d{4}-\d{2}-\d{2}$') THEN PARSE_DATE('%Y-%m-%d', @p2) ELSE NULL END,
  CASE WHEN REGEXP_CONTAINS(@p3, r'^\d{4}-\d{2}-\d{2}$') THEN PARSE_DATE('%Y-%m-%d', @p3) ELSE NULL END,
  CASE WHEN REGEXP_CONTAINS(@p4, r'^[+-]?\d+(\.\d+)?$') THEN CAST(@p4 AS NUMERIC) ELSE NULL END,
  CASE WHEN @p5 IN ('ACTIVE','COMPLETED','ON_HOLD','CANCELLED') THEN @p5 ELSE 'ACTIVE' END
FROM UNNEST([1]) AS dummy;

-- Insert project assignment only if both employee and project exist (parameters: @p1=emp_id, @p2=project_id, @p3=role, @p4=hours_allocated_str)
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
SELECT
  e.emp_id,
  p.project_id,
  @p3,
  CASE WHEN REGEXP_CONTAINS(@p4, r'^\d+$') THEN CAST(@p4 AS INT64) ELSE NULL END
FROM employees AS e
JOIN projects AS p ON p.project_id = @p2
WHERE e.emp_id = @p1;

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