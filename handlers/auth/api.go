package auth

import (
	"database/sql"
	"errors"
	"github/com/vybraan/snare/components/toast"
	"github/com/vybraan/snare/helpers"
	"github/com/vybraan/snare/models"
	"github/com/vybraan/snare/ui/partials"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

func LoginAPI(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		var credentials models.User
		if err := c.ShouldBind(&credentials); err != nil {

			template := toast.Toast(toast.Props{
				Description: "Invalid input",
				Title:       "Error",
				Variant:     toast.VariantError,
			})
			helpers.Render(c, http.StatusBadRequest, template)

			// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		username := strings.TrimSpace(credentials.Username)
		password := strings.TrimSpace(credentials.Password)
		if username == "" || password == "" {

			template := toast.Toast(toast.Props{
				Description: "Username and password cannot be empty",
				Title:       "Error",
				Variant:     toast.VariantError,
			})

			helpers.Render(c, http.StatusBadRequest, template)
			// c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password cannot be empty"})
			return
		}

		// Fetch user record
		var user models.User
		query := "SELECT id, username, password FROM users WHERE username = ?"
		err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {

				template := toast.Toast(toast.Props{
					Description: "Invalid username or password",
					Title:       "Error",
					Variant:     toast.VariantError,
				})

				helpers.Render(c, http.StatusUnauthorized, template)

				// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			} else {
				template := toast.Toast(toast.Props{
					Description: "Internal server error",
					Title:       "Error",
					Variant:     toast.VariantError,
				})
				helpers.Render(c, http.StatusInternalServerError, template)
				// c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}

		// Compare hashed passwords
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
			template := toast.Toast(toast.Props{
				Description: "Invalid username or password",
				Title:       "Error",
				Variant:     toast.VariantError,
			})
			helpers.Render(c, http.StatusUnauthorized, template)

			// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		session.Clear()
		session.Options(sessions.Options{
			MaxAge:   3600,
			Path:     "/",
			HttpOnly: true,
		})

		// Store user ID in session
		session.Set(helpers.UserKey, int64(user.ID))
		if err := session.Save(); err != nil {
			template := toast.Toast(toast.Props{
				Description: "Failed to save session",
				Title:       "Error",
				Variant:     toast.VariantError,
			})
			helpers.Render(c, http.StatusInternalServerError, template)

			// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		// c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated"})
		// Render a success template or redirect as needed
		template := partials.LoginSuccess()
		helpers.Render(c, http.StatusOK, template)

	}
}

func RegisterAPI(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBind(&user); err != nil {
			template := toast.Toast(toast.Props{
				Description: "Invalid input",
				Title:       "Error",
				Variant:     toast.VariantError,
			})
			helpers.Render(c, http.StatusBadRequest, template)
			return
		}

		if user.Username == "" || user.Password == "" {
			template := toast.Toast(toast.Props{
				Description: "Username and password cannot be empty",
				Title:       "Error",
				Variant:     toast.VariantError,
			})
			helpers.Render(c, http.StatusBadRequest, template)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			template := toast.Toast(toast.Props{
				Description: "Failed to hash password",
				Title:       "Error",
				Variant:     toast.VariantError,
			})
			helpers.Render(c, http.StatusInternalServerError, template)
			return
		}
		user.Password = string(hashedPassword)

		query := "INSERT INTO users (username, password) VALUES (?, ?)"
		result, err := db.Exec(query, user.Username, user.Password)
		if err != nil {
			template := toast.Toast(toast.Props{
				Description: "Failed to register user",
				Title:       "Error",
				Variant:     toast.VariantError,
			})
			helpers.Render(c, http.StatusInternalServerError, template)
			return
		}

		_, err = result.LastInsertId()
		if err != nil {
			template := toast.Toast(toast.Props{
				Description: "Failed to retrieve user ID",
				Title:       "Error",
				Variant:     toast.VariantError,
			})
			helpers.Render(c, http.StatusInternalServerError, template)
			return
		}

		session := sessions.Default(c)
		session.Clear()
		session.Options(sessions.Options{
			MaxAge:   3600,
			Path:     "/",
			HttpOnly: true,
		})

		// Store user ID in session
		session.Set(helpers.UserKey, int64(user.ID))
		if err := session.Save(); err != nil {
			template := toast.Toast(toast.Props{
				Description: "Failed to save session",
				Title:       "Error",
				Variant:     toast.VariantError,
			})
			helpers.Render(c, http.StatusInternalServerError, template)

			return
		}

		// Render a success template or redirect as needed
		template := partials.RegisterSuccess()
		helpers.Render(c, http.StatusOK, template)
	}
}

func LogoutAPI(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(helpers.UserKey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(helpers.UserKey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	template := partials.LogoutSuccess()
	helpers.Render(c, http.StatusOK, template)

	// c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func StatusAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}
