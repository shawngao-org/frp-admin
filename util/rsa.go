package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"frp-admin/logger"
)

func Decrypted(key string, base64String string) string {
	originStr, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		logger.LogErr("Unable to decode public key.")
		logger.LogErr("%s", err)
		return ""
	}
	cipherText := originStr
	privateKeyBlock, _ := pem.Decode([]byte(key))
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

func Encrypted(key string, str string) string {
	plainText := []byte(str)
	publicKeyBlock, _ := pem.Decode([]byte(key))
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
