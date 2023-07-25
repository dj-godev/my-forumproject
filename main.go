package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dj-godev/my-forumproject/models"
	"github.com/dj-godev/my-forumproject/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database connection
	dbConn := "root:@tcp(localhost:3306)/go-webdb" // Replace with your actual database connection string
	err := models.InitDB(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	// Set Gin to production mode for better performance
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	// Initialize Gin router
	r := gin.Default()

	// Initialize the session middleware with a secret key
	store := cookie.NewStore([]byte("forumIAS@123#"))
	r.Use(sessions.Sessions("mysession", store))

	// Initialize the routes
	routes.InitRoutes(r)

	// Serve static files (CSS, JS, etc.)
	r.Static("/static", "./static")

	// Load HTML templates from the "templates" folder
	r.LoadHTMLGlob("templates/*.html")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/login")
	})

	fmt.Println("Server started at http://localhost:8082")
	err = r.Run(":8082")
	if err != nil {
		log.Fatal(err)
	}
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
