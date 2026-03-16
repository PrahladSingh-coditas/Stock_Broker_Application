package handlers

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/models"
	"net/http"
	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils/validations"

	"github.com/gin-gonic/gin"
)

type ChangePasswordHandler struct {
	service *business.ChangePasswordService
}

func NewChangePasswordHandler(service *business.ChangePasswordService) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		service: service,
	}
}

// HandleCreaterUser handles the user signin request.
// @Summary Change Password
// @Description Authenticates user credentials and returns success if valid
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BFFChangePasswordRequest true "User Change Password Request"
// @Success 200 {string} string "Password changed in successfully"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 401 {object} models.ErrorAPIResponse "Invalid Uername"
// @Failure 500 {object} models.ErrorAPIResponse "Authentication failed"
// @Router /api/auth/changepassword [post]
func (controller *ChangePasswordHandler) HandleChangePassword(ctx *gin.Context) {

	var bffChangePasswordRequest models.BFFChangePasswordRequest
	if err := ctx.ShouldBind(&bffChangePasswordRequest); err != nil {
		errorMsgs := genericModels.ErrorMessage{
			Key:          "request",
			ErrorMessage: constants.InvalidPayloadError,
		}
		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.InvalidPayloadError,
		})
		return
	}

	if err := validations.GetBFFValidator().Struct(&bffChangePasswordRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)
		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}

	username := ctx.GetString(constants.Username)
	err := controller.service.ChangePassword(ctx, ctx.Request.Context(), bffChangePasswordRequest, username)
	if err != nil {

		if err.Error() == constants.UserNotFoundError {
			ctx.IndentedJSON(http.StatusUnauthorized, genericModels.ErrorAPIResponse{
				Error: constants.UserNotFoundError,
			})
			return
		}

		ctx.IndentedJSON(http.StatusInternalServerError, genericModels.ErrorAPIResponse{
			Error: constants.AuthenticationFailedError,
		})
		return
	}
	ctx.IndentedJSON(http.StatusOK, constants.PasswordChangeSuccessMsg)

}
