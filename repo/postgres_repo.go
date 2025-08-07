package repo

import (
	"context"
	"database/sql"
	"time"

	"sql-parser/models"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{DB: db}
}

func (r *PostgresRepo) CreateTables() error {
	_, err := r.DB.Exec(PostgresDDL)
	return err
}

func (r *PostgresRepo) InsertSampleData() (int, int, int, error) {
	var deptID int
	err := r.DB.QueryRow(InsertDepartmentSQL, "Engineering", "New York").Scan(&deptID)
	if err != nil {
		return 0, 0, 0, err
	}

	var empID int
	err = r.DB.QueryRow(InsertEmployeeSQL,
		"John", "Doe", "john.doe@example.com", time.Now(), 75000.00, deptID).Scan(&empID)
	if err != nil {
		return 0, 0, 0, err
	}

	var projectID int
	err = r.DB.QueryRow(InsertProjectSQL,
		"Database Migration", time.Now(), time.Now().AddDate(0, 3, 0), 50000.00, "ACTIVE").Scan(&projectID)
	if err != nil {
		return 0, 0, 0, err
	}

	_, err = r.DB.Exec(InsertProjectAssignmentSQL,
		empID, projectID, "Lead Developer", 40)
	if err != nil {
		return 0, 0, 0, err
	}

	return deptID, empID, projectID, nil
}

func (r *PostgresRepo) QueryEmployeeDetails() ([]models.EmployeeDetails, error) {
	query := QueryEmployeeDetailsSQL

	var details []models.EmployeeDetails
	err := sqlscan.Select(context.Background(), r.DB, &details, query)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (r *PostgresRepo) CleanupDB() error {
	_, err := r.DB.Exec(CleanupDBSQL)
	return err
}
