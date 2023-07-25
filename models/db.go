// models/db.go
package models

import (
	"database/sql"
	"fmt"

	// Import the necessary database driver package (e.g., "github.com/go-sql-driver/mysql")
	_ "github.com/go-sql-driver/mysql"
)

// Define a global variable to hold the database connection.
var db *sql.DB

// InitDB initializes the database connection.
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Check if the database connection is successful.
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
