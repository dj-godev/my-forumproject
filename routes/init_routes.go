// routes/init_routes.go
package routes

import (
	"github.com/dj-godev/my-forumproject/middlewares"
	"github.com/gin-gonic/gin"
)

// InitRoutes initializes all the routes for the application
func InitRoutes(r *gin.Engine) {
	// Routes for login
	r.GET("/login", middlewares.EnsureNotLoggedIn(), showLoginPage)
	r.POST("/login", handleLogin)

	// Routes for signup
	r.GET("/signup", showSignupPage)
	r.POST("/signup", handleSignup)

	// Routes for logout
	r.GET("/logout", handleLogout)

	// Routes for success page
	r.GET("/login-success", middlewares.EnsureLoggedIn(), showSuccessPage)

	// ... Other routes ...
}
