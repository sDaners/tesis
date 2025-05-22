package main

import (
	"log"
)

func main() {
	db, err := ConnectDB()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	if err := CreateTables(db); err != nil {
		log.Fatal("Error creating tables: ", err)
	}

	_, _, _, err = InsertSampleData(db)
	if err != nil {
		log.Fatal("Error inserting sample data: ", err)
	}

	results, err := QueryEmployeeDetails(db)
	if err != nil {
		log.Fatal("Error querying data: ", err)
	}

	for _, employee := range results {
		log.Printf("Employee: %s %s, Department: %s, Project: %s", employee.FirstName, employee.LastName, employee.DeptName, employee.ProjectName.String)
	}
}
