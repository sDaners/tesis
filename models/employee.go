package models

import (
	"database/sql"
	"time"
)

type Employee struct {
	EmpID       int       `db:"emp_id"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Email       string    `db:"email"`
	HireDate    time.Time `db:"hire_date"`
	Salary      float64   `db:"salary"`
	DeptID      int       `db:"dept_id"`
	ManagerID   *int      `db:"manager_id"`
	PhoneNumber string    `db:"phone_number"`
}

type EmployeeDetails struct {
	EmpID            int            `db:"emp_id"`
	FirstName        string         `db:"first_name"`
	LastName         string         `db:"last_name"`
	Email            string         `db:"email"`
	DeptName         string         `db:"dept_name"`
	ManagerFirstName sql.NullString `db:"manager_first_name"`
	ManagerLastName  sql.NullString `db:"manager_last_name"`
	ProjectName      sql.NullString `db:"project_name"`
}
