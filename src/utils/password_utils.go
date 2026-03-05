package utils

import (
	"strconv"
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

// this function compares the otp in the user request and the otp from the database and returns a bolean value
func CompareUserRequestOTP(OtpFromDB uint64, OtpFromRequest string) bool {
	parsedOTP, err := strconv.ParseUint(OtpFromRequest, 10, 64)
	if err != nil {
		return false
	}
	return OtpFromDB == parsedOTP
}

// this funciton compares the time of request and the expiry time of the otp from the db and returns if the otp in the request is expired or not in a boolean value
func CheckOtpExpiry(otpExpiresAtEpoch uint64, requestArrivalTime time.Time) bool {
	requestArrivedAtEpoch := uint64(requestArrivalTime.Unix())
	return requestArrivedAtEpoch <= otpExpiresAtEpoch
}
