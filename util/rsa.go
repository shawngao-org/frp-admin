package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"frp-admin/config"
	"frp-admin/logger"
	"os"
)

var PrivateKey string
var PublicKey string

func GetKeys() {
	pri, err := os.ReadFile(config.Conf.Security.Rsa.Private)
	if err != nil {
		logger.LogErr("Unable to read private key file.")
		os.Exit(-1)
	}
	PrivateKey = string(pri)
	pub, err := os.ReadFile(config.Conf.Security.Rsa.Public)
	if err != nil {
		logger.LogErr("Unable to read public key file.")
		os.Exit(-1)
	}
	PublicKey = string(pub)
}

func Decrypted(str string) string {
	originStr, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		logger.LogErr("Unable to decode public key.")
		logger.LogErr("%s", err)
		return ""
	}
	cipherText := originStr
	privateKeyBlock, _ := pem.Decode([]byte(PrivateKey))
	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		logger.LogErr("Unable to parse private key.")
		logger.LogErr("%s", err)
		return ""
	}
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		logger.LogErr("Wrong private key type.")
		return ""
	}
	decryptedText, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, cipherText)
	if err != nil {
		logger.LogErr("%s", err)
		return ""
	}
	return string(decryptedText)
}

func Encrypted(str string) string {
	plainText := []byte(str)
	publicKeyBlock, _ := pem.Decode([]byte(PublicKey))
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		logger.LogErr("Unable to parse public key.")
		logger.LogErr("%s", err)
		return ""
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		logger.LogErr("Wrong public key type.")
		return ""
	}
	encryptedText, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, plainText)
	if err != nil {
		logger.LogErr("%s", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(encryptedText)
}
