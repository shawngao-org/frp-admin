package util

import (
	"errors"
	"frp-admin/config"
	"frp-admin/logger"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func GenerateToken(id int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = id
	claims["exp"] = time.Now().Add(time.Duration(config.Conf.Security.Jwt.Timeout) * time.Second).Unix()
	secret := []byte(config.Conf.Security.Jwt.Secret)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		logger.LogErr("Token generate")
		logger.LogErr("%s", err)
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, errors.New("invalid token")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(config.Conf.Security.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
