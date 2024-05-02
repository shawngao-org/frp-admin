package api

import (
	"frp-admin/service"
	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      Login
// @Tags         User
// @Accept       multipart/form-data
// @Produce      application/json
// @Success      200  {object}  string
// @Param        email formData string true "Email"
// @Param        password formData string true "Password(RSA Encrypted)"
// @Router       /api/v1/user/login [post]
func Login(ctx *gin.Context) {
	service.Login(ctx)
}

// Register godoc
// @Summary      Register
// @Tags         User
// @Accept       multipart/form-data
// @Produce      application/json
// @Success      200  {object}  string
// @Param        name formData string true "Name"
// @Param        email formData string true "Email"
// @Param        password formData string true "Password(RSA Encrypted)"
// @Router       /api/v1/user/register [post]
func Register(ctx *gin.Context) {
	service.RegisterUser(ctx)
}

// SendTestMail godoc
// @Summary      Send Test Mail
// @Tags         User
// @Accept       multipart/form-data
// @Produce      application/json
// @Success      200  {object}  string
// @Param        email formData string true "Email"
// @Router       /api/v1/test-mail [post]
func SendTestMail(ctx *gin.Context) {
	service.SendTestMail(ctx)
}
