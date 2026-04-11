package handlers

import (
	"authentication/business"
	"authentication/commons"
	"authentication/commons/constants"
	"authentication/models"
	"errors"
	"net/http"

	genericConstants "stock_broker_application/src/constants"
	"stock_broker_application/src/utils/validations"

	"github.com/gin-gonic/gin"
)

type changePasswordHandler struct {
	changePasswordService *business.ChangePasswordService
}

func NewChangePasswordHandler(changePasswordService *business.ChangePasswordService) *changePasswordHandler {
	return &changePasswordHandler{
		changePasswordService: changePasswordService,
	}
}

// this fucntion handles user requests with password and JWT for changing password
// Handles password-change functionality
// @Summary Changes user's password
// @Description verifies the JWT and user request and changes user's password
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.BFFChangePasswordRequest true "User Change-Password Request"
// @Success 200 {object} models.BFFChangePasswordResponse "Password Changed successfully"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 401 {object} models.ErrorAPIResponse "Case A: JWT expired / Case B: JWT invalid / Case C: Unauthorized request"
// @Failure 404 {object} models.ErrorAPIResponse "User does not exist"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Security BearerAuth
// @Router /api/auth/change-password [post]
func (controller changePasswordHandler) HandleChangePassword(ctx *gin.Context) {
	var bffChangePasswordRequest models.BFFChangePasswordRequest

	if bindingErr := ctx.ShouldBindJSON(&bffChangePasswordRequest); bindingErr != nil {
		errorMessage := models.ErrorMessage{
			Key:          bindingErr.Error(),
			ErrorMessage: constants.ErrUnexpectedValue,
		}
		ctx.IndentedJSON(http.StatusBadRequest, models.ErrorAPIResponse{
			Message: errorMessage,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	if validationErr := validations.GetBFFValidator().Struct(&bffChangePasswordRequest); validationErr != nil {
		validationErrors, _ := validations.FormatValidationErrors(validationErr)
		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}

	username := ctx.GetString("username")

	errChangePassword := controller.changePasswordService.ChangePassword(ctx, username, bffChangePasswordRequest.NewPassword, bffChangePasswordRequest.ConfirmPassword)
	if errChangePassword != nil {

		if errors.Is(errChangePassword, commons.NewPasswordMismatchError) {
			paswordsMismatchErr := models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          "new_password",
					ErrorMessage: genericConstants.ErrNewPasswordMatch,
				},
				Error: constants.ErrPasswordChangeFailed,
			}
			ctx.IndentedJSON(http.StatusBadRequest, paswordsMismatchErr)
			return
		}

		if errors.Is(errChangePassword, commons.UserNotFoundError) {
			userNotFoundErr := models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          commons.Username,
					ErrorMessage: constants.ErrUserNotFound,
				},
				Error: constants.ErrPasswordChangeFailed,
			}
			ctx.IndentedJSON(http.StatusNotFound, userNotFoundErr)
			return
		}

		if errors.Is(errChangePassword, commons.InvalidTokenError) {
			invalidTokenErr := models.ErrorAPIResponse{
				Message: models.ErrorMessage{
					Key:          genericConstants.Token,
					ErrorMessage: constants.ErrInvalidToken,
				},
				Error: constants.ErrPasswordChangeFailed,
			}
			ctx.IndentedJSON(http.StatusUnauthorized, invalidTokenErr)
			return
		}
	}

	ctx.IndentedJSON(http.StatusOK, models.BFFChangePasswordResponse{
		Message: constants.PasswordChangedSuccessMsg,
	})
}
