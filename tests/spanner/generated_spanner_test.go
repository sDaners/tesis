package spanner_test

import (
	"database/sql"
	"testing"

	"postgres-example/repo"
	"postgres-example/tools"
)

type SpannerDBTeardown struct {
	db        *sql.DB
	repo      repo.Database
	t         *testing.T
	terminate func()
}

func setupSpannerDB(t *testing.T) *SpannerDBTeardown {
	db, terminate, err := tools.GetDB(true)
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	r := repo.NewSpannerRepo(db)
	if err := r.CleanupDB(); err != nil {
		t.Fatalf("Failed to cleanup DB: %v", err)
	}
	if err := r.CreateTables(); err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}
	return &SpannerDBTeardown{db: db, repo: r, t: t, terminate: terminate}
}

func (d *SpannerDBTeardown) Close() {
	if err := d.repo.CleanupDB(); err != nil {
		d.t.Errorf("Failed to cleanup DB: %v", err)
	}
	d.db.Close()
	d.terminate()
}

func TestSpannerDatabaseOperations(t *testing.T) {
	dbT := setupSpannerDB(t)
	defer dbT.Close()

	// Insert sample data
	deptID, empID, projectID, err := dbT.repo.InsertSampleData()
	if err != nil {
		t.Fatalf("InsertSampleData failed: %v", err)
	}
	if deptID == 0 || empID == 0 || projectID == 0 {
		t.Errorf("Expected non-zero IDs, got deptID=%d, empID=%d, projectID=%d", deptID, empID, projectID)
	}

	// Query and check results
	details, err := dbT.repo.QueryEmployeeDetails()
	if err != nil {
		t.Fatalf("QueryEmployeeDetails failed: %v", err)
	}
	if len(details) == 0 {
		t.Error("Expected at least one employee detail result")
	}

	// Check if we found our test employee
	found := false
	for _, detail := range details {
		if detail.FirstName == "John" &&
			detail.LastName == "Doe" &&
			detail.DeptName == "Engineering" &&
			detail.ProjectName.String == "Database Migration" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find inserted employee and project in results")
	}
}

func TestSpannerEmployeeDetails(t *testing.T) {
	dbT := setupSpannerDB(t)
	defer dbT.Close()

	// Insert test data
	_, _, _, err := dbT.repo.InsertSampleData()
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Query employee details
	details, err := dbT.repo.QueryEmployeeDetails()
	if err != nil {
		t.Fatalf("Failed to query employee details: %v", err)
	}

	// Verify the structure of the results
	if len(details) == 0 {
		t.Error("Expected at least one employee detail")
	}

	for _, detail := range details {
		// Check required fields
		if detail.FirstName == "" {
			t.Error("Expected non-empty FirstName")
		}
		if detail.LastName == "" {
			t.Error("Expected non-empty LastName")
		}
		if detail.Email == "" {
			t.Error("Expected non-empty Email")
		}
		if detail.DeptName == "" {
			t.Error("Expected non-empty DeptName")
		}

		// Check nullable fields
		if detail.ManagerFirstName.Valid {
			t.Logf("Manager First Name: %s", detail.ManagerFirstName.String)
		}
		if detail.ManagerLastName.Valid {
			t.Logf("Manager Last Name: %s", detail.ManagerLastName.String)
		}
		if detail.ProjectName.Valid {
			t.Logf("Project Name: %s", detail.ProjectName.String)
		}
	}
}
