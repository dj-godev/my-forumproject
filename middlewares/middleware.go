// middlewares/middleware.go
package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Middleware to ensure the user is not logged in
func EnsureNotLoggedIn() gin.HandlerFunc {
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
func EnsureLoggedIn() gin.HandlerFunc {
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
