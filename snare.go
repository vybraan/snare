package snare

import (
	"database/sql"
	"github.com/vybraan/snare/handlers/auth"
	"github.com/vybraan/snare/helpers"
	"github.com/vybraan/snare/middlewares"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type AuthModule struct {
	DB *sql.DB
}

func New(db *sql.DB) *AuthModule {
	return &AuthModule{DB: db}
}

func (a *AuthModule) RegisterRoutes(r *gin.Engine) {
	store := cookie.NewStore([]byte(helpers.Secret))
	store.Options(sessions.Options{
		MaxAge:   3600,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	r.Use(sessions.Sessions("mysession", store))

	r.Static("/snare", "./assets")

	r.GET("/auth/login", auth.LoginPage)
	r.GET("/auth/logout", middlewares.AuthRequired, auth.LogoutPage)
	r.GET("/auth/register", auth.RegisterPage)

	api := r.Group("/api")
	api.POST("/auth/login", auth.LoginAPI(a.DB))
	api.POST("/auth/register", auth.RegisterAPI(a.DB))
	api.GET("/auth/logout", auth.LogoutAPI)

	api.Use(middlewares.AuthRequired)
	api.GET("/auth/state", auth.StatusAPI)
}
