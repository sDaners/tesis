package main

import (
	"database/sql"
	"testing"
)

func setupDB(t *testing.T) *DBTeardown {
	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	if err := CleanupDB(db); err != nil {
		t.Fatalf("Failed to cleanup DB: %v", err)
	}
	if err := CreateTables(db); err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}
	return &DBTeardown{db: db, t: t}
}

type DBTeardown struct {
	db *sql.DB
	t  *testing.T
}

func (d *DBTeardown) Close() {
	if err := CleanupDB(d.db); err != nil {
		d.t.Errorf("Failed to cleanup DB: %v", err)
	}
	d.db.Close()
}

func TestDatabaseOperations(t *testing.T) {
	dbT := setupDB(t)
	defer dbT.Close()
	db := dbT.db

	// Insert sample data
	deptID, empID, projectID, err := InsertSampleData(db)
	if err != nil {
		t.Fatalf("InsertSampleData failed: %v", err)
	}
	if deptID == 0 || empID == 0 || projectID == 0 {
		t.Errorf("Expected non-zero IDs, got deptID=%d, empID=%d, projectID=%d", deptID, empID, projectID)
	}

	// Query and check results
	details, err := QueryEmployeeDetails(db)
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

func TestEmployeeDetails(t *testing.T) {
	dbT := setupDB(t)
	defer dbT.Close()
	db := dbT.db

	// Insert test data
	_, _, _, err := InsertSampleData(db)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Query employee details
	details, err := QueryEmployeeDetails(db)
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
