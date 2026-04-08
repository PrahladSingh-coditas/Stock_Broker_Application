package middleware

import (
	"authentication/commons"
	"authentication/commons/constants"
	"net/http"
	"stock_broker_application/src/models"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ValidateToken(ctx); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          commons.Token,
					ErrorMessage: err.Error(),
				},
				Error: constants.ErrInvalidToken,
			})
			return
		}
		ctx.Next()
	}
}
