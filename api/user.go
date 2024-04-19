package api

import (
	"frp-admin/config"
	"frp-admin/redis"
	"frp-admin/service"
	"frp-admin/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Login godoc
// @Summary      Login
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  string
// @Param        email query string true "Email"
// @Param        password query string true "Password(RSA Encrypted)"
// @Router       /api/v1/login [post]
func Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	passwd := ctx.PostForm("password")
	passwd = util.Decrypted(passwd)
	user, err := service.GetUserByEmailAndPasswd(email, passwd)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message":   http.StatusText(http.StatusForbidden),
			"exception": err.Error(),
		})
		return
	}
	token, err := util.GenerateToken(user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message":   http.StatusText(http.StatusInternalServerError),
			"exception": err.Error(),
		})
		return
	}
	status := redis.Client.Set(ctx, user.Id, token, time.Duration(config.Conf.Security.Jwt.Timeout))
	if status.Err() != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message":   http.StatusText(http.StatusInternalServerError),
			"exception": status.Err(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": http.StatusText(http.StatusOK),
		"token":   token,
	})
}
