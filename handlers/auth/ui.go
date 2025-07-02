package auth

import (
	"github/com/vybraan/snare/helpers"
	"github/com/vybraan/snare/ui"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterPage(c *gin.Context) {

	template := ui.Register()

	helpers.Render(c, http.StatusOK, template)
}
func LoginPage(c *gin.Context) {

	template := ui.Login()

	helpers.Render(c, http.StatusOK, template)
}

func LogoutPage(c *gin.Context) {
	template := ui.Logout()

	helpers.Render(c, http.StatusOK, template)
}
