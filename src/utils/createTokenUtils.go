package utils

import (
	"errors"
	"log"
	"stock_broker_application/src/constants"
	"stock_broker_application/src/models"
	"stock_broker_application/src/utils/configs"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var secretKey *models.JWT

func GenerateToken(username string) (string, string, error) {

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Minute * 2).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(secretKey.AccessSecretKey))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Minute * 1).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey.RefreshSecretKey))
	if err != nil {
		return "", "", err
	}
	return accessTokenString, refreshTokenString, nil
}

func InitJWTConfig(configPath string) error {
	var err error
	secretKey, err = configs.LoadConfig[models.JWT](configPath, constants.JWT, constants.Yaml)
	if err != nil {
		return err
	}
	return nil
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey.AccessSecretKey), nil
	})
	log.Println(err)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiry, ok := claims["exp"].(float64)
		if !ok {
			return "", errors.New("expiry missing in token")
		}
		exp := strconv.FormatFloat(expiry, 'f', -1, 64)
		return exp, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("username missing in token")
		}
		return username, nil
	}
	return "", errors.New("Invalid token")
}
