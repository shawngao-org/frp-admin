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

func UpdateUserById(user entity.User) error {
	var result entity.User
	err := GetById(result, user.Id)
	if err != nil {
		return err
	}
	util.CopyProperties(&user, &result)
	db.Db.Save(result)
	return nil
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
		GroupId:      "47bbe440-dfcb-435f-b7ef-dba7b54a2135",
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

func SendForgetPasswordMail(ctx *gin.Context) {
	email := ctx.PostForm("email")
	user, err := CheckEmail(email)
	if err != nil {
		ErrHandle(ctx, err, http.StatusBadRequest)
		return
	}
	code, err := GenerateResetPasswordTmpCode2Redis(user.Id)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	content := util.MailContent{
		Title: "Reset Password",
		Content: "Click the button below to reset your account password. " +
			"Note: If you did not initiate the request, please ignore this email so as not to cause unnecessary trouble.",
		BtnLink: fmt.Sprintf("%s/reset-password/%s", config.Conf.Server.FrontEndAddr, code),
		BtnText: "Click to change your password",
		Author:  "FRP-Admin",
		Note:    util.DefaultFooterNote,
	}
	util.SendDefaultMail(user.Email, "Reset Password", &content)
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
	_, err = VerifyResetPasswordTmpCodeRedis(user.Id, verifyCode)
	if err != nil {
		ErrHandle(ctx, 123, http.StatusBadRequest)
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
	err = UpdateUserById(user)
	if err != nil {
		ErrHandle(ctx, err, http.StatusInternalServerError)
		return
	}
	OkHandle(ctx)
}

func GenerateResetPasswordTmpCode2Redis(userId string) (string, error) {
	tmpCode := entity.TmpCode{
		Model:  gorm.Model{},
		Id:     utils.NewUUID().String(),
		IsUsed: false,
	}
	result := db.Db.Create(&tmpCode)
	if result.Error != nil {
		return "", result.Error
	}
	redisKey := getResetPasswordRedisKey(userId)
	if _, err := redis.Client.Get(common.Context, redisKey).Result(); err == nil {
		return "", errors.New("sending too often, please try again later")
	}
	status := redis.Client.Set(common.Context, redisKey, tmpCode.Id, 5*time.Minute)
	if status.Err() != nil {
		return "", status.Err()
	}
	return tmpCode.Id, nil
}

func VerifyResetPasswordTmpCodeRedis(userId string, code string) (bool, error) {
	redisKey := getResetPasswordRedisKey(userId)
	val, err := redis.Client.Get(common.Context, redisKey).Result()
	if err != nil {
		return false, err
	}
	if val != code {
		return false, errors.New("the verification code is not correct")
	}
	result, err := GetTmpCodeById(val)
	if err != nil {
		return false, err
	}
	if result.IsUsed {
		return false, errors.New("the verification code is invalid")
	}
	result.IsUsed = true
	err = UpdateTmpCodeById(result)
	if err != nil {
		return false, err
	}
	redis.Client.Del(common.Context, redisKey)
	return true, nil
}

func getResetPasswordRedisKey(userId string) string {
	return fmt.Sprintf("rp-%s", userId)
}
