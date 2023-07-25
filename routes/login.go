// routes/login.go
package routes

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	// Import other necessary packages

	"github.com/dj-godev/my-forumproject/models" // Import the models package to access GetUserByEmail function
)

func showLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func handleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Check if the user with the given email exists
	user, err := models.GetUserByEmail(email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Error retrieving user data.",
		})
		return
	}

	// Verify the password
	if user == nil || user.Password != models.GetMD5Hash(password) {
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
