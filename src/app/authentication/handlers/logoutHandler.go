package handlers

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/models"
	"net/http"
	genericModels "stock_broker_application/src/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type LogoutHandler struct {
	service *business.LogoutUserService
}

func LogoutUserHandler(service *business.LogoutUserService) *LogoutHandler {
	return &LogoutHandler{
		service: service,
	}
}

// HandleUserLogout handles the user logout request.
// @Summary User Logout
// @Description Logout user using Redis by Blacklisting Token
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.BFFLogoutUserResponse "logged out successfully"
// @Failure 401 {object} models.ErrorAPIResponse "Unauthorized"
// @Failure 500 {object} models.ErrorAPIResponse "Authentication failed"
// @Router /api/auth/logout [post]
func (controller *LogoutHandler) HandleUserLogout(ctx *gin.Context) {
	var bffLogoutUserResponse models.BFFLogoutUserResponse
	expiry := ctx.GetInt64(constants.Expiry)
	token := ctx.GetString(constants.Token)

	err := controller.service.LogoutUser(ctx, token, expiry)
	if err != nil {
		if strings.Contains(err.Error(), constants.RedisConnectionFailedError) {
			err := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Redis,
					ErrorMessage: constants.RedisConnectionFailedError,
				},
				Error: constants.AuthenticationFailedError,
			}
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		if strings.Contains(err.Error(), constants.RedisOperationFailedError) {
			err := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Redis,
					ErrorMessage: constants.RedisOperationFailedError,
				},
				Error: constants.AuthenticationFailedError,
			}
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		ctx.IndentedJSON(http.StatusUnauthorized, genericModels.ErrorAPIResponse{
			Error: constants.LogoutFailedError,
		})
		return
	}
	bffLogoutUserResponse.Message = "Logout Successful and Token Blacklisted"
	ctx.IndentedJSON(http.StatusOK, bffLogoutUserResponse)
}
