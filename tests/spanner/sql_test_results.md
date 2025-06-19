# SQL Test Results

Generated: 2025-06-19 14:49:31

## Summary

- **Total SQL Files**: 6
- **Total Statements**: 123
- **Successfully Executed**: 61
- **Failed**: 62
- **Overall Success Rate**: 49.6%

## Error Code Summary

| Error Code | Total Occurrences | Description |
|------------|-------------------|-------------|
| InvalidArgument | 43 | Invalid SQL syntax or unsupported features |
| NotFound | 19 | Referenced table, column, or object not found |

## Detailed Results

| File | Total | CREATE | INSERT | SELECT | DROP | Executed | Failed | Success Rate | Execution Time |
|------|-------|--------|--------|--------|------|----------|--------|--------------|----------------|
| [gpt-4o-2.sql](../../generated_sql/gpt-4o-2.sql) | 11 | 9 | 1 | 1 | 0 | 6 | 5 | 54.5% | 5.491s |
| [gpt-4o-p2.sql](../../generated_sql/gpt-4o-p2.sql) | 27 | 9 | 8 | 1 | 9 | 9 | 18 | 33.3% | 3.444s |
| [gpt-4o.sql](../../generated_sql/gpt-4o.sql) | 10 | 8 | 1 | 1 | 0 | 2 | 8 | 20.0% | 3.389s |
| [gpt-o3.sql](../../generated_sql/gpt-o3.sql) | 29 | 12 | 4 | 1 | 12 | 12 | 17 | 41.4% | 3.428s |
| [gpt-o4-mini-high.sql](../../generated_sql/gpt-o4-mini-high.sql) | 23 | 9 | 4 | 1 | 9 | 9 | 14 | 39.1% | 3.472s |
| [valid_spanner_database.sql](../../generated_sql/valid_spanner_database.sql) | 23 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 3.472s |

## Error Details

### gpt-4o-2.sql

**Error Rate**: 45.5% (5/11 failed)

**Error Codes**:
- `InvalidArgument`: 4 occurrences
- `NotFound`: 1 occurrences

**Sample Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE TABLE departments ( dept_id INT64 NOT NULL, dept_name STRING(50) NOT NULL, location STRING(100), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(), ) PRIMARY KEY (dept_id) : Syntax error on line 1, column 134: Expecting '(' but found 'CURRENT_TIMESTAMP'
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
3. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = View `employee_details` is missing the SQL SECURITY clause.
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_id, dept_name, location) VALUES (NULL, 'Enginee...\n            ^", requestID = "1.4c9dbd79443e577b.1.4.5.1"
5. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:179]\n...p.project_name FROM employees e LEFT JOIN departments d ON e.dept_id = d.d...\n                                             ^", requestID = "1.4c9dbd79443e577b.1.4.6.1"

### gpt-4o-p2.sql

**Error Rate**: 66.7% (18/27 failed)

**Error Codes**:
- `InvalidArgument`: 13 occurrences
- `NotFound`: 5 occurrences

**Sample Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `departments`.`dept_id`: Expected type INT64; found STRING [at 1:1]
GENERATE_UUID()
^
2. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `employees`.`emp_id`: Expected type INT64; found STRING [at 1:1]
GENERATE_UUID()
^
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
4. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE TABLE projects ( project_id INT64 NOT NULL DEFAULT (GENERATE_UUID()), project_name STRING(100) NOT NULL, start_date DATE, end_date DATE, budget NUMERIC, status STRING(20) DEFAULT 'ACTIVE', CONSTRAINT check_dates CHECK (end_date > start_date), CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED')) ) PRIMARY KEY (project_id) : Syntax error on line 1, column 187: Expecting '(' but found ''ACTIVE''
5. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
6. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
7. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
8. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
9. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = View `employee_details` is missing the SQL SECURITY clause.
10. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_id, dept_name, location) VALUES (NULL, 'Enginee...\n            ^", requestID = "1.4c9dbd79443e577b.2.2.5.1"
11. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.4c9dbd79443e577b.2.2.6.1"
12. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_id, dept_name, location) VALUES (NULL, 'Enginee...\n            ^", requestID = "1.4c9dbd79443e577b.2.2.7.1"
13. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.4c9dbd79443e577b.2.2.8.1"
14. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salar...\n            ^", requestID = "1.4c9dbd79443e577b.2.2.9.1"
15. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salar...\n            ^", requestID = "1.4c9dbd79443e577b.2.2.10.1"
16. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.4c9dbd79443e577b.2.2.11.1"
17. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.4c9dbd79443e577b.2.2.12.1"
18. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.4c9dbd79443e577b.2.2.13.1"

### gpt-4o.sql

**Error Rate**: 80.0% (8/10 failed)

**Error Codes**:
- `InvalidArgument`: 4 occurrences
- `NotFound`: 4 occurrences

**Sample Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE TABLE departments ( dept_id INT64 NOT NULL, dept_name STRING(50) NOT NULL, location STRING(100), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ) PRIMARY KEY (dept_id) : Syntax error on line 1, column 134: Expecting '(' but found 'CURRENT_TIMESTAMP'
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
5. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
6. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = View `employee_details` is missing the SQL SECURITY clause.
7. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_id, dept_name, location) VALUES (NULL, 'Enginee...\n            ^", requestID = "1.4c9dbd79443e577b.3.2.5.1"
8. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.4c9dbd79443e577b.3.2.6.1"

### gpt-o3.sql

**Error Rate**: 58.6% (17/29 failed)

**Error Codes**:
- `InvalidArgument`: 12 occurrences
- `NotFound`: 5 occurrences

**Sample Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE SEQUENCE departments_seq OPTIONS (sequence_kind = 'positive') : Unsupported sequence kind: positive
2. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE SEQUENCE employees_seq OPTIONS (sequence_kind = 'positive') : Unsupported sequence kind: positive
3. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE SEQUENCE projects_seq OPTIONS (sequence_kind = 'positive') : Unsupported sequence kind: positive
4. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE TABLE departments ( dept_id INT64 NOT NULL DEFAULT (NEXTVAL(SEQUENCE departments_seq)), dept_name STRING(50) NOT NULL, location STRING(100), created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP() ) PRIMARY KEY (dept_id) : Syntax error on line 1, column 187: Expecting '(' but found 'CURRENT_TIMESTAMP'
5. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `employees`.`emp_id`: Function not found: NEXTVAL [at 1:1]
NEXTVAL(SEQUENCE employees_seq)
^
6. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE TABLE projects ( project_id INT64 NOT NULL DEFAULT (NEXTVAL(SEQUENCE projects_seq)), project_name STRING(100) NOT NULL, start_date DATE, end_date DATE, budget NUMERIC, status STRING(20) NOT NULL DEFAULT 'ACTIVE', CONSTRAINT chk_dates CHECK (end_date > start_date), CONSTRAINT chk_status CHECK (status IN ('ACTIVE','COMPLETED','ON_HOLD','CANCELLED')) ) PRIMARY KEY (project_id) : Syntax error on line 1, column 211: Expecting '(' but found ''ACTIVE''
7. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
8. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
9. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
10. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
11. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
12. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = View `employee_details` is missing the SQL SECURITY clause.
13. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_name, location) VALUES ('Engineering', 'New Yor...\n            ^", requestID = "1.4c9dbd79443e577b.4.4.5.1"
14. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_name, start_date, end_date, budget, status) VAL...\n            ^", requestID = "1.4c9dbd79443e577b.4.4.6.1"
15. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_...\n            ^", requestID = "1.4c9dbd79443e577b.4.4.7.1"
16. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.4c9dbd79443e577b.4.4.8.1"
17. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees AS e LEFT JO...\n                                                       ^", requestID = "1.4c9dbd79443e577b.4.4.9.1"

### gpt-o4-mini-high.sql

**Error Rate**: 60.9% (14/23 failed)

**Error Codes**:
- `InvalidArgument`: 10 occurrences
- `NotFound`: 4 occurrences

**Sample Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE TABLE IF NOT EXISTS Departments ( dept_id INT64 NOT NULL GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY, dept_name STRING(50) NOT NULL, location STRING(100), created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP() ) : Syntax error on line 1, column 202: Expecting '(' but found 'CURRENT_TIMESTAMP'
2. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = The sequence kind of an identity column emp_id is not specified. Please specify the sequence kind explicitly or set the database option `default_sequence_kind`.
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: Employees
4. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE TABLE IF NOT EXISTS Projects ( project_id INT64 NOT NULL GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY, project_name STRING(100) NOT NULL, start_date DATE, end_date DATE, budget NUMERIC, status STRING(20) NOT NULL DEFAULT 'ACTIVE', CONSTRAINT check_dates CHECK (end_date > start_date), CONSTRAINT check_status CHECK (status IN ('ACTIVE','COMPLETED','ON_HOLD','CANCELLED')) ) : Syntax error on line 1, column 229: Expecting '(' but found ''ACTIVE''
5. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE TABLE IF NOT EXISTS ProjectAssignments ( emp_id INT64 NOT NULL, project_id INT64 NOT NULL, role STRING(50), hours_allocated INT64, PRIMARY KEY (emp_id, project_id), FOREIGN KEY (emp_id) REFERENCES Employees(emp_id), FOREIGN KEY (project_id) REFERENCES Projects(project_id) ) : Syntax error on line 1, column 151: Expecting ')' but found '('
6. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: Employees
7. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: Departments
8. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: Projects
9. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the definition of view `employee_details`: Table not found: Employees [at 1:255]
...m.last_name AS manager_last_name FROM Employees AS e LEFT JOIN Departments...
                                         ^
10. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: Departments [at 1:13]\nINSERT INTO Departments (dept_name, location) VALUES ('Sample Name', 'New Yor...\n            ^", requestID = "1.4c9dbd79443e577b.5.4.5.1"
11. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: Projects [at 1:13]\nINSERT INTO Projects (project_name, start_date, end_date, budget, status) VAL...\n            ^", requestID = "1.4c9dbd79443e577b.5.4.6.1"
12. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: Employees [at 1:13]\nINSERT INTO Employees (first_name, last_name, email, hire_date, salary, dept_...\n            ^", requestID = "1.4c9dbd79443e577b.5.4.7.1"
13. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: ProjectAssignments [at 1:13]\nINSERT INTO ProjectAssignments (emp_id, project_id, role, hours_allocated) VA...\n            ^", requestID = "1.4c9dbd79443e577b.5.4.8.1"
14. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: Employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM Employees AS e LEFT JO...\n                                                       ^", requestID = "1.4c9dbd79443e577b.5.4.9.1"

## Compatibility Insights

### Common Issues Found

- **DEFAULT value syntax**: Found in 2 files
- **FOREIGN KEY constraints**: Found in 1 files
- **CURRENT_TIMESTAMP syntax**: Found in 4 files
- **Table dependency issues**: Found in 5 files
- **Missing SQL SECURITY clause in views**: Found in 4 files
- **GENERATE_UUID() compatibility**: Found in 1 files
