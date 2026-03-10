package handlers

import (
	"authentication/business"
	"authentication/commons"
	"authentication/commons/constants"
	"authentication/models"
	"encoding/json"
	"errors"
	"net/http"
	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils/validations"

	"github.com/gin-gonic/gin"
)

type ValidateUserOtpHandler struct {
	service *business.ValidateUserOtpService
}

func NewValidateUserOtpHandler(service *business.ValidateUserOtpService) *ValidateUserOtpHandler {
	return &ValidateUserOtpHandler{
		service: service,
	}
}

// this fucntion handles user requests and responses by
// Handles user OTP validation
// @Summary Validates user OTP
// @Description Validates user OTP and return clear success/ failure message
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.BFFValidateUserOtpRequest true "User OTP Validation Request"
// @Success 200 {object} models.BFFValidateUserOtpResponse "OTP validation successful"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 401 {object} models.ErrorAPIResponse "Case A: Incorrect OTP / Case B: Expired OTP"
// @Failure 404 {object} models.ErrorAPIResponse "User does not exist"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /api/auth/validateotp [post]
func (controller *ValidateUserOtpHandler) HandleValidateUserOtp(ctx *gin.Context) {
	var bffValidateUserOtpRequest models.BFFValidateUserOtpRequest

	var bFFValidateUserOtpResponse models.BFFValidateUserOtpResponse
	if err := ctx.ShouldBind(&bffValidateUserOtpRequest); err != nil {
		errorMessage := genericModels.ErrorMessage{
			Key:          err.(*json.UnmarshalTypeError).Field,
			ErrorMessage: constants.ErrUnexpectedValue,
		}

		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMessage,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	if err := validations.GetBFFValidator().Struct(&bffValidateUserOtpRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)
		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}

	token, err := controller.service.ValidateUserOtp(ctx, ctx.Request.Context(), bffValidateUserOtpRequest)
	if err != nil {
		if errors.Is(err, commons.UserNotFoundError) {
			err := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.Username,
					ErrorMessage: constants.UserNotFoundError,
				},
				Error: constants.AuthenticationFailedError,
			}
			ctx.IndentedJSON(http.StatusNotFound, err)
			return
		}

		if errors.Is(err, commons.IncorrectOTPError) {
			err := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.Otp,
					ErrorMessage: constants.IncorrectOTPError,
				},
				Error: constants.AuthenticationFailedError,
			}
			ctx.IndentedJSON(http.StatusUnauthorized, err)
			return
		}

		if errors.Is(err, commons.OtpExpiredError) {
			err := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.Otp,
					ErrorMessage: constants.OtpExpiredError,
				},
				Error: constants.AuthenticationFailedError,
			}
			ctx.IndentedJSON(http.StatusUnauthorized, err)
			return
		}

		if errors.Is(err, commons.TokenGenerationError) {
			err := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.Token,
					ErrorMessage: constants.TokenGenerationError,
				},
				Error: constants.AuthenticationFailedError,
			}
			ctx.IndentedJSON(http.StatusBadRequest, err)
			return
		}

		ctx.IndentedJSON(http.StatusUnauthorized, genericModels.ErrorAPIResponse{
			Error: constants.SigninFailedError,
		})
		return
	}

	bFFValidateUserOtpResponse.Message = constants.TokenGeneratedSuccessMsg
	bFFValidateUserOtpResponse.Token = token

	ctx.IndentedJSON(http.StatusOK, bFFValidateUserOtpResponse)
}
