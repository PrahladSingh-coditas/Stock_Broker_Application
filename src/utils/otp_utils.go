package utils

import (
	"strconv"
	"time"
)

func CompareUserRequestOTP(OtpFromDB uint64, OtpFromRequest string) bool {
	parsedOTP, err := strconv.ParseUint(OtpFromRequest, 10, 64)
	if err != nil {
		return false
	}
	return OtpFromDB == parsedOTP
}

func CheckOtpExpiry(otpExpiresAtEpoch uint64, requestArrivalTime time.Time) bool {
	requestArrivedAtEpoch := uint64(requestArrivalTime.Unix())
	return requestArrivedAtEpoch <= otpExpiresAtEpoch
}
