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
  hire_date STRING(50),
  salary STRING(100),
  dept_id STRING(36),
  manager_id STRING(36),
  phone_number STRING(20),
  CONSTRAINT fk_employees_dept FOREIGN KEY (dept_id) REFERENCES departments(dept_id),
  CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
) PRIMARY KEY (emp_id);

CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name STRING(100) NOT NULL,
  start_date STRING(50),
  end_date STRING(50),
  budget STRING(100),
  status STRING(20) DEFAULT ('ACTIVE')
) PRIMARY KEY (project_id);

CREATE TABLE project_assignments (
  assignment_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  emp_id STRING(36),
  project_id STRING(36),
  role STRING(50),
  hours_allocated STRING(50),
  CONSTRAINT fk_pa_emp FOREIGN KEY (emp_id) REFERENCES employees(emp_id),
  CONSTRAINT fk_pa_project FOREIGN KEY (project_id) REFERENCES projects(project_id)
) PRIMARY KEY (assignment_id);

CREATE UNIQUE INDEX idx_emp_email ON employees(email);
CREATE INDEX idx_emp_name ON employees(last_name, first_name);
CREATE INDEX idx_dept_location ON departments(location);
CREATE INDEX idx_project_status ON projects(status);
CREATE UNIQUE INDEX uq_project_assignment_emp_project ON project_assignments(emp_id, project_id);

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

-- Insert a department (parameters: dept_name, location)
INSERT INTO departments (dept_name, location)
VALUES (@p1, @p2);

-- Insert an employee.
-- Parameters expected:
--   @p1 first_name
--   @p2 last_name
--   @p3 email
--   @p4 hire_date (kept as STRING to accept varied inputs)
--   @p5 salary (kept as STRING to accept varied inputs)
--   @p6 dept_name (we resolve dept_id by dept_name so FK is satisfied when dept exists; otherwise dept_id will be NULL)
INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
SELECT
  @p1,
  @p2,
  @p3,
  @p4,
  @p5,
  (SELECT dept_id FROM departments WHERE dept_name = @p6 LIMIT 1);

-- Insert a project (parameters: project_name, start_date, end_date, budget, status)
INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES (@p1, @p2, @p3, @p4, @p5);

-- Insert a project assignment.
-- Parameters expected:
--   @p1 employee email (will resolve to emp_id if exists)
--   @p2 project_name (will resolve to project_id if exists)
--   @p3 role
--   @p4 hours_allocated (kept as STRING to accept varied inputs)
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
SELECT
  (SELECT emp_id FROM employees WHERE email = @p1 LIMIT 1),
  (SELECT project_id FROM projects WHERE project_name = @p2 LIMIT 1),
  @p3,
  @p4;

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

DROP VIEW employee_details;
DROP INDEX idx_project_status;
DROP INDEX idx_dept_location;
DROP INDEX idx_emp_name;
DROP INDEX idx_emp_email;
DROP INDEX uq_project_assignment_emp_project;
DROP TABLE project_assignments;
DROP TABLE projects;
DROP TABLE employees;
DROP TABLE departments;