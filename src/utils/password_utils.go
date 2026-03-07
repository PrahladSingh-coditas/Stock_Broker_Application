package utils

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func CompareHashPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func CompareUserRequestOTP(OtpFromDB uint64, OtpFromRequest uint64) bool {
	return OtpFromDB == OtpFromRequest
}

func CheckOtpExpiry(otpExpiresAt uint64, requestTime time.Time) bool {
	convertedotpExpiresAt := time.Unix(int64(otpExpiresAt), 0)
	if requestTime.After(convertedotpExpiresAt) {
		return false
	}
	return true
}
