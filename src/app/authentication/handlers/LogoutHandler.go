package handlers

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LogoutUserHandler struct {
	logoutUserService *business.LogoutUserService
}

func NewLogoutUserHandler(logout *business.LogoutUserService) *LogoutUserHandler {
	return &LogoutUserHandler{
		logoutUserService: logout,
	}
}

// Swagger annotations
// @Summary logout user
// @Description handles logout request of the user
// @Tags User
// @produce json
// @Security BearerAuth
// @Success 200 {object} models.BFFLogoutResponse "User logout successfully"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /api/auth/logout [post]
func (handler *LogoutUserHandler) HandleLogoutUser(ctx *gin.Context) {
	logger := logrus.New()
	expiryTime := ctx.GetInt64(constants.ExpiryTime)
	tokenString := ctx.GetString(constants.Token)

	err := handler.logoutUserService.LogoutUser(ctx, logger, tokenString, expiryTime)

	if err != nil {
		if strings.Contains(err.Error(), constants.RedisSetOperationError) {

			logger.Error(constants.RedisSetOperationError)

			ctx.IndentedJSON(http.StatusInternalServerError, models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          constants.Redis,
					ErrorMessage: constants.RedisSetOperationError,
				},
				Error: constants.LogoutFailed,
			})
			return
		}
	}

	logger.Error(constants.LogoutSuccessMsg)

	ctx.IndentedJSON(http.StatusOK, models.BFFLogoutResponse{
		Status: constants.LogoutSuccessMsg,
	})
}
