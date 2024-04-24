package api

import (
	"frp-admin/config"
	"frp-admin/db"
	"frp-admin/entity"
	"frp-admin/redis"
	"frp-admin/service"
	"frp-admin/util"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// Login godoc
// @Summary      Login
// @Tags         User
// @Accept       multipart/form-data
// @Produce      application/json
// @Success      200  {object}  string
// @Param        email formData string true "Email"
// @Param        password formData string true "Password(RSA Encrypted)"
// @Router       /api/v1/login [post]
func Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	passwd := ctx.PostForm("password")
	passwd = util.Decrypted(passwd)
	validMsg := util.VerificationEmailAndPassword(email, passwd)
	if !util.IsPassValid(validMsg) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": validMsg,
		})
		return
	}
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

// Register godoc
// @Summary      Register
// @Tags         User
// @Accept       multipart/form-data
// @Produce      application/json
// @Success      200  {object}  string
// @Param        name formData string true "Name"
// @Param        email formData string true "Email"
// @Param        password formData string true "Password(RSA Encrypted)"
// @Router       /api/v1/register [post]
func Register(ctx *gin.Context) {
	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	password = util.Decrypted(password)
	validMsg := util.UnifiedVerificationOfBasicUserInfo(name, email, password)
	if !util.IsPassValid(validMsg) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": validMsg,
		})
		return
	}
	password, err := util.PasswordEncrypt(password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	if service.CheckUserExists(name, email) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "User already exists.",
		})
		return
	}
	user := entity.User{
		Model:        gorm.Model{},
		Id:           utils.NewUUID().String(),
		Name:         name,
		Email:        email,
		Password:     password,
		IsValid:      false,
		RegisterTime: time.Now(),
		Ip:           ctx.RemoteIP(),
		Key:          utils.NewUUID().String(),
		GroupId:      "47bbe440-dfcb-435f-b7ef-dba7b54a2135",
	}
	result := db.Db.Create(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": result.Error,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": http.StatusText(http.StatusOK),
	})
}
