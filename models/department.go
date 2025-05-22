package models

import "time"

type Department struct {
	DeptID    int       `db:"dept_id"`
	DeptName  string    `db:"dept_name"`
	Location  string    `db:"location"`
	CreatedAt time.Time `db:"created_at"`
}
