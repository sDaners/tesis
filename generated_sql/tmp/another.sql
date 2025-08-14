-- 64.4% parse, 42.1 success.

CREATE TABLE departments (
  dept_id INT64 NOT NULL,
  dept_name STRING(50) NOT NULL,
  location STRING(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
) PRIMARY KEY (dept_id);

CREATE SEQUENCE seq_departments START WITH 1;
ALTER TABLE departments ALTER COLUMN dept_id SET DEFAULT NEXTVAL(seq_departments);

CREATE TABLE employees (
  emp_id INT64 NOT NULL,
  first_name STRING(50) NOT NULL,
  last_name STRING(50) NOT NULL,
  email STRING(150),
  hire_date DATE NOT NULL,
  salary NUMERIC,
  dept_id INT64,
  manager_id INT64,
  phone_number STRING(20),
) PRIMARY KEY (emp_id);

CREATE SEQUENCE seq_employees START WITH 1;
ALTER TABLE employees ALTER COLUMN emp_id SET DEFAULT NEXTVAL(seq_employees);

ALTER TABLE employees ADD CONSTRAINT fk_dept FOREIGN KEY (dept_id) REFERENCES departments (dept_id);
ALTER TABLE employees ADD CONSTRAINT fk_manager FOREIGN KEY (manager_id) REFERENCES employees (emp_id);

CREATE TABLE projects (
  project_id INT64 NOT NULL,
  project_name STRING(100) NOT NULL,
  start_date DATE,
  end_date DATE,
  budget NUMERIC,
  status STRING(20) DEFAULT 'ACTIVE',
) PRIMARY KEY (project_id);

CREATE SEQUENCE seq_projects START WITH 1;
ALTER TABLE projects ALTER COLUMN project_id SET DEFAULT NEXTVAL(seq_projects);

ALTER TABLE projects ADD CONSTRAINT check_dates CHECK (end_date > start_date);
ALTER TABLE projects ADD CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'));

CREATE TABLE project_assignments (
  emp_id INT64 NOT NULL,
  project_id INT64 NOT NULL,
  role STRING(50),
  hours_allocated INT64,
) PRIMARY KEY (emp_id, project_id);

ALTER TABLE project_assignments ADD CONSTRAINT fk_emp FOREIGN KEY (emp_id) REFERENCES employees (emp_id);
ALTER TABLE project_assignments ADD CONSTRAINT fk_proj FOREIGN KEY (project_id) REFERENCES projects (project_id);

CREATE UNIQUE INDEX idx_emp_email ON employees(email);
CREATE INDEX idx_emp_name ON employees(last_name, first_name);
CREATE INDEX idx_dept_location ON departments(location);
CREATE INDEX idx_project_status ON projects(status);

CREATE OR REPLACE VIEW employee_details AS
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
VALUES (@dept_name, @location)
RETURNING dept_id;

INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES (@first_name, @last_name, @email, @hire_date, @salary, @dept_id)
RETURNING emp_id;

INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES (@project_name, @start_date, @end_date, @budget, @status)
RETURNING project_id;

INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES (@emp_id, @project_id, @role, @hours_allocated);

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

DROP VIEW employee_details;
DROP INDEX idx_project_status;
DROP INDEX idx_dept_location;
DROP INDEX idx_emp_name;
DROP INDEX idx_emp_email;

DROP TABLE project_assignments;
DROP TABLE projects;
DROP TABLE employees;
DROP TABLE departments;

DROP SEQUENCE seq_projects;
DROP SEQUENCE seq_employees;
DROP SEQUENCE seq_departments;
