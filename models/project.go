package models

import "time"

type Project struct {
	ProjectID   int       `db:"project_id"`
	ProjectName string    `db:"project_name"`
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
	Budget      float64   `db:"budget"`
	Status      string    `db:"status"`
}

type ProjectAssignment struct {
	EmpID          int    `db:"emp_id"`
	ProjectID      int    `db:"project_id"`
	Role           string `db:"role"`
	HoursAllocated int    `db:"hours_allocated"`
}
