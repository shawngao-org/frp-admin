package util

import (
	"crypto"
	hmac2 "crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"frp-admin/config"
	bcrypt2 "golang.org/x/crypto/bcrypt"
	"hash"
	"strings"
)

func PasswordEncrypt(pwd string) (string, error) {
	method := config.Conf.Security.Password.Method
	switch strings.ToUpper(method) {
	case "HMAC512":
		return hmac(pwd, 512), nil
	case "HMAC256":
		return hmac(pwd, 256), nil
	case "SHA224":
		return sha(pwd, 224), nil
	case "SHA256":
		return sha(pwd, 256), nil
	case "SHA384":
		return sha(pwd, 384), nil
	case "SHA512":
		return sha(pwd, 512), nil
	case "BCRYPT":
		return bcrypt(pwd)
	}
	return "", errors.New("non-existent encryption method")
}

func hmac(pwd string, m int) string {
	var hmacObj hash.Hash
	if m == 256 {
		hmacObj = hmac2.New(sha256.New, []byte(config.Conf.Security.Password.Secret))
	} else {
		hmacObj = hmac2.New(sha512.New, []byte(config.Conf.Security.Password.Secret))
	}
	hmacObj.Write([]byte(pwd))
	sign := hex.EncodeToString(hmacObj.Sum(nil))
	return sign
}

func sha(pwd string, m int) string {
	var shaObj hash.Hash
	switch m {
	case 224:
		shaObj = crypto.SHA224.New()
	case 256:
		shaObj = crypto.SHA256.New()
	case 384:
		shaObj = crypto.SHA384.New()
	case 512:
		shaObj = crypto.SHA512.New()
	}
	shaObj.Write([]byte(pwd))
	sign := hex.EncodeToString(shaObj.Sum(nil))
	return sign
}

func bcrypt(pwd string) (string, error) {
	cost := config.Conf.Security.Password.Cost
	bytes, err := bcrypt2.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckBcrypt(pwd string, key string) error {
	return bcrypt2.CompareHashAndPassword([]byte(key), []byte(pwd))
}
