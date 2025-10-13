--  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

-- Total statements: 23
-- Successfully parsed: 23
-- Parse errors: 0
-- Executed: 23
-- Execution errors: 0
-- Parse success rate: 100.0%
-- Execution success rate (of parsed): 100.0%
-- Overall success rate: 100.0%

CREATE TABLE departments (
    dept_id STRING(36) DEFAULT (GENERATE_UUID()),
    dept_name STRING(50) NOT NULL,
    location STRING(100),
    created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (dept_id);

CREATE TABLE employees (
    emp_id STRING(36) DEFAULT (GENERATE_UUID()),
    first_name STRING(50) NOT NULL,
    last_name STRING(50) NOT NULL,
    email STRING(150),
    hire_date TIMESTAMP NOT NULL,
    salary FLOAT64,
    dept_id STRING(36),
    manager_id STRING(36),
    phone_number STRING(20),
    CONSTRAINT fk_dept FOREIGN KEY (dept_id) REFERENCES departments(dept_id),
    CONSTRAINT fk_manager FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
) PRIMARY KEY (emp_id);

CREATE TABLE projects (
    project_id STRING(36) DEFAULT (GENERATE_UUID()),
    project_name STRING(100) NOT NULL,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    budget FLOAT64,
    status STRING(20) DEFAULT ('ACTIVE'),
    CONSTRAINT check_dates CHECK (end_date > start_date),
    CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id);

CREATE TABLE project_assignments (
    emp_id STRING(36) NOT NULL,
    project_id STRING(36) NOT NULL,
    role STRING(50),
    hours_allocated INT64,
    CONSTRAINT fk_emp FOREIGN KEY (emp_id) REFERENCES employees(emp_id),
    CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects(project_id)
) PRIMARY KEY (emp_id, project_id);

CREATE UNIQUE INDEX idx_emp_email ON employees(email);
CREATE INDEX idx_emp_name ON employees(last_name, first_name);
CREATE INDEX idx_dept_location ON departments(location);
CREATE INDEX idx_project_status ON projects(status);

CREATE OR REPLACE VIEW employee_details
    SQL SECURITY INVOKER
    AS SELECT 
        e.emp_id,
        e.first_name,
        e.last_name,
        e.email,
        d.dept_name,
        m.first_name as manager_first_name,
        m.last_name as manager_last_name
    FROM employees e
    LEFT JOIN departments d ON e.dept_id = d.dept_id
    LEFT JOIN employees m ON e.manager_id = m.emp_id;

INSERT INTO departments (dept_name, location)
VALUES (@dept_name, @location)
THEN RETURN dept_id;

INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES (@first_name, @last_name, @email, @hire_date, @salary, @dept_id)
THEN RETURN emp_id;

INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES (@project_name, @start_date, @end_date, @budget, @status)
THEN RETURN project_id;

INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES (@emp_id, @project_id, @role, @hours);

SELECT e.emp_id, e.first_name, e.last_name, e.email, d.dept_name, 
       m.first_name as manager_first_name, m.last_name as manager_last_name,
       p.project_name
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
LEFT JOIN project_assignments pa ON e.emp_id = pa.emp_id
LEFT JOIN projects p ON pa.project_id = p.project_id;

DROP VIEW employee_details;
DROP INDEX  idx_project_status;
DROP INDEX  idx_dept_location;
DROP INDEX  idx_emp_name;
DROP INDEX  idx_emp_email;
DROP TABLE  project_assignments;
DROP TABLE  projects;
DROP TABLE  employees;
DROP TABLE  departments; 