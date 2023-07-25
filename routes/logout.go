// routes/logout.go
package routes

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	// Import other necessary packages
)

func handleLogout(c *gin.Context) {
	session := sessions.Default(c)

	// Clear the session data
	session.Clear()
	session.Save()

	// Redirect to the login page after logout
	c.Redirect(http.StatusSeeOther, "/login")
}
