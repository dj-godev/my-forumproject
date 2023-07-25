// models/user.go
package models

// Import necessary packages for database handling
import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

// Define the User struct
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

// Define functions for user retrieval, insertion, etc.

// GetUserByEmail fetches a user from the database by email
func GetUserByEmail(email string) (*User, error) {
	// Prepare the SQL query to fetch the user by email
	query := "SELECT id, fullName, email, password FROM users WHERE email = ?"

	// Execute the query and scan the result into a User struct
	row := db.QueryRow(query, email)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows are returned, the user with the given email doesn't exist
			return nil, nil
		}
		// If there was an error other than no rows, return it
		return nil, err
	}

	// Return the user
	return user, nil
}

// InsertUserIntoDB inserts a new user into the database
func InsertUserIntoDB(user *User) error {
	// Prepare the SQL query to insert a new user into the "users" table
	query := "INSERT INTO users (fullName, email, password) VALUES (?, ?, ?)"

	// Execute the query with the user data
	_, err := db.Exec(query, user.Name, user.Email, user.Password)

	if err != nil {
		// If there was an error inserting the user, return it
		return err
	}

	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
