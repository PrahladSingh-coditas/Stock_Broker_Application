package utils

import (
	"stock_broker_application/src/constants"
	"stock_broker_application/src/models"
	"stock_broker_application/src/utils/configs"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var secretKey *models.JWT

func GenerateToken(username string, purpose string) (string, error) {

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     username,
		"purpose": purpose,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(secretKey.AccessSecretKey))
	if err != nil {
		return "", err
	}

	// refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"username": username,
	// 	"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
	// })

	// refreshTokenString, err := refreshToken.SignedString([]byte(secretKey.RefreshSecretKey))
	// if err != nil {
	// 	return "", "", err
	// }
	return accessTokenString, nil
}

func InitJWTConfig(configPath string) error {
	var err error
	secretKey, err = configs.LoadConfig[models.JWT](configPath, constants.JWT, constants.Yaml)
	if err != nil {
		return err
	}
	return nil
}
