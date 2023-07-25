// routes/success.go
package routes

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	// Import other necessary packages
)

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
