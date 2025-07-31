#!/bin/bash

# Script to generate and view Allure reports for SQL testing

echo "=== SQL Test Allure Report Generator ==="

# Clean up any existing results
echo "Cleaning up previous results..."
rm -rf allure-results*
rm -rf allure-report

# Run the main SQL compatibility tests
echo "Running SQL compatibility tests..."
cd tests/spanner
go test -run TestGeneratedSQLFiles
cd ../..

# Run atomic SQL tests by statement type
echo "Running atomic SQL tests..."
cd tests/spanner
go test -run TestAtomicCREATEStatements
go test -run TestAtomicINSERTStatements
go test -run TestAtomicSELECTStatements
go test -run TestAtomicDROPStatements
cd ../..

echo "=== Test Results Generated ==="
echo "Allure results directories:"
ls -la | grep allure-results

echo ""
echo "=== To View Allure Reports ==="
echo "Generate and serve reports:"
echo "  allure serve allure-results                    # For SQL compatibility tests"
echo "  allure serve allure-results-create             # For CREATE statement tests"  
echo "  allure serve allure-results-insert             # For INSERT statement tests"
echo "  allure serve allure-results-select             # For SELECT statement tests"
echo "  allure serve allure-results-drop               # For DROP statement tests"
echo ""
echo "Or generate static reports:"
echo "  allure generate allure-results -o allure-report --clean"
echo "  open allure-report/index.html" 