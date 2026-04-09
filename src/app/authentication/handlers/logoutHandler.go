package handlers

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/models"
	"fmt"
	"log"
	"net/http"
	genericModels "stock_broker_application/src/models"
	"strings"
	"time"

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
	fmt.Println()
	fmt.Println()
	duration := expiry - time.Now().Unix()
	log.Println(duration)
	timeToLive := time.Duration(duration) * time.Second
	log.Println(timeToLive)

	token := ctx.GetString(constants.Token)

	err := controller.service.LogoutUser(ctx, ctx.Request.Context(), token, timeToLive)
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

		if strings.Contains(err.Error(), constants.RedisSetOperationError) {
			err := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Redis,
					ErrorMessage: constants.RedisSetOperationError,
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
