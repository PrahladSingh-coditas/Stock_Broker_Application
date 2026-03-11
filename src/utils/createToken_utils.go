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

	return accessTokenString, nil
}

func ParseToken(tokenstring string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.WrongSigningAlgorithmError
		}

		return []byte(secretKey.AccessSecretKey), nil
	})

	if err != nil {
		return nil, constants.ParsingFailedError
	}

	return token, nil
}

func VerifyToken(token *jwt.Token) (jwt.MapClaims, error) {

	if !token.Valid {
		return nil, constants.TokenInvalidError
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, constants.ClaimMappingFailedError
	}

	return claims, nil
}

func InitJWTConfig(configPath string) error {
	var err error
	secretKey, err = configs.LoadConfig[models.JWT](configPath, constants.JWT, constants.Yaml)
	if err != nil {
		return err
	}
	return nil
}
