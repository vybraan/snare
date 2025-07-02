package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github/com/vybraan/snare/helpers"
)

// AuthRequired is a middleware that checks if the user has a valid session.
// It should be used on routes that require authentication.
// If no valid session exists, it aborts the request with 401 Unauthorized.
func AuthRequired(c *gin.Context) {
	// Get the session from the request context
	session := sessions.Default(c)
	// c.Abort()
	// Try to get the user from the session
	if user := session.Get(helpers.UserKey); user == nil {
		// No user in session, abort the request
		//redirect to login
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// User is authenticated, continue to the next handler
	c.Next()
}
