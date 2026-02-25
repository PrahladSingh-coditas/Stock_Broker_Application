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

// we have created a struct that have a pointer field to the service
type ForgotPasswordHandler struct {
	service *business.ForgotPasswordService
}

func NewForgotPasswordHandler(service *business.ForgotPasswordService) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		service: service,
	}
}

// HandleCreaterUser handles the user signin request.
// @Summary Forgot Password
// @Description Authenticates user credentials to generate OTP for forgot password
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.BFFForgotPasswordRequest true "Request OTP for forgot password"
// @Success 200 {string} string "OTP Sent successfully"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 401 {object} models.ErrorAPIResponse "Invalid credentials"
// @Failure 404 {object} models.ErrorAPIResponse "User Not Found"
// @Failure 500 {object} models.ErrorAPIResponse "Authentication failed"
// @Router /api/auth/forgotpassword [post]
func (controller *ForgotPasswordHandler) HandleForgotPassword(ctx *gin.Context) {
	var bffForgotPasswordRequest models.BFFForgotPasswordRequest
	//check for api binding errors
	if err := ctx.ShouldBind(&bffForgotPasswordRequest); err != nil {
		errorMsgs := genericModels.ErrorMessage{
			Key:          "request",
			ErrorMessage: constants.ErrInvalidPayload,
		}
		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	// valiation errors
	if err := validations.GetBFFValidator().Struct(&bffForgotPasswordRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)
		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}

	// then we call the function that is in service
	err := controller.service.ForgotPassword(ctx, ctx.Request.Context(), bffForgotPasswordRequest)
	if err != nil {
		if err.Error() == constants.UserNotFoundError {
			ctx.JSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Error: constants.UserNotFoundError,
			})
			return
		}
		ctx.JSON(http.StatusUnauthorized, genericModels.ErrorAPIResponse{
			Error: constants.InvalidCredentialsError,
		})
		return
	}
	ctx.IndentedJSON(http.StatusOK, constants.ForgotPasswordGenerateOtpSuccessMsg)
}
