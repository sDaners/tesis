package repo

import (
	"context"
	"database/sql"
	"time"

	"postgres-example/models"
	"postgres-example/repo/valid_repo"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type SpannerRepo struct {
	DB *sql.DB
}

func NewSpannerRepo(db *sql.DB) *SpannerRepo {
	return &SpannerRepo{DB: db}
}

func (r *SpannerRepo) CreateTables() error {
	for _, ddl := range valid_repo.SpannerDDL {
		_, err := r.DB.Exec(ddl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SpannerRepo) InsertSampleData() (int, int, int, error) {
	tx, err := r.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return 0, 0, 0, err
	}
	defer tx.Rollback()

	var deptID int64
	err = sqlscan.Get(context.Background(), tx, &deptID, valid_repo.SpannerInsertDepartmentSQL,
		sql.Named("dept_name", "Engineering"),
		sql.Named("location", "New York"))
	if err != nil {
		return 0, 0, 0, err
	}

	var empID int64
	err = sqlscan.Get(context.Background(), tx, &empID, valid_repo.SpannerInsertEmployeeSQL,
		sql.Named("first_name", "John"),
		sql.Named("last_name", "Doe"),
		sql.Named("email", "john.doe@example.com"),
		sql.Named("hire_date", time.Now()),
		sql.Named("salary", 75000.00),
		sql.Named("dept_id", deptID))
	if err != nil {
		return 0, 0, 0, err
	}

	var projectID int64
	err = sqlscan.Get(context.Background(), tx, &projectID, valid_repo.SpannerInsertProjectSQL,
		sql.Named("project_name", "Database Migration"),
		sql.Named("start_date", time.Now()),
		sql.Named("end_date", time.Now().AddDate(0, 3, 0)),
		sql.Named("budget", 50000.00),
		sql.Named("status", "ACTIVE"))
	if err != nil {
		return 0, 0, 0, err
	}

	_, err = tx.Exec(valid_repo.SpannerInsertProjectAssignmentSQL,
		sql.Named("emp_id", empID),
		sql.Named("project_id", projectID),
		sql.Named("role", "Lead Developer"),
		sql.Named("hours", 40))
	if err != nil {
		return 0, 0, 0, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, 0, 0, err
	}

	return int(deptID), int(empID), int(projectID), nil
}

func (r *SpannerRepo) QueryEmployeeDetails() ([]models.EmployeeDetails, error) {
	var details []models.EmployeeDetails
	err := sqlscan.Select(context.Background(), r.DB, &details, valid_repo.SpannerQueryEmployeeDetailsSQL)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (r *SpannerRepo) CleanupDB() error {
	// Drop objects in reverse order of dependencies
	for _, stmt := range valid_repo.SpannerCleanupStatements {
		_, err := r.DB.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
