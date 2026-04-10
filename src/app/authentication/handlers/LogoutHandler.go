package handlers

import (
	"authentication/business"
	"authentication/commons"
	"authentication/commons/constants"
	"net/http"
	genericModels "stock_broker_application/src/models"
	"strings"
	"time"

	commonErrors "stock_broker_application/src/constants"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LogoutUserHandler struct {
	service *business.LogoutUserService
}

func NewLogoutUserHandler(service *business.LogoutUserService) *LogoutUserHandler {
	return &LogoutUserHandler{
		service: service,
	}
}

// HandleSigninUser handles the user signin request.
// @Summary Logout User
// @Description Logout the user
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.BFFLogoutUserResponse "logout successfully"
// @Failure 401 {object} models.ErrorAPIResponse "failed to blacklist/unauthorized"
// @Failure 500 {object} models.ErrorAPIResponse "redis connection failed/ internal server error"
// @Router /api/auth/logout [post]
func (controller *LogoutUserHandler) HandleLogoutUser(ctx *gin.Context) {
	start := time.Now()
	logger := logrus.New()

	tokenString := ctx.GetString(commons.Token)
	tokenExpiry := ctx.GetInt64(commons.Expiry)

	err := controller.service.LogoutUser(ctx, tokenString, tokenExpiry)
	if err != nil {
		errorString := err.Error()

		// 500 error for redis fail
		if strings.Contains(errorString, commonErrors.ErrRedisInitFailed) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          "redis",
					ErrorMessage: commonErrors.ErrRedisInitFailed,
				},
				Error: constants.ErrInternalServer,
			}

			logger.WithFields(logrus.Fields{
				"user":    tokenString,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrPasswordMismatch)

			ctx.IndentedJSON(http.StatusUnauthorized, errorResponse)
			return
		}

		//401 unauthorized if failed to blacklist token or if token is already expired
		if strings.Contains(errorString, constants.ErrFailedToBlacklist) || strings.Contains(errorString, constants.ErrTokenExpired) {
			var errorMsg string

			switch {
			case strings.Contains(errorString, constants.ErrFailedToBlacklist):
				errorMsg = constants.ErrFailedToBlacklist
			case strings.Contains(errorString, constants.ErrTokenExpired):
				errorMsg = constants.ErrTokenExpired
			}

			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          "token",
					ErrorMessage: errorMsg,
				},
				Error: constants.ErrUnauthorized,
			}

			logger.WithFields(logrus.Fields{
				"user":    tokenString,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrUnauthorized)

			ctx.IndentedJSON(http.StatusUnauthorized, errorResponse)
			return
		}
		logger.WithFields(logrus.Fields{
			"user":    tokenString,
			"latency": time.Since(start).Milliseconds(),
		}).Info(constants.ErrInternalServer)

		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, constants.UserLogoutSuccess)

}
