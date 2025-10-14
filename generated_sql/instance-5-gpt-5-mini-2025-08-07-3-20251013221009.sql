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
  CONSTRAINT check_dates CHECK (start_date IS NULL OR end_date IS NULL OR end_date > start_date),
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

-- Insert a department (params: @p1 = dept_name, @p2 = location)
INSERT INTO departments (dept_name, location)
VALUES (@p1, @p2);

-- Insert an employee (params: @p1=first_name, @p2=last_name, @p3=email, @p4=hire_date, @p5=salary, @p6=dept_id)
-- Ensure referenced department exists if a dept_id was supplied
INSERT INTO departments (dept_id, dept_name)
SELECT @p6, 'Imported Department'
WHERE @p6 IS NOT NULL AND NOT EXISTS (SELECT 1 FROM departments WHERE dept_id = @p6);

INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES (
  @p1,
  @p2,
  @p3,
  CASE
    WHEN @p4 IS NOT NULL AND REGEXP_CONTAINS(@p4, r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$') THEN CAST(@p4 AS DATE)
    ELSE CURRENT_DATE()
  END,
  CASE
    WHEN @p5 IS NOT NULL AND REGEXP_CONTAINS(@p5, r'^-?[0-9]+(\.[0-9]+)?$') THEN CAST(@p5 AS NUMERIC)
    ELSE NULL
  END,
  @p6
);

-- Insert a project (params: @p1=project_name, @p2=start_date, @p3=end_date, @p4=budget, @p5=status)
INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES (
  @p1,
  CASE
    WHEN @p2 IS NOT NULL AND REGEXP_CONTAINS(@p2, r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$') THEN CAST(@p2 AS DATE)
    ELSE NULL
  END,
  CASE
    WHEN @p3 IS NOT NULL AND REGEXP_CONTAINS(@p3, r'^[0-9]{4}-[0-9]{2}-[0-9]{2}$') THEN CAST(@p3 AS DATE)
    ELSE NULL
  END,
  CASE
    WHEN @p4 IS NOT NULL AND REGEXP_CONTAINS(@p4, r'^-?[0-9]+(\.[0-9]+)?$') THEN CAST(@p4 AS NUMERIC)
    ELSE NULL
  END,
  CASE
    WHEN @p5 IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED') THEN @p5
    ELSE 'ACTIVE'
  END
);

-- Insert a project assignment (params: @p1=emp_id, @p2=project_id, @p3=role, @p4=hours_allocated)
-- Ensure referenced employee and project exist; if they don't, create placeholder rows so FK constraints are satisfied.
INSERT INTO employees (emp_id, first_name, last_name, hire_date)
SELECT @p1, 'Imported', 'Employee', CURRENT_DATE()
WHERE @p1 IS NOT NULL AND NOT EXISTS (SELECT 1 FROM employees WHERE emp_id = @p1);

INSERT INTO projects (project_id, project_name, status)
SELECT @p2, 'Imported Project', 'ACTIVE'
WHERE @p2 IS NOT NULL AND NOT EXISTS (SELECT 1 FROM projects WHERE project_id = @p2);

INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES (
  @p1,
  @p2,
  @p3,
  CASE
    WHEN @p4 IS NOT NULL AND REGEXP_CONTAINS(@p4, r'^[0-9]+$') THEN CAST(@p4 AS INT64)
    ELSE 0
  END
);

-- Example select joining data
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