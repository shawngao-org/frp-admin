package service

import (
	"errors"
	"frp-admin/config"
	"frp-admin/db"
	"frp-admin/entity"
	"frp-admin/util"
	"reflect"
	"strings"
)

func GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	db.Db.First(&user, "email = ?", email)
	if reflect.DeepEqual(user, entity.User{}) {
		return entity.User{}, errors.New("user not found")
	}
	return user, nil
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
