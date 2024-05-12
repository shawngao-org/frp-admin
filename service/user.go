package service

import (
	"errors"
	"fmt"
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

func CheckEmail(email string) (entity.User, error) {
	validMsg := util.ValidateEmail(email)
	if !util.IsPassValid(validMsg) {
		return entity.User{}, errors.New(validMsg)
	}
	user, err := GetUserByEmail(email)
	if err != nil {
		return entity.User{}, errors.New(err.Error())
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
	password, err := util.Decrypted(password)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	validMsg := util.UnifiedVerificationOfBasicUserInfo(name, email, password)
	if !util.IsPassValid(validMsg) {
		ErrHandle(ctx, validMsg, http.StatusBadRequest)
		return
	}
	password, err = util.PasswordEncrypt(password)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	if CheckUserExists(name, email) {
		ErrHandle(ctx, "User already exists.", http.StatusInternalServerError)
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
		GroupId:      config.Conf.Data.GroupId,
	}
	result := db.Db.Create(&user)
	if result.Error != nil {
		ErrHandle(ctx, result.Error, http.StatusInternalServerError)
		return
	}
	OkHandle(ctx)
}

func Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	passwd := ctx.PostForm("password")
	passwd, err := util.Decrypted(passwd)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	validMsg := util.VerificationEmailAndPassword(email, passwd)
	if !util.IsPassValid(validMsg) {
		ErrHandle(ctx, validMsg, http.StatusBadRequest)
		return
	}
	user, err := GetUserByEmailAndPasswd(email, passwd)
	if err != nil {
		ErrHandle(ctx, err, http.StatusForbidden)
		return
	}
	token, err := util.GenerateToken(user.Id)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	status := redis.Client.Set(common.Context, user.Id, token, time.Duration(config.Conf.Security.Jwt.Timeout)*time.Second)
	if status.Err() != nil {
		ErrHandle(ctx, status.Err(), http.StatusInternalServerError)
		return
	}
	OkHandleByPayload(ctx, gin.H{"token": token})
}

func SendRegisterVerifyMail(ctx *gin.Context) {
	SendVerificationMail(ctx, "register")
}

func SendForgetPasswordMail(ctx *gin.Context) {
	SendVerificationMail(ctx, "resetPassword")
}

func SendVerificationMail(ctx *gin.Context, mailType string) {
	email := ctx.PostForm("email")
	user, err := CheckEmail(email)
	if err != nil {
		ErrHandle(ctx, err, http.StatusBadRequest)
		return
	}
	if user.IsValid {
		ErrHandle(ctx, "Do not duplicate verification.", http.StatusBadRequest)
		return
	}
	var redisKey, title, btnLinkPath, btnText string
	switch mailType {
	case "register":
		redisKey = getRegisterVerifyRedisKey(user.Id)
		title = "Register Verify"
		btnLinkPath = "/verify-register/"
		btnText = "Click to verify"
	case "resetPassword":
		redisKey = getResetPasswordRedisKey(user.Id)
		title = "Reset Password"
		btnLinkPath = "/reset-password/"
		btnText = "Click to change your password"
	default:
		ErrHandle(ctx, fmt.Errorf("unknown mail type"), http.StatusBadRequest)
		return
	}
	code, err := GenerateTmpCode2Redis(redisKey)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	content := util.MailContent{
		Title: title,
		Content: fmt.Sprintf("Click the button below to %s your account. "+
			"Note: If you did not initiate the request, please ignore this email so as not to cause unnecessary trouble.",
			strings.ToLower(title)),
		BtnLink: fmt.Sprintf("%s%s%s", config.Conf.Server.FrontEndAddr, btnLinkPath, code),
		BtnText: btnText,
		Author:  "FRP-Admin",
		Note:    util.DefaultFooterNote,
	}
	util.SendDefaultMail(user.Email, content.Title, &content)
	OkHandle(ctx)
}

func ConfirmVerifyRegister(ctx *gin.Context) {
	email := ctx.PostForm("email")
	code := ctx.PostForm("code")
	user, err := CheckEmail(email)
	if err != nil {
		ErrHandle(ctx, err, http.StatusBadRequest)
		return
	}
	redisKey := getRegisterVerifyRedisKey(user.Id)
	_, err = VerifyTmpCodeRedis(code, redisKey)
	if err != nil {
		ErrHandle(ctx, err, http.StatusBadRequest)
		return
	}
	user.IsValid = true
	err = UpdateById(&user)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	OkHandle(ctx)
}

func ResetPassword(ctx *gin.Context) {
	email := ctx.PostForm("email")
	encryptedPassword := ctx.PostForm("password")
	verifyCode := ctx.PostForm("code")
	user, err := CheckEmail(email)
	if err != nil {
		ErrHandle(ctx, err, http.StatusBadRequest)
		return
	}
	redisKey := getResetPasswordRedisKey(user.Id)
	_, err = VerifyTmpCodeRedis(redisKey, verifyCode)
	if err != nil {
		ErrHandle(ctx, err, http.StatusBadRequest)
		return
	}
	decryptedPassword, err := util.Decrypted(encryptedPassword)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	if validMsg := util.ValidatePassword(decryptedPassword); validMsg != "" {
		ErrHandle(ctx, validMsg, http.StatusBadRequest)
		return
	}
	passwordResult, err := util.PasswordEncrypt(decryptedPassword)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	user.Password = passwordResult
	err = UpdateById(&user)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	OkHandle(ctx)
}

func GenerateTmpCode2Redis(redisKey string) (string, error) {
	tmpCode := entity.TmpCode{
		Model:  gorm.Model{},
		Id:     utils.NewUUID().String(),
		IsUsed: false,
	}
	result := db.Db.Create(&tmpCode)
	if result.Error != nil {
		return "", result.Error
	}
	if _, err := redis.Client.Get(common.Context, redisKey).Result(); err == nil {
		return "", errors.New("sending too often, please try again later")
	}
	status := redis.Client.Set(common.Context, redisKey, tmpCode.Id, 5*time.Minute)
	if status.Err() != nil {
		return "", status.Err()
	}
	return tmpCode.Id, nil
}

func VerifyTmpCodeRedis(code string, redisKey string) (bool, error) {
	val, err := redis.Client.Get(common.Context, redisKey).Result()
	if err != nil {
		return false, err
	}
	if val != code {
		return false, errors.New("the verification code is not correct")
	}
	var result entity.TmpCode
	err = GetById(&result, val)
	if err != nil {
		return false, err
	}
	if result.IsUsed {
		return false, errors.New("the verification code is invalid")
	}
	result.IsUsed = true
	err = UpdateById(&result)
	if err != nil {
		return false, err
	}
	redis.Client.Del(common.Context, redisKey)
	return true, nil
}

func getResetPasswordRedisKey(userId string) string {
	return fmt.Sprintf("rp-%s", userId)
}

func getRegisterVerifyRedisKey(userId string) string {
	return fmt.Sprintf("reg-%s", userId)
}
