package handlers

import (
	"authentication/business"
	"authentication/commons"
	authConstants "authentication/commons/constants"
	"net/http"
	genericModels "stock_broker_application/src/models"
	"strings"
	"time"

	"stock_broker_application/src/constants"

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
	redisKey := ctx.GetString(commons.RedisKey)

	err := controller.service.LogoutUser(ctx, tokenString, tokenExpiry,redisKey)
	if err != nil {
		errorString := err.Error()

		// 500 error for redis fail
		if strings.Contains(errorString, constants.ErrRedisInitFailed) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Redis,
					ErrorMessage: constants.ErrRedisInitFailed,
				},
				Error: constants.ErrInternalServer,
			}

			logger.WithFields(logrus.Fields{
				"user":    tokenString,
				"latency": time.Since(start).Milliseconds(),
			}).Info(authConstants.ErrInternalServer)

			ctx.IndentedJSON(http.StatusInternalServerError, errorResponse)
			return
		}

		//401 unauthorized if failed to blacklist token or if token is already expired
		if strings.Contains(errorString, authConstants.ErrFailedToBlacklist) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Token,
					ErrorMessage: authConstants.ErrFailedToBlacklist,
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

		//500 error
		logger.WithFields(logrus.Fields{
			"user":    tokenString,
			"latency": time.Since(start).Milliseconds(),
		}).Info(constants.ErrInternalServer)

		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.IndentedJSON(http.StatusOK, authConstants.UserLogoutSuccess)
}
