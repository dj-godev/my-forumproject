// routes/signup.go
package routes

import (
	"net/http"

	models "github.com/dj-godev/my-forumproject/models" // Import the models package to access InsertUserIntoDB function
	"github.com/gin-gonic/gin"
)

func showSignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

func handleSignup(c *gin.Context) {
	// Get the form data from the signup form
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Create a new user with the form data
	newUser := &models.User{
		Name:     name,
		Email:    email,
		Password: password, // Remember to use a secure password hashing method
	}

	// Insert the new user into the database
	err := models.InsertUserIntoDB(newUser)
	if err != nil {
		// Handle the error, show an error message, or redirect to an error page
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"ErrorMessage": "Error creating user account.",
		})
		return
	}

	// Registration successful, redirect the user to the login page or a success page
	c.Redirect(http.StatusSeeOther, "/login")
}
