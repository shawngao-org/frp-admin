package util

import (
	"regexp"
)

func IsPassValid(msg string) bool {
	return len(msg) == 0
}

func UnifiedVerificationOfBasicUserInfo(name string, email string, password string) string {
	nameStr := ValidateUsername(name)
	if len(nameStr) != 0 {
		return nameStr
	}
	emailStr := ValidateEmail(email)
	if len(emailStr) != 0 {
		return emailStr
	}
	passwdStr := ValidatePassword(password)
	if len(passwdStr) != 0 {
		return passwdStr
	}
	return ""
}

func VerificationEmailAndPassword(email string, password string) string {
	emailStr := ValidateEmail(email)
	if len(emailStr) != 0 {
		return emailStr
	}
	passwdStr := ValidatePassword(password)
	if len(passwdStr) != 0 {
		return passwdStr
	}
	return ""
}

func VerificationUsernameAndPassword(name string, password string) string {
	nameStr := ValidateUsername(name)
	if len(nameStr) != 0 {
		return nameStr
	}
	passwdStr := ValidatePassword(password)
	if len(passwdStr) != 0 {
		return passwdStr
	}
	return ""
}

func ValidateUsername(name string) string {
	nameReg := regexp.MustCompile(`^[a-zA-Z]\w{3,12}$`)
	if !nameReg.MatchString(name) {
		return "The user name can contain only 4 to 12 letters, digits, and underscores (_), and must start with a letter."
	}
	return ""
}

func ValidateEmail(email string) string {
	emailReg := regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	if !emailReg.MatchString(email) {
		return "Invalid email address."
	}
	return ""
}

func ValidatePassword(passwd string) string {
	hasNumber := bool2Int(regexp.MustCompile(`[0-9]`).MatchString(passwd))
	hasLower := bool2Int(regexp.MustCompile(`[a-z]`).MatchString(passwd))
	hasUpper := bool2Int(regexp.MustCompile(`[A-Z]`).MatchString(passwd))
	lengthValid := bool2Int(regexp.MustCompile(`^.{8,32}$`).MatchString(passwd))
	sum := 1 << hasNumber << hasLower << hasUpper << lengthValid
	if sum < 8 {
		return "The password must contain 8 to 16 characters, including two types of uppercase letters, lowercase letters, digits, and symbols."
	}
	return ""
}

func bool2Int(flag bool) int {
	if flag {
		return 1
	}
	return 0
}
