package utils

import (
	"errors"
	"log"
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
		"exp":     time.Now().Add(time.Minute * 2).Unix(),
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

func ValidateToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//to check if it belongs to same family of signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		//if it does then only send the secret key to verify
		return []byte(secretKey.AccessSecretKey), nil
	})

	log.Println(err)
	if err != nil {
		return nil, err
	}

	// to extract claims from token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims,nil
	}

	return nil, errors.New("Invalid token")
}

func ExtractUsername(tokenString string) (string,error){
	claims, err := ValidateToken(tokenString)
	if err !=nil{
		return "",err
	}
	
	username, usernameFetched := claims["sub"].(string)
	if !usernameFetched{
		return "",errors.New("Username missing or invalid")
	}

	return username,nil
}

func ExtractExpiry(tokenString string) (int64,error){
	claims, err := ValidateToken(tokenString)
	if err !=nil{
		return 0,err
	}
	
	tokenExpiry, expiryFetched := claims["exp"].(float64) //jwt expiry is float64
	if !expiryFetched{
		return 0,errors.New("Expiry missing or invalid")
	}

	return int64(tokenExpiry),nil
}