package repo

import "sql-parser/models"

// Database defines the interface for database operations.
type Database interface {
	CreateTables() error
	InsertSampleData() (int, int, int, error)
	QueryEmployeeDetails() ([]models.EmployeeDetails, error)
	CleanupDB() error
}
