package helpers

import (
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	UserKey = "user"    // key used to store the username in the session
	Secret  = "secreta" // random and secure key used to encrypt the session cookie
)

func GetUserIDFromSession(c *gin.Context) (int64, error) {
	session := sessions.Default(c)
	val := session.Get(UserKey)
	if val == nil {
		return 0, errors.New("user not logged in")
	}

	userID, ok := val.(int64)
	if !ok {
		return 0, errors.New("invalid user ID in session")
	}

	return userID, nil
}
