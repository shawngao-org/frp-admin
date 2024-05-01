package service

import (
	"errors"
	"frp-admin/common"
	"frp-admin/config"
	"frp-admin/db"
	"frp-admin/entity"
	"frp-admin/redis"
	"frp-admin/util"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	db.Db.First(&user, "email = ?", email)
	if reflect.DeepEqual(user, entity.User{}) {
		return entity.User{}, errors.New("user not found")
	}
	return user, nil
}

func CheckUserExists(name string, email string) bool {
	var user entity.User
	db.Db.First(&user, "name = ? OR email = ?", name, email)
	return !reflect.DeepEqual(user, entity.User{})
}

func GetUserByEmailAndPasswd(email string, password string) (entity.User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return entity.User{}, errors.New("user not found")
	}
	if strings.ToUpper(config.Conf.Security.Password.Method) == "BCRYPT" {
		e := util.CheckBcrypt(password, user.Password)
		if e == nil {
			user.Password = ""
			return user, nil
		} else {
			return entity.User{}, errors.New("invalid password")
		}
	}
	str, err := util.PasswordEncrypt(password)
	if err != nil {
		return entity.User{}, err
	}
	if str == user.Password {
		user.Password = ""
		return user, nil
	}
	return entity.User{}, errors.New("invalid password")
}

func RegisterUser(ctx *gin.Context) {

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
	if CheckUserExists(name, email) {
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
	user, err := GetUserByEmailAndPasswd(email, passwd)
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
	status := redis.Client.Set(common.Context, user.Id, token, time.Duration(config.Conf.Security.Jwt.Timeout)*time.Second)
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

func SendTestMail(ctx *gin.Context) {
	email := ctx.PostForm("email")
	util.SendTestMail(email)
}
