package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// User represents a user record in the database
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func main() {
	// Connect to the MySQL database
	dbConn, err := sql.Open("mysql", "root:@tcp(localhost:3306)/go-webdb")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	db = dbConn

	// Set Gin to production mode for better performance
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	// Initialize Gin router
	r := gin.Default()

	// Initialize the session middleware with a secret key
	store := cookie.NewStore([]byte("forumIAS@123#"))
	r.Use(sessions.Sessions("mysession", store))

	// Serve static files (CSS, JS, etc.)
	r.Static("/static", "./static")

	// Load HTML templates from the "templates" folder
	r.LoadHTMLGlob("templates/*.html")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/login")
	})

	//r.GET("/login", showLoginPage)
	r.GET("/login", ensureNotLoggedIn(), showLoginPage)
	r.POST("/login", handleLogin)

	r.GET("/signup", showSignupPage)
	r.POST("/signup", handleSignup)

	//r.GET("/login-success", showSuccessPage)
	r.GET("/login-success", ensureLoggedIn(), showSuccessPage)

	r.GET("/logout", handleLogout) // Logout route

	fmt.Println("Server started at http://localhost:8082")
	err = r.Run(":8082")
	if err != nil {
		log.Fatal(err)
	}
}

// Middleware to ensure the user is not logged in
func ensureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userId")
		if userID != nil {
			// User is already logged in, redirect to login-success page
			c.Redirect(http.StatusSeeOther, "/login-success")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Middleware to ensure the user is logged in
func ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userId")
		if userID == nil {
			// User is not logged in, redirect to login page
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func handleLogout(c *gin.Context) {
	session := sessions.Default(c)

	// Clear the session data
	session.Clear()
	session.Save()

	// Redirect to the login page after logout
	c.Redirect(http.StatusSeeOther, "/login")
}

func showLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func handleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Check if the user with the given email exists
	user, err := GetUserByEmail(db, email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Error retrieving user data.",
		})
		return
	}

	// Verify the password
	if user == nil || user.Password != GetMD5Hash(password) {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error": "Invalid email or password.",
		})
		return
	}

	// User is authenticated, handle login (e.g., set session cookie)
	// For simplicity, we'll just redirect to a success page
	//c.Redirect(http.StatusSeeOther, "/login-success")

	// Render the "login-success" page with user data
	//showSuccessPage(c, user.Name)
	//return

	// User is authenticated, handle login (e.g., set session variable)
	session := sessions.Default(c)
	session.Set("userId", user.ID)
	session.Set("fullName", user.Name)
	session.Save()

	// Redirect to the "login-success" page
	c.Redirect(http.StatusSeeOther, "/login-success")
}

func showSignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

func handleSignup(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Insert the new user into the database
	err := InsertUserIntoDB(db, name, email, password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Error creating user.",
		})
		return
	}

	// User is signed up, handle signup success (e.g., set session cookie)
	// For simplicity, we'll just redirect to a success page
	c.Redirect(http.StatusSeeOther, "/signup-success")
}

func showSuccessPage(c *gin.Context) {
	session := sessions.Default(c)

	// Get the user name from the session
	userName := session.Get("fullName")
	userId := session.Get("userId")

	c.HTML(http.StatusOK, "success.html", gin.H{
		"FullName": userName,
		"UserId":   userId,
	})
}

// func showSuccessPage(c *gin.Context, name string) {
// 	c.HTML(http.StatusOK, "success.html", gin.H{
// 		"Name": name,
// 	})
// }

// func showSuccessPage(c *gin.Context) {
// 	// Retrieve the user's name from the database using the email provided during login
// 	email := c.PostForm("email")
// 	user, err := GetUserByEmail(db, email)

// 	if err != nil || user == nil {
// 		//log.Fatal(user)
// 		fmt.Println("emialll", email)
// 		errorMessage := "Error retrieving user data."
// 		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
// 			"Error": errorMessage,
// 		})
// 		return
// 	}

// 	// Get the user's name from the user object
// 	name := user.Name

// 	c.HTML(http.StatusOK, "success.html", gin.H{
// 		"Name": name,
// 	})
// }

// GetUserByEmail fetches a user from the database based on their email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	// Query the database for the user with the given email
	row := db.QueryRow("SELECT id, fullName, email, password FROM users WHERE email = ?", email)

	// Create a new User object to store the retrieved data
	user := &User{}

	// Scan the result row into the User object
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// User with the given email not found
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// InsertUserIntoDB inserts a new user into the database
func InsertUserIntoDB(db *sql.DB, name, email, password string) error {
	// Prepare the INSERT statement
	stmt, err := db.Prepare("INSERT INTO users (fullName, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the INSERT statement with the user data
	_, err = stmt.Exec(name, email, password)
	if err != nil {
		return err
	}

	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
