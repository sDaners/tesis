# SQL Test Results

Generated: 2025-10-12 23:51:16

## Summary

- **Total SQL Files**: 63
- **Total Statements**: 1502
- **Successfully Parsed**: 1402
- **Parse Errors**: 100
- **Successfully Executed**: 1148
- **Execution Errors**: 254
- **Parse Success Rate**: 93.3%
- **Execution Success Rate** (of parsed): 81.9%
- **Overall Success Rate**: 76.4%

## Parse Error Summary

| Parse Error Type | Total Occurrences | Description |
|------------------|-------------------|-------------|
| Syntax Error: Expected Token | 66 | Missing expected tokens in SQL syntax. FIX: Add required keywords, punctuation, or identifiers where expected. After DEFAULT remember to wrap the value in parentheses |
| Syntax Error: String Literal Quotes | 17 | String literals in SQL statements are using incorrect quote types. FIX: Use single quotes for string literals instead of double quotes. Change "ACTIVE" to 'ACTIVE'. Spanner requires single quotes for string constants. |
| Syntax Error: PRIMARY/FOREIGN KEY Placement | 11 | PRIMARY KEY constraints MUST be placed OUTSIDE the column definition parentheses in Spanner. FIX: Move PRIMARY KEY clause to after the closing parenthesis of column definitions. CORRECT: ') PRIMARY KEY (column_name);' WRONG: 'PRIMARY KEY (column_name)' inside column list. This is a critical Spanner-specific syntax requirement. |
| Syntax Error: General | 6 | General SQL syntax errors not matching specific patterns. FIX: Review statement structure, check for typos and syntax compliance with Spanner SQL |

## Execution Error Code Summary

| Error Code | Total Occurrences | Description |
|------------|-------------------|-------------|
| InvalidArgument | 132 | Invalid SQL syntax or unsupported features. FIX: Check for Spanner-specific syntax requirements (e.g., CURRENT_TIMESTAMP vs CURRENT_TIMESTAMP(), required clauses in views) |
| NotFound | 90 | Referenced table, column, or object not found. FIX: Ensure all tables/columns exist before referencing them, or create them first in dependency order |
| FailedPrecondition | 30 | Constraint violations or prerequisite not met. FIX: Check for NOT NULL constraints, foreign key violations, or missing required data |
| OutOfRange | 2 | Unknown error code |

## Error Category Analysis

| Error Category | Total Occurrences | Files Affected | Percentage |
|----------------|-------------------|----------------|------------|
| Table Not Found (InvalidArgument) | 95 | 28/46 | 37.4% |
| NotFound | 90 | 28/46 | 35.4% |
| FailedPrecondition | 30 | 18/46 | 11.8% |
| InvalidArgument: Other | 9 | 5/46 | 3.5% |
| Unsupported Feature: Sequence Kind | 6 | 2/46 | 2.4% |
| Syntax Error: General | 6 | 6/46 | 2.4% |
| Type Mismatch: GENERATE_UUID on INT64 | 5 | 2/46 | 2.0% |
| Function Not Found: NEXTVAL | 4 | 2/46 | 1.6% |
| Default Value: Parsing Error | 3 | 1/46 | 1.2% |
| Missing Clause: General | 3 | 1/46 | 1.2% |
| OutOfRange | 2 | 2/46 | 0.8% |
| Identity Column: Missing Sequence Kind | 1 | 1/46 | 0.4% |

## Detailed Results

| File | Total | Parsed | Parse Errors | CREATE | INSERT | SELECT | DROP | Executed | Exec Errors | Parse Success | Exec Success | Total Time |
|------|-------|--------|--------------|--------|--------|--------|------|----------|-------------|---------------|--------------|------------|
| [concurrent_5.sql](../../generated_sql/concurrent_5.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 100.0% | 4.274s |
| [fix-gpt-4o-2.sql](../../generated_sql/fix-gpt-4o-2.sql) | 17 | 17 | 0 | 9 | 1 | 1 | 0 | 16 | 1 | 100.0% | 94.1% | 3.639s |
| [fix-gpt-4o-p2.sql](../../generated_sql/fix-gpt-4o-p2.sql) | 27 | 27 | 0 | 9 | 8 | 1 | 9 | 9 | 18 | 100.0% | 33.3% | 3.599s |
| [fix-gpt-4o.sql](../../generated_sql/fix-gpt-4o.sql) | 10 | 10 | 0 | 8 | 1 | 1 | 0 | 4 | 6 | 100.0% | 40.0% | 3.634s |
| [fix-gpt-o3.sql](../../generated_sql/fix-gpt-o3.sql) | 29 | 29 | 0 | 12 | 4 | 1 | 12 | 12 | 17 | 100.0% | 41.4% | 3.565s |
| [fix-gpt-o4-mini-high.sql](../../generated_sql/fix-gpt-o4-mini-high.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 21 | 2 | 100.0% | 91.3% | 3.687s |
| [gpt-4o-2.sql](../../generated_sql/gpt-4o-2.sql) | 17 | 15 | 2 | 7 | 1 | 1 | 0 | 11 | 4 | 88.2% | 73.3% | 3.559s |
| [gpt-4o-p2.sql](../../generated_sql/gpt-4o-p2.sql) | 27 | 25 | 2 | 7 | 8 | 1 | 9 | 9 | 16 | 92.6% | 36.0% | 3.599s |
| [gpt-4o-recommendations.sql](../../generated_sql/gpt-4o-recommendations.sql) | 23 | 16 | 7 | 5 | 1 | 1 | 9 | 9 | 7 | 69.6% | 56.2% | 3.55s |
| [gpt-4o-recommendations2.sql](../../generated_sql/gpt-4o-recommendations2.sql) | 31 | 21 | 10 | 7 | 4 | 1 | 9 | 15 | 6 | 67.7% | 71.4% | 3.593s |
| [gpt-5-3.sql](../../generated_sql/gpt-5-3.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 100.0% | 3.698s |
| [gpt-5-3shot.sql](../../generated_sql/gpt-5-3shot.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 19 | 4 | 100.0% | 82.6% | 3.603s |
| [gpt-5.sql](../../generated_sql/gpt-5.sql) | 23 | 18 | 5 | 5 | 4 | 1 | 8 | 8 | 10 | 78.3% | 44.4% | 3.602s |
| [gpt-o3.sql](../../generated_sql/gpt-o3.sql) | 29 | 26 | 3 | 9 | 4 | 1 | 12 | 12 | 14 | 89.7% | 46.2% | 3.616s |
| [gpt-o4-mini-high.sql](../../generated_sql/gpt-o4-mini-high.sql) | 23 | 16 | 7 | 6 | 0 | 1 | 9 | 9 | 7 | 69.6% | 56.2% | 3.567s |
| [instance-1-chatgpt-4o-latest-2-20251005112832.sql](../../generated_sql/instance-1-chatgpt-4o-latest-2-20251005112832.sql) | 23 | 22 | 1 | 8 | 4 | 1 | 9 | 17 | 5 | 95.7% | 77.3% | 3.617s |
| [instance-1-chatgpt-4o-latest-3-20251007203321.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007203321.sql) | 23 | 22 | 1 | 8 | 4 | 1 | 9 | 15 | 7 | 95.7% | 68.2% | 3.573s |
| [instance-1-chatgpt-4o-latest-3-20251007203554.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007203554.sql) | 23 | 22 | 1 | 8 | 4 | 1 | 9 | 15 | 7 | 95.7% | 68.2% | 3.564s |
| [instance-1-chatgpt-4o-latest-3-20251007203716.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007203716.sql) | 27 | 27 | 0 | 9 | 4 | 1 | 9 | 26 | 1 | 100.0% | 96.3% | 3.632s |
| [instance-1-chatgpt-4o-latest-3-20251007210819.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007210819.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 22 | 1 | 100.0% | 95.7% | 3.651s |
| [instance-1-chatgpt-4o-latest-3-20251007213639.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007213639.sql) | 25 | 25 | 0 | 9 | 4 | 1 | 9 | 23 | 2 | 100.0% | 92.0% | 3.62s |
| [instance-1-chatgpt-4o-latest-3-20251007213715.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007213715.sql) | 27 | 27 | 0 | 9 | 4 | 1 | 9 | 24 | 3 | 100.0% | 88.9% | 3.699s |
| [instance-1-chatgpt-4o-latest-3-20251007213745.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007213745.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 22 | 1 | 100.0% | 95.7% | 3.647s |
| [instance-1-chatgpt-4o-latest-3-20251007213927.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007213927.sql) | 29 | 29 | 0 | 9 | 4 | 1 | 9 | 27 | 2 | 100.0% | 93.1% | 3.713s |
| [instance-1-chatgpt-4o-latest-3-20251007221438.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007221438.sql) | 19 | 18 | 1 | 8 | 0 | 1 | 9 | 15 | 3 | 94.7% | 83.3% | 3.589s |
| [instance-1-chatgpt-4o-latest-3-20251007232200.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007232200.sql) | 19 | 19 | 0 | 9 | 0 | 1 | 9 | 19 | 0 | 100.0% | 100.0% | 3.599s |
| [instance-1-chatgpt-4o-latest-3-20251007232334.sql](../../generated_sql/instance-1-chatgpt-4o-latest-3-20251007232334.sql) | 29 | 29 | 0 | 9 | 4 | 1 | 9 | 29 | 0 | 100.0% | 100.0% | 3.639s |
| [instance-1-gpt-4o-mini-3-20251009005208.sql](../../generated_sql/instance-1-gpt-4o-mini-3-20251009005208.sql) | 23 | 15 | 8 | 5 | 0 | 1 | 9 | 9 | 6 | 65.2% | 60.0% | 3.551s |
| [instance-1-gpt-5-2025-08-07-3-20251006235414.sql](../../generated_sql/instance-1-gpt-5-2025-08-07-3-20251006235414.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 22 | 1 | 100.0% | 95.7% | 3.701s |
| [instance-1-gpt-5-2025-08-07-3-20251007000603.sql](../../generated_sql/instance-1-gpt-5-2025-08-07-3-20251007000603.sql) | 24 | 24 | 0 | 9 | 4 | 1 | 9 | 24 | 0 | 100.0% | 100.0% | 3.659s |
| [instance-1-gpt-5-2025-08-07-3-20251007001238.sql](../../generated_sql/instance-1-gpt-5-2025-08-07-3-20251007001238.sql) | 26 | 26 | 0 | 12 | 4 | 1 | 9 | 26 | 0 | 100.0% | 100.0% | 3.655s |
| [instance-1-gpt-5-2025-08-07-3-20251007002717.sql](../../generated_sql/instance-1-gpt-5-2025-08-07-3-20251007002717.sql) | 29 | 29 | 0 | 12 | 4 | 1 | 12 | 29 | 0 | 100.0% | 100.0% | 3.611s |
| [instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql](../../generated_sql/instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql) | 23 | 15 | 8 | 4 | 1 | 1 | 9 | 9 | 6 | 65.2% | 60.0% | 3.583s |
| [instance-1-gpt-5-mini-2025-08-07-1-20251007233356.sql](../../generated_sql/instance-1-gpt-5-mini-2025-08-07-1-20251007233356.sql) | 26 | 21 | 5 | 4 | 4 | 4 | 9 | 9 | 12 | 80.8% | 42.9% | 3.571s |
| [instance-2-chatgpt-4o-latest-2-20251005112832.sql](../../generated_sql/instance-2-chatgpt-4o-latest-2-20251005112832.sql) | 23 | 22 | 1 | 8 | 4 | 1 | 9 | 17 | 5 | 95.7% | 77.3% | 3.658s |
| [instance-2-chatgpt-4o-latest-3-20251007203405.sql](../../generated_sql/instance-2-chatgpt-4o-latest-3-20251007203405.sql) | 23 | 22 | 1 | 8 | 4 | 1 | 9 | 15 | 7 | 95.7% | 68.2% | 3.623s |
| [instance-2-chatgpt-4o-latest-3-20251007203554.sql](../../generated_sql/instance-2-chatgpt-4o-latest-3-20251007203554.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 22 | 1 | 100.0% | 95.7% | 3.642s |
| [instance-2-gpt-4o-mini-3-20251009005208.sql](../../generated_sql/instance-2-gpt-4o-mini-3-20251009005208.sql) | 23 | 14 | 9 | 4 | 0 | 1 | 9 | 9 | 5 | 60.9% | 64.3% | 3.595s |
| [instance-2-gpt-5-2025-08-07-3-20251006235414.sql](../../generated_sql/instance-2-gpt-5-2025-08-07-3-20251006235414.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 100.0% | 3.632s |
| [instance-2-gpt-5-2025-08-07-3-20251007000603.sql](../../generated_sql/instance-2-gpt-5-2025-08-07-3-20251007000603.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 21 | 2 | 100.0% | 91.3% | 3.638s |
| [instance-2-gpt-5-2025-08-07-3-20251007001238.sql](../../generated_sql/instance-2-gpt-5-2025-08-07-3-20251007001238.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 100.0% | 3.673s |
| [instance-2-gpt-5-2025-08-07-3-20251007002717.sql](../../generated_sql/instance-2-gpt-5-2025-08-07-3-20251007002717.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 100.0% | 3.75s |
| [instance-3-chatgpt-4o-latest-2-20251005112832.sql](../../generated_sql/instance-3-chatgpt-4o-latest-2-20251005112832.sql) | 23 | 22 | 1 | 8 | 4 | 1 | 9 | 17 | 5 | 95.7% | 77.3% | 3.644s |
| [instance-3-gpt-4o-mini-3-20251009005208.sql](../../generated_sql/instance-3-gpt-4o-mini-3-20251009005208.sql) | 23 | 15 | 8 | 5 | 0 | 1 | 9 | 9 | 6 | 65.2% | 60.0% | 3.596s |
| [instance-3-gpt-5-2025-08-07-3-20251006235414.sql](../../generated_sql/instance-3-gpt-5-2025-08-07-3-20251006235414.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 21 | 2 | 100.0% | 91.3% | 3.611s |
| [instance-3-gpt-5-2025-08-07-3-20251007000603.sql](../../generated_sql/instance-3-gpt-5-2025-08-07-3-20251007000603.sql) | 26 | 26 | 0 | 12 | 4 | 1 | 9 | 26 | 0 | 100.0% | 100.0% | 3.661s |
| [instance-3-gpt-5-2025-08-07-3-20251007001238.sql](../../generated_sql/instance-3-gpt-5-2025-08-07-3-20251007001238.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 100.0% | 3.623s |
| [instance-3-gpt-5-2025-08-07-3-20251007002717.sql](../../generated_sql/instance-3-gpt-5-2025-08-07-3-20251007002717.sql) | 26 | 26 | 0 | 12 | 4 | 1 | 9 | 12 | 14 | 100.0% | 46.2% | 3.694s |
| [instance-4-chatgpt-4o-latest-2-20251005112832.sql](../../generated_sql/instance-4-chatgpt-4o-latest-2-20251005112832.sql) | 23 | 22 | 1 | 8 | 4 | 1 | 9 | 16 | 6 | 95.7% | 72.7% | 3.609s |
| [instance-4-gpt-4o-mini-3-20251009005208.sql](../../generated_sql/instance-4-gpt-4o-mini-3-20251009005208.sql) | 23 | 14 | 9 | 4 | 0 | 1 | 9 | 9 | 5 | 60.9% | 64.3% | 3.59s |
| [instance-4-gpt-5-2025-08-07-3-20251006235414.sql](../../generated_sql/instance-4-gpt-5-2025-08-07-3-20251006235414.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 100.0% | 3.705s |
| [instance-4-gpt-5-2025-08-07-3-20251007000603.sql](../../generated_sql/instance-4-gpt-5-2025-08-07-3-20251007000603.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 21 | 2 | 100.0% | 91.3% | 6.788s |
| [instance-4-gpt-5-2025-08-07-3-20251007001238.sql](../../generated_sql/instance-4-gpt-5-2025-08-07-3-20251007001238.sql) | 33 | 33 | 0 | 12 | 4 | 1 | 12 | 33 | 0 | 100.0% | 100.0% | 3.645s |
| [instance-4-gpt-5-2025-08-07-3-20251007002717.sql](../../generated_sql/instance-4-gpt-5-2025-08-07-3-20251007002717.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 20 | 3 | 100.0% | 87.0% | 3.616s |
| [instance-5-chatgpt-4o-latest-2-20251005112832.sql](../../generated_sql/instance-5-chatgpt-4o-latest-2-20251005112832.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 21 | 2 | 100.0% | 91.3% | 3.605s |
| [instance-5-chatgpt-4o-latest-3-20251007203405.sql](../../generated_sql/instance-5-chatgpt-4o-latest-3-20251007203405.sql) | 23 | 22 | 1 | 8 | 4 | 1 | 9 | 15 | 7 | 95.7% | 68.2% | 3.64s |
| [instance-5-gpt-4o-mini-3-20251009005208.sql](../../generated_sql/instance-5-gpt-4o-mini-3-20251009005208.sql) | 23 | 16 | 7 | 6 | 0 | 1 | 9 | 14 | 2 | 69.6% | 87.5% | 3.553s |
| [instance-5-gpt-5-2025-08-07-3-20251006235414.sql](../../generated_sql/instance-5-gpt-5-2025-08-07-3-20251006235414.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 100.0% | 3.683s |
| [instance-5-gpt-5-2025-08-07-3-20251007000603.sql](../../generated_sql/instance-5-gpt-5-2025-08-07-3-20251007000603.sql) | 30 | 30 | 0 | 12 | 4 | 1 | 9 | 27 | 3 | 100.0% | 90.0% | 3.691s |
| [instance-5-gpt-5-2025-08-07-3-20251007001238.sql](../../generated_sql/instance-5-gpt-5-2025-08-07-3-20251007001238.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 100.0% | 3.691s |
| [instance-5-gpt-5-2025-08-07-3-20251007002717.sql](../../generated_sql/instance-5-gpt-5-2025-08-07-3-20251007002717.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 20 | 3 | 100.0% | 87.0% | 3.656s |
| [instance-7-chatgpt-4o-latest-3-20251007203405.sql](../../generated_sql/instance-7-chatgpt-4o-latest-3-20251007203405.sql) | 23 | 22 | 1 | 8 | 4 | 1 | 9 | 17 | 5 | 95.7% | 77.3% | 3.616s |
| [valid_spanner_database.sql](../../generated_sql/valid_spanner_database.sql) | 23 | 23 | 0 | 9 | 4 | 1 | 9 | 23 | 0 | 100.0% | 100.0% | 3.605s |

## Error Details

### fix-gpt-4o-2.sql

**Execution Error Rate**: 5.9% (1/17 parsed statements failed)

**Execution Error Codes**:
- `FailedPrecondition`: 1 occurrences

**Error Categories**:
- `FailedPrecondition`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: departments.dept_id in table: departments", requestID = "1.c3438849fd14d720.2.4.5.1"

### fix-gpt-4o-p2.sql

**Execution Error Rate**: 66.7% (18/27 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 13 occurrences
- `NotFound`: 5 occurrences

**Error Categories**:
- `Type Mismatch: GENERATE_UUID on INT64`: 3 occurrences
- `NotFound`: 5 occurrences
- `Table Not Found (InvalidArgument)`: 10 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `departments`.`dept_id`: Expected type INT64; found STRING [at 1:1]
GENERATE_UUID()
^
2. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `employees`.`emp_id`: Expected type INT64; found STRING [at 1:1]
GENERATE_UUID()
^
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
4. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `projects`.`project_id`: Expected type INT64; found STRING [at 1:1]
GENERATE_UUID()
^
5. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
6. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
7. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
8. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
9. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error analyzing the definition of view `employee_details`: Table not found: employees [at 1:196]
...m.last_name AS manager_last_name FROM employees e LEFT JOIN departments d ...
                                         ^
10. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_id, dept_name, location) VALUES (NULL, 'Enginee...\n            ^", requestID = "1.c3438849fd14d720.3.1.5.1"
11. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.3.1.6.1"
12. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_id, dept_name, location) VALUES (NULL, 'Enginee...\n            ^", requestID = "1.c3438849fd14d720.3.1.7.1"
13. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.3.1.8.1"
14. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salar...\n            ^", requestID = "1.c3438849fd14d720.3.1.9.1"
15. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salar...\n            ^", requestID = "1.c3438849fd14d720.3.1.10.1"
16. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.3.1.11.1"
17. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.3.1.12.1"
18. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.3.1.13.1"

### fix-gpt-4o.sql

**Execution Error Rate**: 60.0% (6/10 parsed statements failed)

**Execution Error Codes**:
- `FailedPrecondition`: 2 occurrences
- `NotFound`: 2 occurrences
- `InvalidArgument`: 2 occurrences

**Error Categories**:
- `Table Not Found (InvalidArgument)`: 2 occurrences
- `FailedPrecondition`: 2 occurrences
- `NotFound`: 2 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = FailedPrecondition desc = Table employees does not reference parent key column dept_id.
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
4. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error analyzing the definition of view `employee_details`: Table not found: employees [at 1:196]
...m.last_name AS manager_last_name FROM employees e LEFT JOIN departments d ...
                                         ^
5. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: departments.dept_id in table: departments", requestID = "1.c3438849fd14d720.4.4.5.1"
6. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.4.4.6.1"

### fix-gpt-o3.sql

**Execution Error Rate**: 58.6% (17/29 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 12 occurrences
- `NotFound`: 5 occurrences

**Error Categories**:
- `NotFound`: 5 occurrences
- `Table Not Found (InvalidArgument)`: 6 occurrences
- `Unsupported Feature: Sequence Kind`: 3 occurrences
- `Function Not Found: NEXTVAL`: 3 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE SEQUENCE departments_seq OPTIONS (sequence_kind = 'positive') : Unsupported sequence kind: positive
2. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE SEQUENCE employees_seq OPTIONS (sequence_kind = 'positive') : Unsupported sequence kind: positive
3. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE SEQUENCE projects_seq OPTIONS (sequence_kind = 'positive') : Unsupported sequence kind: positive
4. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `departments`.`dept_id`: Function not found: NEXTVAL [at 1:1]
NEXTVAL(SEQUENCE departments_seq)
^
5. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `employees`.`emp_id`: Function not found: NEXTVAL [at 1:1]
NEXTVAL(SEQUENCE employees_seq)
^
6. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `projects`.`project_id`: Function not found: NEXTVAL [at 1:1]
NEXTVAL(SEQUENCE projects_seq)
^
7. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
8. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
9. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
10. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
11. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
12. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error analyzing the definition of view `employee_details`: Table not found: employees [at 1:196]
...m.last_name AS manager_last_name FROM employees AS e LEFT JOIN departments...
                                         ^
13. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_name, location) VALUES ('Engineering', 'New Yor...\n            ^", requestID = "1.c3438849fd14d720.5.3.5.1"
14. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_name, start_date, end_date, budget, status) VAL...\n            ^", requestID = "1.c3438849fd14d720.5.3.6.1"
15. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_...\n            ^", requestID = "1.c3438849fd14d720.5.3.7.1"
16. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.5.3.8.1"
17. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees AS e LEFT JO...\n                                                       ^", requestID = "1.c3438849fd14d720.5.3.9.1"

### fix-gpt-o4-mini-high.sql

**Execution Error Rate**: 8.7% (2/23 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 2 occurrences

**Error Categories**:
- `InvalidArgument: Other`: 2 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Value has type STRING which cannot be inserted into column dept_id, which has type INT64 [at 1:148]\n...Doe', 'test@example.com', '2024-01-01', 75000.0, '4611686018427387904') TH...\n                                                    ^", requestID = "1.c3438849fd14d720.6.3.9.1"
2. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Value has type STRING which cannot be inserted into column project_id, which has type INT64 [at 1:90]\n...project_id, role, hours_allocated) VALUES (NULL, '4611686018427387904', 'D...\n                                                    ^", requestID = "1.c3438849fd14d720.6.3.10.1"

### gpt-4o-2.sql

**Parse Error Rate**: 11.8% (2/17 failed to parse)

**Parse Error Types**:
- `Syntax Error: Expected Token`: 2 occurrences

**Parse Errors**:
- syntax error: gpt-4o-2.sql:5:32: expected token: (, but: <ident>
  Statement: `CREATE TABLE departments (
  dept_id INT64 NOT NULL,
  dept_name STRING(50) NOT NULL,
  location STRING(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
) PRIMARY KEY (dept_id)`

- syntax error: gpt-4o-2.sql:1:41: expected token: <ident>, but: AS
  Statement: `CREATE OR REPLACE VIEW employee_details AS
SELECT 
  e.emp_id,
  e.first_name,
  e.last_name,
 ...
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id`

**Execution Error Rate**: 26.7% (4/15 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 2 occurrences
- `InvalidArgument`: 2 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 2 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
2. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_id, dept_name, location) VALUES (NULL, 'Enginee...\n            ^", requestID = "1.c3438849fd14d720.7.2.5.1"
3. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:179]\n...p.project_name FROM employees e LEFT JOIN departments d ON e.dept_id = d.d...\n                                             ^", requestID = "1.c3438849fd14d720.7.2.6.1"
4. statement failed: executing ALTER statement: rpc error: code = NotFound desc = Table not found: departments

### gpt-4o-p2.sql

**Parse Error Rate**: 7.4% (2/27 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences
- `Syntax Error: Expected Token`: 1 occurrences

**Parse Errors**:
- syntax error: gpt-4o-p2.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id INT64 NOT NULL DEFAULT (GENERATE_UUID()),
  project_name STRI...status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

- syntax error: gpt-4o-p2.sql:1:30: expected token: <ident>, but: AS
  Statement: `CREATE VIEW employee_details AS
SELECT 
  e.emp_id,
  e.first_name,
  e.last_name,
  e.email,
...
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id`

**Execution Error Rate**: 64.0% (16/25 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 11 occurrences
- `NotFound`: 5 occurrences

**Error Categories**:
- `Type Mismatch: GENERATE_UUID on INT64`: 2 occurrences
- `NotFound`: 5 occurrences
- `Table Not Found (InvalidArgument)`: 9 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `departments`.`dept_id`: Expected type INT64; found STRING [at 1:1]
GENERATE_UUID()
^
2. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `employees`.`emp_id`: Expected type INT64; found STRING [at 1:1]
GENERATE_UUID()
^
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
5. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
6. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
7. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
8. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_id, dept_name, location) VALUES (NULL, 'Enginee...\n            ^", requestID = "1.c3438849fd14d720.8.3.5.1"
9. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.8.3.6.1"
10. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_id, dept_name, location) VALUES (NULL, 'Enginee...\n            ^", requestID = "1.c3438849fd14d720.8.3.7.1"
11. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.8.3.8.1"
12. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salar...\n            ^", requestID = "1.c3438849fd14d720.8.3.9.1"
13. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salar...\n            ^", requestID = "1.c3438849fd14d720.8.3.10.1"
14. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.8.3.11.1"
15. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.8.3.12.1"
16. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.8.3.13.1"

### gpt-4o-recommendations.sql

**Parse Error Rate**: 30.4% (7/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: General`: 4 occurrences
- `Syntax Error: Expected Token`: 3 occurrences

**Parse Errors**:
- syntax error: gpt-4o-recommendations.sql:6:29: unknown constraint PRIMARY
  Statement: `CREATE TABLE departments (
  dept_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  dept_name ST...ed_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP()),
  CONSTRAINT pk_departments PRIMARY KEY (dept_id)
)`

- syntax error: gpt-4o-recommendations.sql:11:27: unknown constraint PRIMARY
  Statement: `CREATE TABLE employees (
  emp_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  first_name STRI...ept_id),
  CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
)`

- syntax error: gpt-4o-recommendations.sql:10:26: unknown constraint PRIMARY
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name...CTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED')),
  CONSTRAINT pk_projects PRIMARY KEY (project_id)
)`

- syntax error: gpt-4o-recommendations.sql:6:37: unknown constraint PRIMARY
  Statement: `CREATE TABLE project_assignments (
  emp_id STRING(36) NOT NULL,
  project_id STRING(36) NOT NULL,...ees(emp_id),
  CONSTRAINT fk_pa_project FOREIGN KEY (project_id) REFERENCES projects(project_id)
)`

- syntax error: gpt-4o-recommendations.sql:3:1: expected token: <eof>, but: <ident>
  Statement: `INSERT INTO departments (dept_name, location)
VALUES (@dept_name, @location)
RETURNING dept_id`

- syntax error: gpt-4o-recommendations.sql:3:1: expected token: <eof>, but: <ident>
  Statement: `INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES (@first_name, @last_name, @email, @hire_date, @salary, @dept_id)
RETURNING emp_id`

- syntax error: gpt-4o-recommendations.sql:3:1: expected token: <eof>, but: <ident>
  Statement: `INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES (@project_name, @start_date, @end_date, @budget, @status)
RETURNING project_id`

**Execution Error Rate**: 43.8% (7/16 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 4 occurrences
- `InvalidArgument`: 3 occurrences

**Error Categories**:
- `Table Not Found (InvalidArgument)`: 3 occurrences
- `NotFound`: 4 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
5. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error analyzing the definition of view `employee_details`: Table not found: employees [at 1:196]
...m.last_name AS manager_last_name FROM employees e LEFT JOIN departments d ...
                                         ^
6. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.9.2.5.1"
7. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.9.2.6.1"

### gpt-4o-recommendations2.sql

**Parse Error Rate**: 32.3% (10/31 failed to parse)

**Parse Error Types**:
- `Syntax Error: General`: 2 occurrences
- `Syntax Error: Expected Token`: 8 occurrences

**Parse Errors**:
- syntax error: gpt-4o-recommendations2.sql:12:3: expected pseudo keyword: OPTIONS, but: FOREIGN
  Statement: `CREATE TABLE employees (
  emp_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  first_name STRI... (dept_id) REFERENCES departments(dept_id),
  FOREIGN KEY (manager_id) REFERENCES employees(emp_id)`

- syntax error: gpt-4o-recommendations2.sql:7:3: expected pseudo keyword: OPTIONS, but: FOREIGN
  Statement: `CREATE TABLE project_assignments (
  emp_id STRING(36) NOT NULL,
  project_id STRING(36) NOT NULL,...Y (emp_id) REFERENCES employees(emp_id),
  FOREIGN KEY (project_id) REFERENCES projects(project_id)`

- syntax error: gpt-4o-recommendations2.sql:1:1: unexpected token: <ident>
  Statement: `BEGIN TRANSACTION`

- syntax error: gpt-4o-recommendations2.sql:1:1: unexpected token: <ident>
  Statement: `COMMIT TRANSACTION`

- syntax error: gpt-4o-recommendations2.sql:1:1: unexpected token: <ident>
  Statement: `BEGIN TRANSACTION`

- syntax error: gpt-4o-recommendations2.sql:1:1: unexpected token: <ident>
  Statement: `COMMIT TRANSACTION`

- syntax error: gpt-4o-recommendations2.sql:1:1: unexpected token: <ident>
  Statement: `BEGIN TRANSACTION`

- syntax error: gpt-4o-recommendations2.sql:1:1: unexpected token: <ident>
  Statement: `COMMIT TRANSACTION`

- syntax error: gpt-4o-recommendations2.sql:1:1: unexpected token: <ident>
  Statement: `BEGIN TRANSACTION`

- syntax error: gpt-4o-recommendations2.sql:1:1: unexpected token: <ident>
  Statement: `COMMIT TRANSACTION`

**Execution Error Rate**: 28.6% (6/21 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 2 occurrences
- `InvalidArgument`: 4 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 4 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
3. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error analyzing the definition of view `employee_details`: Table not found: employees [at 1:196]
...m.last_name AS manager_last_name FROM employees e LEFT JOIN departments d ...
                                         ^
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salar...\n            ^", requestID = "1.c3438849fd14d720.10.3.11.1"
5. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.10.3.12.1"
6. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.10.3.13.1"

### gpt-5-3shot.sql

**Execution Error Rate**: 17.4% (4/23 parsed statements failed)

**Execution Error Codes**:
- `FailedPrecondition`: 4 occurrences

**Error Categories**:
- `FailedPrecondition`: 4 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: departments.dept_id in table: departments", requestID = "1.c3438849fd14d720.12.4.5.1"
2. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: projects.project_id in table: projects", requestID = "1.c3438849fd14d720.12.4.6.1"
3. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: employees.emp_id in table: employees", requestID = "1.c3438849fd14d720.12.4.7.1"
4. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: project_assignments.emp_id in table: project_assignments", requestID = "1.c3438849fd14d720.12.4.8.1"

### gpt-5.sql

**Parse Error Rate**: 21.7% (5/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences
- `Syntax Error: Expected Token`: 1 occurrences
- `Syntax Error: PRIMARY/FOREIGN KEY Placement`: 3 occurrences

**Parse Errors**:
- syntax error: gpt-5.sql:7:15: expected token: ), but: (
  Statement: `CREATE TABLE IF NOT EXISTS departments (
  dept_id   INT64 NOT NULL
            GENERATED BY DEFAU...ion  STRING(100),
  created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP()),
  PRIMARY KEY (dept_id)
)`

- syntax error: gpt-5.sql:12:15: expected token: ), but: (
  Statement: `CREATE TABLE IF NOT EXISTS employees (
  emp_id     INT64 NOT NULL
             GENERATED BY DEFAU...
  CONSTRAINT fk_emp_manager
    FOREIGN KEY (manager_id) REFERENCES employees(emp_id) ENFORCED
)`

- syntax error: gpt-5.sql:8:35: expected token: (, but: <string>
  Statement: `CREATE TABLE IF NOT EXISTS projects (
  project_id   INT64 NOT NULL
               GENERATED BY DE...ONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED')) ENFORCED
)`

- syntax error: gpt-5.sql:6:15: expected token: ), but: (
  Statement: `CREATE TABLE IF NOT EXISTS project_assignments (
  emp_id          INT64 NOT NULL,
  project_id   ...  CONSTRAINT fk_pa_project
    FOREIGN KEY (project_id) REFERENCES projects(project_id) ENFORCED
)`

- syntax error: gpt-5.sql:1:11: expected token: <ident>, but: IF
  Statement: `DROP VIEW IF EXISTS employee_details`

**Execution Error Rate**: 55.6% (10/18 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 4 occurrences
- `InvalidArgument`: 6 occurrences

**Error Categories**:
- `NotFound`: 4 occurrences
- `Table Not Found (InvalidArgument)`: 6 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
5. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error analyzing the definition of view `employee_details`: Table not found: employees [at 1:196]
...m.last_name AS manager_last_name FROM employees e LEFT JOIN departments d ...
                                         ^
6. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_name, location) VALUES ('Engineering', 'New Yor...\n            ^", requestID = "1.c3438849fd14d720.13.3.5.1"
7. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_name, start_date, end_date, budget, status) VAL...\n            ^", requestID = "1.c3438849fd14d720.13.3.6.1"
8. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_...\n            ^", requestID = "1.c3438849fd14d720.13.3.7.1"
9. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.13.3.8.1"
10. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.13.3.9.1"

### gpt-o3.sql

**Parse Error Rate**: 10.3% (3/29 failed to parse)

**Parse Error Types**:
- `Syntax Error: Expected Token`: 2 occurrences
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: gpt-o3.sql:5:42: expected token: (, but: <ident>
  Statement: `CREATE TABLE departments (
  dept_id     INT64 NOT NULL DEFAULT (NEXTVAL(SEQUENCE departments_seq))... STRING(100),
  created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP()
) PRIMARY KEY (dept_id)`

- syntax error: gpt-o3.sql:7:44: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id   INT64 NOT NULL DEFAULT (NEXTVAL(SEQUENCE projects_seq)),
  ...k_status  CHECK (status IN ('ACTIVE','COMPLETED','ON_HOLD','CANCELLED'))
) PRIMARY KEY (project_id)`

- syntax error: gpt-o3.sql:1:41: expected token: <ident>, but: AS
  Statement: `CREATE OR REPLACE VIEW employee_details AS
SELECT
  e.emp_id,
  e.first_name,
  e.last_name,
  ... departments AS d ON e.dept_id    = d.dept_id
LEFT JOIN employees   AS m ON e.manager_id = m.emp_id`

**Execution Error Rate**: 53.8% (14/26 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 5 occurrences
- `InvalidArgument`: 9 occurrences

**Error Categories**:
- `Unsupported Feature: Sequence Kind`: 3 occurrences
- `Function Not Found: NEXTVAL`: 1 occurrences
- `NotFound`: 5 occurrences
- `Table Not Found (InvalidArgument)`: 5 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE SEQUENCE departments_seq OPTIONS (sequence_kind = 'positive') : Unsupported sequence kind: positive
2. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE SEQUENCE employees_seq OPTIONS (sequence_kind = 'positive') : Unsupported sequence kind: positive
3. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing Spanner DDL statement: CREATE SEQUENCE projects_seq OPTIONS (sequence_kind = 'positive') : Unsupported sequence kind: positive
4. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `employees`.`emp_id`: Function not found: NEXTVAL [at 1:1]
NEXTVAL(SEQUENCE employees_seq)
^
5. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
6. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
7. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
8. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
9. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
10. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_name, location) VALUES ('Engineering', 'New Yor...\n            ^", requestID = "1.c3438849fd14d720.14.4.5.1"
11. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_name, start_date, end_date, budget, status) VAL...\n            ^", requestID = "1.c3438849fd14d720.14.4.6.1"
12. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_...\n            ^", requestID = "1.c3438849fd14d720.14.4.7.1"
13. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.14.4.8.1"
14. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees AS e LEFT JO...\n                                                       ^", requestID = "1.c3438849fd14d720.14.4.9.1"

### gpt-o4-mini-high.sql

**Parse Error Rate**: 30.4% (7/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: Expected Token`: 5 occurrences
- `Syntax Error: String Literal Quotes`: 1 occurrences
- `Syntax Error: PRIMARY/FOREIGN KEY Placement`: 1 occurrences

**Parse Errors**:
- syntax error: gpt-o4-mini-high.sql:5:44: expected token: (, but: <ident>
  Statement: `CREATE TABLE IF NOT EXISTS Departments (
  dept_id    INT64    NOT NULL GENERATED BY DEFAULT AS IDE... NULL,
  location   STRING(100),
  created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP()
)`

- syntax error: gpt-o4-mini-high.sql:7:45: expected token: (, but: <string>
  Statement: `CREATE TABLE IF NOT EXISTS Projects (
  project_id   INT64    NOT NULL GENERATED BY DEFAULT AS IDEN..._date),
  CONSTRAINT check_status CHECK (status IN ('ACTIVE','COMPLETED','ON_HOLD','CANCELLED'))
)`

- syntax error: gpt-o4-mini-high.sql:6:15: expected token: ), but: (
  Statement: `CREATE TABLE IF NOT EXISTS ProjectAssignments (
  emp_id          INT64    NOT NULL,
  project_id ...      REFERENCES Employees(emp_id),
  FOREIGN KEY (project_id)   REFERENCES Projects(project_id)
)`

- syntax error: gpt-o4-mini-high.sql:2:9: unexpected token: $ (and 1 other error)
  Statement: `INSERT INTO Departments (dept_name, location)
VALUES ($1, $2)
THEN RETURN dept_id`

- syntax error: gpt-o4-mini-high.sql:2:9: unexpected token: $ (and 5 other errors)
  Statement: `INSERT INTO Employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES ($1, $2, $3, $4, $5, $6)
THEN RETURN emp_id`

- syntax error: gpt-o4-mini-high.sql:2:9: unexpected token: $ (and 4 other errors)
  Statement: `INSERT INTO Projects (project_name, start_date, end_date, budget, status)
VALUES ($1, $2, $3, $4, $5)
THEN RETURN project_id`

- syntax error: gpt-o4-mini-high.sql:2:9: unexpected token: $ (and 3 other errors)
  Statement: `INSERT INTO ProjectAssignments (emp_id, project_id, role, hours_allocated)
VALUES ($1, $2, $3, $4)`

**Execution Error Rate**: 43.8% (7/16 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 3 occurrences
- `NotFound`: 4 occurrences

**Error Categories**:
- `Identity Column: Missing Sequence Kind`: 1 occurrences
- `NotFound`: 4 occurrences
- `Table Not Found (InvalidArgument)`: 2 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = The sequence kind of an identity column emp_id is not specified. Please specify the sequence kind explicitly or set the database option `default_sequence_kind`.
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: Employees
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: Employees
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: Departments
5. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: Projects
6. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error analyzing the definition of view `employee_details`: Table not found: Employees [at 1:255]
...m.last_name AS manager_last_name FROM Employees AS e LEFT JOIN Departments...
                                         ^
7. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: Employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM Employees AS e LEFT JO...\n                                                       ^", requestID = "1.c3438849fd14d720.15.1.5.1"

### instance-1-chatgpt-4o-latest-2-20251005112832.sql

**Parse Error Rate**: 4.3% (1/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: instance-1-chatgpt-4o-latest-2-20251005112832.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name S..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

**Execution Error Rate**: 22.7% (5/22 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 2 occurrences
- `InvalidArgument`: 3 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 3 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
3. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.16.3.8.1"
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.16.3.12.1"
5. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.16.3.13.1"

### instance-1-chatgpt-4o-latest-3-20251007203321.sql

**Parse Error Rate**: 4.3% (1/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: instance-1-chatgpt-4o-latest-3-20251007203321.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name S..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

**Execution Error Rate**: 31.8% (7/22 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 3 occurrences
- `NotFound`: 2 occurrences
- `FailedPrecondition`: 2 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `FailedPrecondition`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 3 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
3. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: departments.dept_id in table: departments", requestID = "1.c3438849fd14d720.17.2.5.1"
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.17.2.6.1"
5. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: employees.emp_id in table: employees", requestID = "1.c3438849fd14d720.17.2.7.1"
6. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.17.2.8.1"
7. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.17.2.9.1"

### instance-1-chatgpt-4o-latest-3-20251007203554.sql

**Parse Error Rate**: 4.3% (1/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: instance-1-chatgpt-4o-latest-3-20251007203554.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name S..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

**Execution Error Rate**: 31.8% (7/22 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 2 occurrences
- `FailedPrecondition`: 2 occurrences
- `InvalidArgument`: 3 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `FailedPrecondition`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 3 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
3. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: departments.dept_id in table: departments", requestID = "1.c3438849fd14d720.18.2.5.1"
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.18.2.6.1"
5. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: employees.emp_id in table: employees", requestID = "1.c3438849fd14d720.18.2.7.1"
6. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.18.2.8.1"
7. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.18.2.9.1"

### instance-1-chatgpt-4o-latest-3-20251007203716.sql

**Execution Error Rate**: 3.7% (1/27 parsed statements failed)

**Execution Error Codes**:
- `FailedPrecondition`: 1 occurrences

**Error Categories**:
- `FailedPrecondition`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: project_assignments.emp_id in table: project_assignments", requestID = "1.c3438849fd14d720.19.1.14.1"

### instance-1-chatgpt-4o-latest-3-20251007210819.sql

**Execution Error Rate**: 4.3% (1/23 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 1 occurrences

**Error Categories**:
- `Syntax Error: General`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Syntax error: Expected \")\" or \",\" but got identifier \"sample_value\" [at 1:118]\n...salary, dept_id ) VALUES ( 'John', 'Doe', 'john.doe'sample_value'.com', DA...\n                                                       ^", requestID = "1.c3438849fd14d720.20.4.11.1"

### instance-1-chatgpt-4o-latest-3-20251007213639.sql

**Execution Error Rate**: 8.0% (2/25 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 1 occurrences
- `FailedPrecondition`: 1 occurrences

**Error Categories**:
- `Syntax Error: General`: 1 occurrences
- `FailedPrecondition`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Syntax error: Expected \")\" or \",\" but got identifier \"sample_value\" [at 1:139]\n...VALUES ('UUID-EMP-001', 'John', 'Doe', 'john.doe'sample_value'.com', DATE ...\n                                                    ^", requestID = "1.c3438849fd14d720.21.2.11.1"
2. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Foreign key `FK_project_assignments_emp_id` constraint violation on table `project_assignments`. Cannot find referenced key `{String(\"UUID-EMP-001\")}` in table `employees`.", requestID = "1.c3438849fd14d720.21.2.12.1"

### instance-1-chatgpt-4o-latest-3-20251007213715.sql

**Execution Error Rate**: 11.1% (3/27 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 1 occurrences
- `FailedPrecondition`: 2 occurrences

**Error Categories**:
- `Syntax Error: General`: 1 occurrences
- `FailedPrecondition`: 2 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Syntax error: Expected \")\" or \",\" but got identifier \"sample_value\" [at 1:115]\n...salary, dept_id) VALUES ('John', 'Doe', 'john.doe'sample_value'.com', DATE...\n                                                     ^", requestID = "1.c3438849fd14d720.22.4.11.1"
2. statement failed: executing ALTER statement: rpc error: code = FailedPrecondition desc = Foreign key `fk_emp` constraint violation on table `project_assignments`. Cannot find referenced key `{String("sample-emp-id")}` in table `employees`.
3. statement failed: executing ALTER statement: rpc error: code = FailedPrecondition desc = Foreign key `fk_project` constraint violation on table `project_assignments`. Cannot find referenced key `{String("sample-project-id")}` in table `projects`.

### instance-1-chatgpt-4o-latest-3-20251007213745.sql

**Execution Error Rate**: 4.3% (1/23 parsed statements failed)

**Execution Error Codes**:
- `FailedPrecondition`: 1 occurrences

**Error Categories**:
- `FailedPrecondition`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: project_assignments.emp_id in table: project_assignments", requestID = "1.c3438849fd14d720.23.4.14.1"

### instance-1-chatgpt-4o-latest-3-20251007213927.sql

**Execution Error Rate**: 6.9% (2/29 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 1 occurrences
- `FailedPrecondition`: 1 occurrences

**Error Categories**:
- `Syntax Error: General`: 1 occurrences
- `FailedPrecondition`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Syntax error: Expected \")\" or \",\" but got identifier \"sample_value\" [at 1:139]\n...VALUES ('uuid_emp_001', 'John', 'Doe', 'john.doe'sample_value'.com', DATE ...\n                                                    ^", requestID = "1.c3438849fd14d720.24.3.11.1"
2. statement failed: executing ALTER statement: rpc error: code = FailedPrecondition desc = Foreign key `fk_project_assignments_emp` constraint violation on table `project_assignments`. Cannot find referenced key `{String("uuid_emp_001")}` in table `employees`.

### instance-1-chatgpt-4o-latest-3-20251007221438.sql

**Parse Error Rate**: 5.3% (1/19 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: instance-1-chatgpt-4o-latest-3-20251007221438.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name S..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

**Execution Error Rate**: 16.7% (3/18 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 2 occurrences
- `InvalidArgument`: 1 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 1 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
3. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.25.2.5.1"

### instance-1-gpt-4o-mini-3-20251009005208.sql

**Parse Error Rate**: 34.8% (8/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: Expected Token`: 6 occurrences
- `Syntax Error: String Literal Quotes`: 1 occurrences
- `Syntax Error: PRIMARY/FOREIGN KEY Placement`: 1 occurrences

**Parse Errors**:
- syntax error: instance-1-gpt-4o-mini-3-20251009005208.sql:5:34: expected token: (, but: <ident>
  Statement: `CREATE TABLE departments (
    dept_id STRING(36) NOT NULL PRIMARY KEY,
    dept_name STRING(50) NOT NULL,
    location STRING(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
)`

- syntax error: instance-1-gpt-4o-mini-3-20251009005208.sql:7:31: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
    project_id STRING(36) NOT NULL PRIMARY KEY,
    project_name STRING(100)...te),
    CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
)`

- syntax error: instance-1-gpt-4o-mini-3-20251009005208.sql:6:17: expected token: ), but: (
  Statement: `CREATE TABLE project_assignments (
    emp_id STRING(36),
    project_id STRING(36),
    role STRING...emp_id) REFERENCES employees(emp_id),
    FOREIGN KEY (project_id) REFERENCES projects(project_id)
)`

- syntax error: instance-1-gpt-4o-mini-3-20251009005208.sql:1:41: expected token: <ident>, but: AS
  Statement: `CREATE OR REPLACE VIEW employee_details AS
SELECT 
    e.emp_id,
    e.first_name,
    e.last_name,
...nts d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
SQL SECURITY INVOKER`

- syntax error: instance-1-gpt-4o-mini-3-20251009005208.sql:2:26: unexpected token: $ (and 2 other errors)
  Statement: `INSERT INTO departments (dept_id, dept_name, location)
VALUES (GENERATE_UUID(), $1, $2)
RETURN dept_id`

- syntax error: instance-1-gpt-4o-mini-3-20251009005208.sql:2:26: unexpected token: $ (and 6 other errors)
  Statement: `INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id)
VALUES (GENERATE_UUID(), $1, $2, $3, $4, $5, $6)
RETURN emp_id`

- syntax error: instance-1-gpt-4o-mini-3-20251009005208.sql:2:26: unexpected token: $ (and 4 other errors)
  Statement: `INSERT INTO projects (project_id, project_name, start_date, end_date, budget, status)
VALUES (GENERATE_UUID(), $1, $2, $3, $4, 'ACTIVE')
RETURN project_id`

- syntax error: instance-1-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 3 other errors)
  Statement: `INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES ($1, $2, $3, $4)`

**Execution Error Rate**: 40.0% (6/15 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 5 occurrences
- `InvalidArgument`: 1 occurrences

**Error Categories**:
- `NotFound`: 5 occurrences
- `Table Not Found (InvalidArgument)`: 1 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
5. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
6. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.28.2.5.1"

### instance-1-gpt-5-2025-08-07-3-20251006235414.sql

**Execution Error Rate**: 4.3% (1/23 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 1 occurrences

**Error Categories**:
- `InvalidArgument: Other`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Query without FROM clause cannot have a WHERE clause [at 1:212]\n...4289-a094-786fd1561b6f' AS STRING), 'Developer', 40 WHERE '6633877e-1899-4...\n                                                       ^", requestID = "1.c3438849fd14d720.29.3.11.1"

### instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql

**Parse Error Rate**: 34.8% (8/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: Expected Token`: 7 occurrences
- `Syntax Error: PRIMARY/FOREIGN KEY Placement`: 1 occurrences

**Parse Errors**:
- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql:2:40: expected token: ), but: <ident>
  Statement: `CREATE TABLE departments (
    dept_id INT64 NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY (STAR...TRING(50) NOT NULL,
    location STRING(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
)`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql:2:39: expected token: ), but: <ident>
  Statement: `CREATE TABLE employees (
    emp_id INT64 NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY (START W...ept_id),
    CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
)`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql:2:43: expected token: ), but: <ident>
  Statement: `CREATE TABLE projects (
    project_id INT64 NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY (STAR...te),
    CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
)`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql:6:17: expected token: ), but: (
  Statement: `CREATE TABLE project_assignments (
    emp_id INT64 NOT NULL,
    project_id INT64 NOT NULL,
    rol...ees(emp_id),
    CONSTRAINT fk_pa_project FOREIGN KEY (project_id) REFERENCES projects(project_id)
)`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql:1:41: expected token: <ident>, but: AS
  Statement: `CREATE OR REPLACE VIEW employee_details AS
SELECT 
    e.emp_id,
    e.first_name,
    e.last_name,
... e
LEFT JOIN departments d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql:3:1: expected token: <eof>, but: <ident>
  Statement: `INSERT INTO departments (dept_name, location)
VALUES (@p1, @p2)
RETURNING dept_id`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql:3:1: expected token: <eof>, but: <ident>
  Statement: `INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES (@p1, @p2, @p3, @p4, @p5, @p6)
RETURNING emp_id`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233257.sql:3:1: expected token: <eof>, but: <ident>
  Statement: `INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES (@p1, @p2, @p3, @p4, @p5)
RETURNING project_id`

**Execution Error Rate**: 40.0% (6/15 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 4 occurrences
- `InvalidArgument`: 2 occurrences

**Error Categories**:
- `Table Not Found (InvalidArgument)`: 2 occurrences
- `NotFound`: 4 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
5. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.33.1.5.1"
6. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.33.1.6.1"

### instance-1-gpt-5-mini-2025-08-07-1-20251007233356.sql

**Parse Error Rate**: 19.2% (5/26 failed to parse)

**Parse Error Types**:
- `Syntax Error: Expected Token`: 4 occurrences
- `Syntax Error: PRIMARY/FOREIGN KEY Placement`: 1 occurrences

**Parse Errors**:
- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233356.sql:2:36: expected token: BY, but: <ident>
  Statement: `CREATE TABLE departments (
  dept_id INT64 NOT NULL GENERATED ALWAYS AS IDENTITY (START WITH 1 INCRE...n STRING(100),
  created_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP())
) PRIMARY KEY (dept_id)`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233356.sql:2:35: expected token: BY, but: <ident>
  Statement: `CREATE TABLE employees (
  emp_id INT64 NOT NULL GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMEN...NT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
) PRIMARY KEY (emp_id)`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233356.sql:2:39: expected token: BY, but: <ident>
  Statement: `CREATE TABLE projects (
  project_id INT64 NOT NULL GENERATED ALWAYS AS IDENTITY (START WITH 1 INCRE..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233356.sql:6:15: expected token: ), but: (
  Statement: `CREATE TABLE project_assignments (
  emp_id INT64 NOT NULL,
  project_id INT64 NOT NULL,
  role STRI... (emp_id) REFERENCES employees(emp_id),
  FOREIGN KEY (project_id) REFERENCES projects(project_id)
)`

- syntax error: instance-1-gpt-5-mini-2025-08-07-1-20251007233356.sql:1:30: expected token: <ident>, but: AS
  Statement: `CREATE VIEW employee_details AS
SELECT
  e.emp_id,
  e.first_name,
  e.last_name,
  e.email,
  d.dep...T JOIN departments AS d ON e.dept_id = d.dept_id
LEFT JOIN employees AS m ON e.manager_id = m.emp_id`

**Execution Error Rate**: 57.1% (12/21 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 4 occurrences
- `InvalidArgument`: 8 occurrences

**Error Categories**:
- `NotFound`: 4 occurrences
- `Table Not Found (InvalidArgument)`: 5 occurrences
- `Missing Clause: General`: 3 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
5. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_name, location) VALUES ('sample_value', 'sample...\n            ^", requestID = "1.c3438849fd14d720.34.4.5.1"
6. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_name, start_date, end_date, budget, status) VAL...\n            ^", requestID = "1.c3438849fd14d720.34.4.6.1"
7. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_...\n            ^", requestID = "1.c3438849fd14d720.34.4.7.1"
8. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.34.4.8.1"
9. SELECT failed: executing SELECT: spanner: code = "InvalidArgument", desc = "missing value for query parameter p1"
10. SELECT failed: executing SELECT: spanner: code = "InvalidArgument", desc = "missing value for query parameter p3"
11. SELECT failed: executing SELECT: spanner: code = "InvalidArgument", desc = "missing value for query parameter p1"
12. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees AS e LEFT JO...\n                                                       ^", requestID = "1.c3438849fd14d720.34.4.9.1"

### instance-2-chatgpt-4o-latest-2-20251005112832.sql

**Parse Error Rate**: 4.3% (1/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: instance-2-chatgpt-4o-latest-2-20251005112832.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name S..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

**Execution Error Rate**: 22.7% (5/22 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 2 occurrences
- `InvalidArgument`: 3 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 3 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
3. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.35.2.8.1"
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.35.2.12.1"
5. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.35.2.13.1"

### instance-2-chatgpt-4o-latest-3-20251007203405.sql

**Parse Error Rate**: 4.3% (1/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: instance-2-chatgpt-4o-latest-3-20251007203405.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name S..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

**Execution Error Rate**: 31.8% (7/22 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 2 occurrences
- `FailedPrecondition`: 2 occurrences
- `InvalidArgument`: 3 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `FailedPrecondition`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 3 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
3. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: departments.dept_id in table: departments", requestID = "1.c3438849fd14d720.36.1.5.1"
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.36.1.6.1"
5. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: employees.emp_id in table: employees", requestID = "1.c3438849fd14d720.36.1.7.1"
6. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.36.1.8.1"
7. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.36.1.9.1"

### instance-2-chatgpt-4o-latest-3-20251007203554.sql

**Execution Error Rate**: 4.3% (1/23 parsed statements failed)

**Execution Error Codes**:
- `FailedPrecondition`: 1 occurrences

**Error Categories**:
- `FailedPrecondition`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: project_assignments.emp_id in table: project_assignments", requestID = "1.c3438849fd14d720.37.2.14.1"

### instance-2-gpt-4o-mini-3-20251009005208.sql

**Parse Error Rate**: 39.1% (9/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: Expected Token`: 8 occurrences
- `Syntax Error: PRIMARY/FOREIGN KEY Placement`: 1 occurrences

**Parse Errors**:
- syntax error: instance-2-gpt-4o-mini-3-20251009005208.sql:2:41: expected token: (, but: <ident>
  Statement: `CREATE TABLE departments (
    dept_id STRING(36) NOT NULL DEFAULT GENERATE_UUID(),
    dept_name ST...ation STRING(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (dept_id)
)`

- syntax error: instance-2-gpt-4o-mini-3-20251009005208.sql:2:40: expected token: (, but: <ident>
  Statement: `CREATE TABLE employees (
    emp_id STRING(36) NOT NULL DEFAULT GENERATE_UUID(),
    first_name STRI...ept_id) REFERENCES departments(dept_id),
    FOREIGN KEY (manager_id) REFERENCES employees(emp_id)
)`

- syntax error: instance-2-gpt-4o-mini-3-20251009005208.sql:2:44: expected token: (, but: <ident>
  Statement: `CREATE TABLE projects (
    project_id STRING(36) NOT NULL DEFAULT GENERATE_UUID(),
    project_name...us CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED')),
    PRIMARY KEY (project_id)
)`

- syntax error: instance-2-gpt-4o-mini-3-20251009005208.sql:6:17: expected token: ), but: (
  Statement: `CREATE TABLE project_assignments (
    emp_id STRING(36),
    project_id STRING(36),
    role STRING...emp_id) REFERENCES employees(emp_id),
    FOREIGN KEY (project_id) REFERENCES projects(project_id)
)`

- syntax error: instance-2-gpt-4o-mini-3-20251009005208.sql:1:41: expected token: <ident>, but: AS
  Statement: `CREATE OR REPLACE VIEW employee_details AS
SELECT 
    e.emp_id,
    e.first_name,
    e.last_name,
...nts d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
SQL SECURITY INVOKER`

- syntax error: instance-2-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 2 other errors)
  Statement: `INSERT INTO departments (dept_name, location)
VALUES ($1, $2)
RETURNING dept_id`

- syntax error: instance-2-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 6 other errors)
  Statement: `INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING emp_id`

- syntax error: instance-2-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 5 other errors)
  Statement: `INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING project_id`

- syntax error: instance-2-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 3 other errors)
  Statement: `INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES ($1, $2, $3, $4)`

**Execution Error Rate**: 35.7% (5/14 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 1 occurrences
- `NotFound`: 4 occurrences

**Error Categories**:
- `NotFound`: 4 occurrences
- `Table Not Found (InvalidArgument)`: 1 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
5. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name as manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.38.2.5.1"

### instance-2-gpt-5-2025-08-07-3-20251007000603.sql

**Execution Error Rate**: 8.7% (2/23 parsed statements failed)

**Execution Error Codes**:
- `FailedPrecondition`: 2 occurrences

**Error Categories**:
- `FailedPrecondition`: 2 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "FailedPrecondition", desc = "Foreign key `FK_employees_departments_F23FD64885E5D64D_1` constraint violation on table `employees`. Cannot find referenced key `{String(\"sample_value\")}` in table `departments`.", requestID = "1.c3438849fd14d720.40.3.9.1"
2. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Foreign key `FK_project_assignments_projects_E405C4ABFCFDD1E6_1` constraint violation on table `project_assignments`. Cannot find referenced key `{String(\"sample_value\")}` in table `projects`.", requestID = "1.c3438849fd14d720.40.3.10.1"

### instance-3-chatgpt-4o-latest-2-20251005112832.sql

**Parse Error Rate**: 4.3% (1/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: instance-3-chatgpt-4o-latest-2-20251005112832.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name S..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

**Execution Error Rate**: 22.7% (5/22 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 2 occurrences
- `InvalidArgument`: 3 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 3 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
3. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.43.1.8.1"
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.43.1.12.1"
5. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.43.1.13.1"

### instance-3-gpt-4o-mini-3-20251009005208.sql

**Parse Error Rate**: 34.8% (8/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: Expected Token`: 6 occurrences
- `Syntax Error: String Literal Quotes`: 1 occurrences
- `Syntax Error: PRIMARY/FOREIGN KEY Placement`: 1 occurrences

**Parse Errors**:
- syntax error: instance-3-gpt-4o-mini-3-20251009005208.sql:5:34: expected token: (, but: <ident>
  Statement: `CREATE TABLE departments (
    dept_id STRING(36) NOT NULL PRIMARY KEY,
    dept_name STRING(50) NOT NULL,
    location STRING(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
)`

- syntax error: instance-3-gpt-4o-mini-3-20251009005208.sql:7:31: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
    project_id STRING(36) NOT NULL PRIMARY KEY,
    project_name STRING(100)...te),
    CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
)`

- syntax error: instance-3-gpt-4o-mini-3-20251009005208.sql:6:17: expected token: ), but: (
  Statement: `CREATE TABLE project_assignments (
    emp_id STRING(36),
    project_id STRING(36),
    role STRING...emp_id) REFERENCES employees(emp_id),
    FOREIGN KEY (project_id) REFERENCES projects(project_id)
)`

- syntax error: instance-3-gpt-4o-mini-3-20251009005208.sql:1:41: expected token: <ident>, but: AS
  Statement: `CREATE OR REPLACE VIEW employee_details AS
SELECT 
    e.emp_id,
    e.first_name,
    e.last_name,
...nts d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
SQL SECURITY INVOKER`

- syntax error: instance-3-gpt-4o-mini-3-20251009005208.sql:2:26: unexpected token: $ (and 2 other errors)
  Statement: `INSERT INTO departments (dept_id, dept_name, location)
VALUES (GENERATE_UUID(), $1, $2)
RETURNING dept_id`

- syntax error: instance-3-gpt-4o-mini-3-20251009005208.sql:2:26: unexpected token: $ (and 6 other errors)
  Statement: `INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id)
VALUES (GENERATE_UUID(), $1, $2, $3, $4, $5, $6)
RETURNING emp_id`

- syntax error: instance-3-gpt-4o-mini-3-20251009005208.sql:2:26: unexpected token: $ (and 4 other errors)
  Statement: `INSERT INTO projects (project_id, project_name, start_date, end_date, budget, status)
VALUES (GENERATE_UUID(), $1, $2, $3, $4, 'ACTIVE')
RETURNING project_id`

- syntax error: instance-3-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 3 other errors)
  Statement: `INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES ($1, $2, $3, $4)`

**Execution Error Rate**: 40.0% (6/15 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 5 occurrences
- `InvalidArgument`: 1 occurrences

**Error Categories**:
- `NotFound`: 5 occurrences
- `Table Not Found (InvalidArgument)`: 1 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
5. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
6. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.44.2.5.1"

### instance-3-gpt-5-2025-08-07-3-20251006235414.sql

**Execution Error Rate**: 8.7% (2/23 parsed statements failed)

**Execution Error Codes**:
- `OutOfRange`: 1 occurrences
- `InvalidArgument`: 1 occurrences

**Error Categories**:
- `OutOfRange`: 1 occurrences
- `InvalidArgument: Other`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "OutOfRange", desc = "Check constraint `projects`.`check_status` is violated for key {String(\"89f0a41e-49d5-4b6c-94f4-536e25c63bd2\")}", requestID = "1.c3438849fd14d720.45.1.7.1"
2. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Query without FROM clause cannot have a WHERE clause [at 1:425]\n...sample_value', SAFE_CAST('sample_value' AS INT64) WHERE EXISTS (SELECT 1 F...\n                                                     ^", requestID = "1.c3438849fd14d720.45.1.10.1"

### instance-3-gpt-5-2025-08-07-3-20251007002717.sql

**Execution Error Rate**: 53.8% (14/26 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 9 occurrences
- `NotFound`: 5 occurrences

**Error Categories**:
- `Default Value: Parsing Error`: 3 occurrences
- `NotFound`: 5 occurrences
- `Table Not Found (InvalidArgument)`: 6 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `departments`.`dept_id`: Unrecognized name: seq_departments_dept_id [at 1:25]
GET_NEXT_SEQUENCE_VALUE(seq_departments_dept_id)
                        ^
2. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `employees`.`emp_id`: Unrecognized name: seq_employees_emp_id [at 1:25]
GET_NEXT_SEQUENCE_VALUE(seq_employees_emp_id)
                        ^
3. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error parsing the default value of column `projects`.`project_id`: Unrecognized name: seq_projects_project_id [at 1:25]
GET_NEXT_SEQUENCE_VALUE(seq_projects_project_id)
                        ^
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
5. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
6. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
7. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
8. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
9. CREATE failed: executing CREATE statement: rpc error: code = InvalidArgument desc = Error analyzing the definition of view `employee_details`: Table not found: employees [at 1:196]
...m.last_name AS manager_last_name FROM employees e LEFT JOIN departments d ...
                                         ^
10. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: departments [at 1:13]\nINSERT INTO departments (dept_name, location) VALUES ('Engineering', 'New Yor...\n            ^", requestID = "1.c3438849fd14d720.48.3.5.1"
11. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_name, start_date, end_date, budget, status) VAL...\n            ^", requestID = "1.c3438849fd14d720.48.3.6.1"
12. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:13]\nINSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_...\n            ^", requestID = "1.c3438849fd14d720.48.3.7.1"
13. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.48.3.8.1"
14. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.48.3.9.1"

### instance-4-chatgpt-4o-latest-2-20251005112832.sql

**Parse Error Rate**: 4.3% (1/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: instance-4-chatgpt-4o-latest-2-20251005112832.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name S..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

**Execution Error Rate**: 27.3% (6/22 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 2 occurrences
- `InvalidArgument`: 4 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 3 occurrences
- `Syntax Error: General`: 1 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
3. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.49.4.8.1"
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Syntax error: Expected \")\" or \",\" but got identifier \"sample_value\" [at 1:132]\n...dept_id) VALUES ('UUID2', 'John', 'Doe', 'john.doe'sample_value'.com', '20...\n                                                      ^", requestID = "1.c3438849fd14d720.49.4.9.1"
5. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.49.4.10.1"
6. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.49.4.11.1"

### instance-4-gpt-4o-mini-3-20251009005208.sql

**Parse Error Rate**: 39.1% (9/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: Expected Token`: 8 occurrences
- `Syntax Error: PRIMARY/FOREIGN KEY Placement`: 1 occurrences

**Parse Errors**:
- syntax error: instance-4-gpt-4o-mini-3-20251009005208.sql:2:41: expected token: (, but: <ident>
  Statement: `CREATE TABLE departments (
    dept_id STRING(36) NOT NULL DEFAULT GENERATE_UUID(),
    dept_name ST...ation STRING(MAX),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (dept_id)
)`

- syntax error: instance-4-gpt-4o-mini-3-20251009005208.sql:2:40: expected token: (, but: <ident>
  Statement: `CREATE TABLE employees (
    emp_id STRING(36) NOT NULL DEFAULT GENERATE_UUID(),
    first_name STRI...ents(dept_id),
    FOREIGN KEY (manager_id) REFERENCES employees(emp_id),
    PRIMARY KEY (emp_id)
)`

- syntax error: instance-4-gpt-4o-mini-3-20251009005208.sql:2:44: expected token: (, but: <ident>
  Statement: `CREATE TABLE projects (
    project_id STRING(36) NOT NULL DEFAULT GENERATE_UUID(),
    project_name...us CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED')),
    PRIMARY KEY (project_id)
)`

- syntax error: instance-4-gpt-4o-mini-3-20251009005208.sql:6:17: expected token: ), but: (
  Statement: `CREATE TABLE project_assignments (
    emp_id STRING(36),
    project_id STRING(36),
    role STRING...emp_id) REFERENCES employees(emp_id),
    FOREIGN KEY (project_id) REFERENCES projects(project_id)
)`

- syntax error: instance-4-gpt-4o-mini-3-20251009005208.sql:1:41: expected token: <ident>, but: AS
  Statement: `CREATE OR REPLACE VIEW employee_details AS
SELECT 
    e.emp_id,
    e.first_name,
    e.last_name,
...nts d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
SQL SECURITY INVOKER`

- syntax error: instance-4-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 2 other errors)
  Statement: `INSERT INTO departments (dept_name, location)
VALUES ($1, $2)
RETURNING dept_id`

- syntax error: instance-4-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 6 other errors)
  Statement: `INSERT INTO employees (first_name, last_name, email, hire_date, salary, dept_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING emp_id`

- syntax error: instance-4-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 5 other errors)
  Statement: `INSERT INTO projects (project_name, start_date, end_date, budget, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING project_id`

- syntax error: instance-4-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 3 other errors)
  Statement: `INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES ($1, $2, $3, $4)`

**Execution Error Rate**: 35.7% (5/14 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 4 occurrences
- `InvalidArgument`: 1 occurrences

**Error Categories**:
- `NotFound`: 4 occurrences
- `Table Not Found (InvalidArgument)`: 1 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: employees
3. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: departments
4. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
5. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: employees [at 1:157]\n...last_name AS manager_last_name, p.project_name FROM employees e LEFT JOIN ...\n                                                       ^", requestID = "1.c3438849fd14d720.50.2.5.1"

### instance-4-gpt-5-2025-08-07-3-20251007000603.sql

**Execution Error Rate**: 8.7% (2/23 parsed statements failed)

**Execution Error Codes**:
- `FailedPrecondition`: 2 occurrences

**Error Categories**:
- `FailedPrecondition`: 2 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "FailedPrecondition", desc = "Foreign key `fk_emp_dept` constraint violation on table `employees`. Cannot find referenced key `{String(\"sample_value\")}` in table `departments`.", requestID = "1.c3438849fd14d720.52.4.9.1"
2. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Foreign key `fk_pa_proj` constraint violation on table `project_assignments`. Cannot find referenced key `{String(\"sample_value\")}` in table `projects`.", requestID = "1.c3438849fd14d720.52.4.10.1"

### instance-4-gpt-5-2025-08-07-3-20251007002717.sql

**Execution Error Rate**: 13.0% (3/23 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 3 occurrences

**Error Categories**:
- `InvalidArgument: Other`: 3 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Could not cast literal \"sample_value\" to type DATE [at 1:99]\n...end_date, budget, status) VALUES ('sample_value', 'sample_value', 'sample_...\n                                                     ^", requestID = "1.c3438849fd14d720.54.4.7.1"
2. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "InvalidArgument", desc = "Could not cast literal \"sample_value\" to type DATE [at 1:138]\n...sample_value', 'sample_value', 'sample_value', 'sample_value', 'sample_val...\n                                                  ^", requestID = "1.c3438849fd14d720.54.4.8.1"
3. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Value has type STRING which cannot be inserted into column hours_allocated, which has type INT64 [at 1:133]\n...sample_value', 'sample_value', 'sample_value', 'sample_value')\n                                                  ^", requestID = "1.c3438849fd14d720.54.4.9.1"

### instance-5-chatgpt-4o-latest-2-20251005112832.sql

**Execution Error Rate**: 8.7% (2/23 parsed statements failed)

**Execution Error Codes**:
- `InvalidArgument`: 1 occurrences
- `FailedPrecondition`: 1 occurrences

**Error Categories**:
- `Syntax Error: General`: 1 occurrences
- `FailedPrecondition`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Syntax error: Expected \")\" or \",\" but got identifier \"sample_value\" [at 1:133]\n...dept_id) VALUES ('<UUID>', 'John', 'Doe', 'john.doe'sample_value'.com', '2...\n                                                       ^", requestID = "1.c3438849fd14d720.55.1.11.1"
2. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Foreign key `FK_project_assignments_projects_E405C4ABFCFDD1E6_1` constraint violation on table `project_assignments`. Cannot find referenced key `{String(\"<PROJECT_UUID>\")}` in table `projects`.", requestID = "1.c3438849fd14d720.55.1.12.1"

### instance-5-chatgpt-4o-latest-3-20251007203405.sql

**Parse Error Rate**: 4.3% (1/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: instance-5-chatgpt-4o-latest-3-20251007203405.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name S..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

**Execution Error Rate**: 31.8% (7/22 parsed statements failed)

**Execution Error Codes**:
- `FailedPrecondition`: 2 occurrences
- `InvalidArgument`: 3 occurrences
- `NotFound`: 2 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `FailedPrecondition`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 3 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
3. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: departments.dept_id in table: departments", requestID = "1.c3438849fd14d720.56.4.5.1"
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_id, project_name, start_date, end_date, budget,...\n            ^", requestID = "1.c3438849fd14d720.56.4.6.1"
5. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: employees.emp_id in table: employees", requestID = "1.c3438849fd14d720.56.4.7.1"
6. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.56.4.8.1"
7. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.56.4.9.1"

### instance-5-gpt-4o-mini-3-20251009005208.sql

**Parse Error Rate**: 30.4% (7/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences
- `Syntax Error: PRIMARY/FOREIGN KEY Placement`: 1 occurrences
- `Syntax Error: Expected Token`: 5 occurrences

**Parse Errors**:
- syntax error: instance-5-gpt-4o-mini-3-20251009005208.sql:7:31: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
    project_id STRING(36) NOT NULL PRIMARY KEY,
    project_name STRING(100)...te),
    CONSTRAINT check_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
)`

- syntax error: instance-5-gpt-4o-mini-3-20251009005208.sql:6:17: expected token: ), but: (
  Statement: `CREATE TABLE project_assignments (
    emp_id STRING(36),
    project_id STRING(36),
    role STRING...emp_id) REFERENCES employees(emp_id),
    FOREIGN KEY (project_id) REFERENCES projects(project_id)
)`

- syntax error: instance-5-gpt-4o-mini-3-20251009005208.sql:1:41: expected token: <ident>, but: AS
  Statement: `CREATE OR REPLACE VIEW employee_details AS
SELECT 
    e.emp_id,
    e.first_name,
    e.last_name,
...nts d ON e.dept_id = d.dept_id
LEFT JOIN employees m ON e.manager_id = m.emp_id
SQL SECURITY INVOKER`

- syntax error: instance-5-gpt-4o-mini-3-20251009005208.sql:2:26: unexpected token: $ (and 2 other errors)
  Statement: `INSERT INTO departments (dept_id, dept_name, location)
VALUES (GENERATE_UUID(), $1, $2)
RETURNING dept_id`

- syntax error: instance-5-gpt-4o-mini-3-20251009005208.sql:2:26: unexpected token: $ (and 6 other errors)
  Statement: `INSERT INTO employees (emp_id, first_name, last_name, email, hire_date, salary, dept_id)
VALUES (GENERATE_UUID(), $1, $2, $3, $4, $5, $6)
RETURNING emp_id`

- syntax error: instance-5-gpt-4o-mini-3-20251009005208.sql:2:26: unexpected token: $ (and 4 other errors)
  Statement: `INSERT INTO projects (project_id, project_name, start_date, end_date, budget, status)
VALUES (GENERATE_UUID(), $1, $2, $3, $4, 'ACTIVE')
RETURNING project_id`

- syntax error: instance-5-gpt-4o-mini-3-20251009005208.sql:2:9: unexpected token: $ (and 3 other errors)
  Statement: `INSERT INTO project_assignments (emp_id, project_id, role, hours_allocated)
VALUES ($1, $2, $3, $4)`

**Execution Error Rate**: 12.5% (2/16 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 1 occurrences
- `InvalidArgument`: 1 occurrences

**Error Categories**:
- `NotFound`: 1 occurrences
- `Table Not Found (InvalidArgument)`: 1 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.57.4.5.1"

### instance-5-gpt-5-2025-08-07-3-20251007000603.sql

**Execution Error Rate**: 10.0% (3/30 parsed statements failed)

**Execution Error Codes**:
- `FailedPrecondition`: 1 occurrences
- `InvalidArgument`: 2 occurrences

**Error Categories**:
- `InvalidArgument: Other`: 2 occurrences
- `FailedPrecondition`: 1 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Query column 4 has type STRING which cannot be inserted into column budget, which has type NUMERIC [at 1:75]\n...project_name, start_date, end_date, budget, status) SELECT 'sample_value',...\n                                                       ^", requestID = "1.c3438849fd14d720.59.4.8.1"
2. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Could not cast literal \"sample_value\" to type DATE [at 1:138]\n...sample_value', 'sample_value', 'sample_value', 'sample_value', 'sample_val...\n                                                  ^", requestID = "1.c3438849fd14d720.59.4.9.1"
3. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Cannot specify a null value for column: project_assignments.emp_id in table: project_assignments", requestID = "1.c3438849fd14d720.59.4.10.1"

### instance-5-gpt-5-2025-08-07-3-20251007002717.sql

**Execution Error Rate**: 13.0% (3/23 parsed statements failed)

**Execution Error Codes**:
- `OutOfRange`: 1 occurrences
- `FailedPrecondition`: 2 occurrences

**Error Categories**:
- `OutOfRange`: 1 occurrences
- `FailedPrecondition`: 2 occurrences

**Execution Errors**:
1. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "OutOfRange", desc = "Check constraint `projects`.`check_status` is violated for key {String(\"b096c561-fde8-49f9-bc2e-93043626be79\")}", requestID = "1.c3438849fd14d720.61.1.7.1"
2. INSERT failed: executing INSERT with THEN RETURN: spanner: code = "FailedPrecondition", desc = "Foreign key `fk_employees_department` constraint violation on table `employees`. Cannot find referenced key `{String(\"sample_value\")}` in table `departments`.", requestID = "1.c3438849fd14d720.61.1.8.1"
3. INSERT failed: executing INSERT: spanner: code = "FailedPrecondition", desc = "Foreign key `fk_pa_project` constraint violation on table `project_assignments`. Cannot find referenced key `{String(\"sample_value\")}` in table `projects`.", requestID = "1.c3438849fd14d720.61.1.9.1"

### instance-7-chatgpt-4o-latest-3-20251007203405.sql

**Parse Error Rate**: 4.3% (1/23 failed to parse)

**Parse Error Types**:
- `Syntax Error: String Literal Quotes`: 1 occurrences

**Parse Errors**:
- syntax error: instance-7-chatgpt-4o-latest-3-20251007203405.sql:7:29: expected token: (, but: <string>
  Statement: `CREATE TABLE projects (
  project_id STRING(36) NOT NULL DEFAULT (GENERATE_UUID()),
  project_name S..._status CHECK (status IN ('ACTIVE', 'COMPLETED', 'ON_HOLD', 'CANCELLED'))
) PRIMARY KEY (project_id)`

**Execution Error Rate**: 22.7% (5/22 parsed statements failed)

**Execution Error Codes**:
- `NotFound`: 2 occurrences
- `InvalidArgument`: 3 occurrences

**Error Categories**:
- `NotFound`: 2 occurrences
- `Table Not Found (InvalidArgument)`: 3 occurrences

**Execution Errors**:
1. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
2. CREATE failed: executing CREATE statement: rpc error: code = NotFound desc = Table not found: projects
3. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: projects [at 1:13]\nINSERT INTO projects (project_name, start_date, end_date, budget, status) VAL...\n            ^", requestID = "1.c3438849fd14d720.62.4.8.1"
4. INSERT failed: executing INSERT: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:13]\nINSERT INTO project_assignments (emp_id, project_id, role, hours_allocated) V...\n            ^", requestID = "1.c3438849fd14d720.62.4.12.1"
5. SELECT failed: reading SELECT results: spanner: code = "InvalidArgument", desc = "Table not found: project_assignments [at 1:277]\n...employees m ON e.manager_id = m.emp_id LEFT JOIN project_assignments pa ON...\n                                                    ^", requestID = "1.c3438849fd14d720.62.4.13.1"

## Compatibility Insights

### Most Common Parse Issues

1. **Syntax Error: Expected Token** (66.0% of parse errors)
   - 66 occurrences across all files
   - Missing expected tokens in SQL syntax. FIX: Add required keywords, punctuation, or identifiers where expected. After DEFAULT remember to wrap the value in parentheses

1. **Syntax Error: String Literal Quotes** (17.0% of parse errors)
   - 17 occurrences across all files
   - String literals in SQL statements are using incorrect quote types. FIX: Use single quotes for string literals instead of double quotes. Change "ACTIVE" to 'ACTIVE'. Spanner requires single quotes for string constants.

1. **Syntax Error: PRIMARY/FOREIGN KEY Placement** (11.0% of parse errors)
   - 11 occurrences across all files
   - PRIMARY KEY constraints MUST be placed OUTSIDE the column definition parentheses in Spanner. FIX: Move PRIMARY KEY clause to after the closing parenthesis of column definitions. CORRECT: ') PRIMARY KEY (column_name);' WRONG: 'PRIMARY KEY (column_name)' inside column list. This is a critical Spanner-specific syntax requirement.

1. **Syntax Error: General** (6.0% of parse errors)
   - 6 occurrences across all files
   - General SQL syntax errors not matching specific patterns. FIX: Review statement structure, check for typos and syntax compliance with Spanner SQL

### Most Common Execution Issues

1. **Table Not Found (InvalidArgument)** (37.4% of execution errors)
   - 95 occurrences across all files
   - Table references that result in InvalidArgument rather than NotFound. FIX: There is likely a error creating the referenced table, so ignore this error

1. **NotFound** (35.4% of execution errors)
   - 90 occurrences across all files
   - Referenced objects (tables, columns, etc.) not found. FIX: There is likely a error creating the referenced table, so ignore this error

1. **FailedPrecondition** (11.8% of execution errors)
   - 30 occurrences across all files
   - Constraint violations or prerequisites not met. FIX: Ensure data meets NOT NULL, foreign key, and other constraints

1. **InvalidArgument: Other** (3.5% of execution errors)
   - 9 occurrences across all files
   - InvalidArgument errors not matching specific patterns. FIX: Review error message details for specific syntax issues

1. **Syntax Error: General** (2.4% of execution errors)
   - 6 occurrences across all files
   - General SQL syntax errors not matching specific patterns. FIX: Check statement structure against Spanner SQL reference

