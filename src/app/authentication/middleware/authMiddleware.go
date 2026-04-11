package middleware

import (
	"authentication/commons"
	"authentication/commons/constants"
	"authentication/models"
	"errors"
	"fmt"
	"net/http"
	genericConstants "stock_broker_application/src/constants"
	"stock_broker_application/src/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//extracting the header
		header := ctx.GetHeader(genericConstants.Authorization)
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          genericConstants.Token,
					ErrorMessage: genericConstants.InvalidTokenError,
				},
				Error: genericConstants.AuthHeaderMissingError,
			})
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != genericConstants.Bearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          genericConstants.Token,
					ErrorMessage: genericConstants.InvalidTokenError,
				},
				Error: genericConstants.InvalidAuthFormatError,
			})
			return
		}

		tokenString := parts[1]

		//validating if token is already blacklisted
		key := fmt.Sprintf("BLACKLISTED_TOKEN_%s", tokenString)
		val, err := redisClient.Get(ctx, key).Result()
		if err == nil && val == "1" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          genericConstants.Token,
					ErrorMessage: genericConstants.InvalidTokenError,
				},
				Error: constants.UserAlreadyLoggedoutError,
			})
			return
		}
		if err != nil && err != redis.Nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		//validatig the token signature
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New(genericConstants.UnexpectedSigningMethod)
			}
			return []byte(utils.SecretKey.AccessSecretKey), nil
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		ctx.Set(genericConstants.Token, tokenString)

		//extracting claims for verification
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          genericConstants.Token,
					ErrorMessage: genericConstants.InvalidTokenError,
				},
				Error: genericConstants.InvalidTokenError,
			})
			return
		}

		//checking expiry
		expiry, ok := claims[genericConstants.Exp].(float64)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          genericConstants.Token,
					ErrorMessage: genericConstants.InvalidTokenError,
				},
				Error: genericConstants.InvalidExpClaimsError,
			})
			return
		}
		if int64(expiry) < time.Now().Unix() {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          genericConstants.Token,
					ErrorMessage: genericConstants.InvalidTokenError,
				},
				Error: genericConstants.TokenExpiredError,
			})
			return
		}

		ctx.Set(commons.TokenExpiry, int64(expiry))

		//extracting username and binding with the context
		username, ok := claims[genericConstants.Sub].(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          genericConstants.Token,
					ErrorMessage: genericConstants.InvalidTokenError,
				},
				Error: genericConstants.InvalidSubClaimsError,
			})
			return
		}

		ctx.Set(commons.Username, username)
		ctx.Next()
	}
}
