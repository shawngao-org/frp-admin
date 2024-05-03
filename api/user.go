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

// SendForgetPasswordMail godoc
// @Summary      Send Forget Password Mail
// @description  Send forget password mail, but front-end must have "http://xxx.xxx.xxx/reset-password/:code" router.
// @Tags         User
// @Accept       multipart/form-data
// @Produce      application/json
// @Success      200  {object}  string
// @Param        email formData string true "Email"
// @Router       /api/v1/user/forget-password [post]
func SendForgetPasswordMail(ctx *gin.Context) {
	service.SendForgetPasswordMail(ctx)
}

// ResetPassword godoc
// @Summary      Verify tmp code and reset password
// @Tags         User
// @Accept       multipart/form-data
// @Produce      application/json
// @Success      200  {object}  string
// @Param        email formData string true "Email"
// @Param        password formData string true "New Password (RSA Encrypted)"
// @Param        code formData string true "Verify code (Temp code)"
// @Router       /api/v1/user/reset-password [post]
func ResetPassword(ctx *gin.Context) {
	service.ResetPassword(ctx)
}
