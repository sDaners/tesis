-- ─────────────────────────────
--  SEQUENCES (one per table)
-- ─────────────────────────────
CREATE SEQUENCE departments_seq OPTIONS (sequence_kind = 'positive');
CREATE SEQUENCE employees_seq   OPTIONS (sequence_kind = 'positive');
CREATE SEQUENCE projects_seq    OPTIONS (sequence_kind = 'positive');
-- ─────────────────────────────
--  TABLE: departments
-- ─────────────────────────────
CREATE TABLE departments (
  dept_id     INT64 NOT NULL DEFAULT (NEXTVAL(SEQUENCE departments_seq)),
  dept_name   STRING(50) NOT NULL,
  location    STRING(100),
  created_at  TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (dept_id);

-- ─────────────────────────────
--  TABLE: employees
-- ─────────────────────────────
CREATE TABLE employees (
  emp_id       INT64 NOT NULL DEFAULT (NEXTVAL(SEQUENCE employees_seq)),
  first_name   STRING(50)  NOT NULL,
  last_name    STRING(50)  NOT NULL,
  email        STRING(150),
  hire_date    DATE        NOT NULL,
  salary       NUMERIC,
  dept_id      INT64,
  manager_id   INT64,
  phone_number STRING(20),

  CONSTRAINT fk_emp_dept    FOREIGN KEY (dept_id)    REFERENCES departments (dept_id),
  CONSTRAINT fk_emp_manager FOREIGN KEY (manager_id) REFERENCES employees   (emp_id)
) PRIMARY KEY (emp_id);

-- ─────────────────────────────
--  TABLE: projects
-- ─────────────────────────────
CREATE TABLE projects (
  project_id   INT64 NOT NULL DEFAULT (NEXTVAL(SEQUENCE projects_seq)),
  project_name STRING(100) NOT NULL,
  start_date   DATE,
  end_date     DATE,
  budget       NUMERIC,
  status       STRING(20) NOT NULL DEFAULT ('ACTIVE'),

  CONSTRAINT chk_dates   CHECK (end_date > start_date),
  CONSTRAINT chk_status  CHECK (status IN ('ACTIVE','COMPLETED','ON_HOLD','CANCELLED'))
) PRIMARY KEY (project_id);

-- ─────────────────────────────
--  TABLE: project_assignments
-- ─────────────────────────────
CREATE TABLE project_assignments (
  emp_id          INT64 NOT NULL,
  project_id      INT64 NOT NULL,
  role            STRING(50),
  hours_allocated INT64,

  CONSTRAINT fk_pa_emp     FOREIGN KEY (emp_id)     REFERENCES employees (emp_id),
  CONSTRAINT fk_pa_project FOREIGN KEY (project_id) REFERENCES projects  (project_id)
) PRIMARY KEY (emp_id, project_id);
-- unique e-mail address per employee
CREATE UNIQUE INDEX idx_emp_email     ON employees(email);

-- alphabetical employee search
CREATE INDEX        idx_emp_name      ON employees(last_name, first_name);

-- location look-up for departments
CREATE INDEX        idx_dept_location ON departments(location);

-- quick filtering by project status
CREATE INDEX        idx_project_status ON projects(status);

CREATE OR REPLACE VIEW employee_details SQL SECURITY INVOKER AS
SELECT
  e.emp_id,
  e.first_name,
  e.last_name,
  e.email,
  d.dept_name,
  m.first_name AS manager_first_name,
  m.last_name  AS manager_last_name
FROM employees   AS e
LEFT JOIN departments AS d ON e.dept_id    = d.dept_id
LEFT JOIN employees   AS m ON e.manager_id = m.emp_id;


--  Department insert
INSERT INTO departments (dept_name, location)
VALUES (@dept_name, @location)
THEN RETURN dept_id;

--  Employee insert
INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES (@first_name, @last_name, @email, @hire_date, @salary, @dept_id)
THEN RETURN emp_id;

--  Project insert
INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES (@project_name, @start_date, @end_date, @budget, @status)
THEN RETURN project_id;

--  Project assignment insert
INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES (@emp_id, @project_id, @role, @hours_allocated);

SELECT
  e.emp_id,
  e.first_name,
  e.last_name,
  e.email,
  d.dept_name,
  m.first_name AS manager_first_name,
  m.last_name  AS manager_last_name,
  p.project_name
FROM employees            AS e
LEFT JOIN departments        d ON e.dept_id     = d.dept_id
LEFT JOIN employees          m ON e.manager_id  = m.emp_id
LEFT JOIN project_assignments pa ON e.emp_id     = pa.emp_id
LEFT JOIN projects            p ON pa.project_id = p.project_id;

-- clean-up (order matters because of FK dependencies)
DROP VIEW    employee_details;

DROP INDEX   idx_project_status;
DROP INDEX   idx_dept_location;
DROP INDEX   idx_emp_name;
DROP INDEX   idx_emp_email;

DROP TABLE   project_assignments;
DROP TABLE   projects;
DROP TABLE   employees;
DROP TABLE   departments;

DROP SEQUENCE projects_seq;
DROP SEQUENCE employees_seq;
DROP SEQUENCE departments_seq; 