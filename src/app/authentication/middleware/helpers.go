package middleware

import (
	"authentication/commons"
	"errors"
	"stock_broker_application/src/utils"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// extracts token claims
func ValidateToken(ctx *gin.Context) error {
	//extracting the header
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return errors.New("missing authorization header")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return errors.New("invalid authorization format")
	}
	tokenString := parts[1]

	//validatig the token signature
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(utils.SecretKey.AccessSecretKey), nil
	})
	if err != nil {
		return err
	}

	//extracting claims for verification
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {

		//checking expiry
		exp, ok := claims["exp"].(float64)
		if !ok {
			return errors.New("invalid exp claim")
		}
		if int64(exp) < time.Now().Unix() {
			return errors.New("token expired")
		}

		//verifying purpose of the token
		purpose, ok := claims["purpose"].(string)
		if !ok {
			return errors.New("invalid purpose claim")
		}
		if purpose != "password_reset" {
			return errors.New("unauthorized")
		}

		//extracting username and binding with the context
		username, ok := claims["sub"].(string)
		if !ok {
			return errors.New("invalid sub claims")
		}
		ctx.Set(commons.Username, username)
		return nil
	}

	return errors.New("invalid token, cannot extract claims")
}

// // checks if a token is valid or not
// func ValidToken(ctx *gin.Context) error {
// 	claims, err := ExtractTokenClaims(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	//checking expiry of the token
// 	exp, ok := claims["exp"].(float64)
// 	if !ok {
// 		return errors.New("invalid exp claim")
// 	}
// 	if int64(exp) < time.Now().Unix() {
// 		return errors.New("token expired")
// 	}

// 	//verifying purpose of the token
// 	purpose, ok := claims["purpose"].(string)
// 	if !ok {
// 		return errors.New("invalid purpose claim")
// 	}
// 	if purpose != "password_reset" {
// 		return errors.New("unauthorized")
// 	}

// 	//extracting username and binding with the context
// 	username, ok := claims["sub"].(string)
// 	if !ok {
// 		return errors.New("invalid sub claims")
// 	}
// 	ctx.Set(commons.Username, username)

// 	return nil
// }
