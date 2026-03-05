package handlers

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/models"
	"encoding/json"
	"errors"
	"net/http"
	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils/validations"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
// @Router /api/auth/validate-otp [post]
func (controller *ValidateUserOtpHandler) HandleValidateUserOtp(ctx *gin.Context) {

	start := time.Now()
	logger := logrus.New()

	var bffValidateUserOtpRequest models.BFFValidateUserOtpRequest
	var userResponse models.BFFValidateUserOtpResponse

	if bindingError := ctx.ShouldBind(&bffValidateUserOtpRequest); bindingError != nil {
		errorMessage := genericModels.ErrorMessage{
			Key:          bindingError.(*json.UnmarshalTypeError).Field,
			ErrorMessage: constants.ErrUnexpectedValue,
		}

		logger.WithFields(logrus.Fields{
			constants.User:    bffValidateUserOtpRequest.Username,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info(constants.ErrBindingFailed)

		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMessage,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	if validationError := validations.GetBFFValidator().Struct(&bffValidateUserOtpRequest); validationError != nil {
		validationErrors, _ := validations.FormatValidationErrors(validationError)

		logger.WithFields(logrus.Fields{
			constants.User:    bffValidateUserOtpRequest.Username,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info(constants.ErrValidationFailed)

		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}

	token,err := controller.service.ValidateUserOtp(ctx, ctx.Request.Context(), bffValidateUserOtpRequest)
	if err != nil {
		if errors.Is(err, constants.UserNotFoundError) {
			errorUserNotFoundResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.User,
					ErrorMessage: constants.ErrUserNotFoundMsg,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				constants.User:    bffValidateUserOtpRequest.Username,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info(constants.ErrUserNotFoundMsg)

			ctx.IndentedJSON(http.StatusNotFound, errorUserNotFoundResponse)
			return
		} else if errors.Is(err, constants.IncorrectOTPError) {
			errorIncorrectOtpResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Otp,
					ErrorMessage: constants.ErrIncorrectOtp,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				constants.User:    bffValidateUserOtpRequest.Username,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info(constants.ErrIncorrectOtp)

			ctx.IndentedJSON(http.StatusUnauthorized, errorIncorrectOtpResponse)
			return
		} else if errors.Is(err, constants.OtpExpiredError) {
			errorExpiredOtpResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Otp,
					ErrorMessage: constants.ErrExpiredOtp,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				constants.User:    bffValidateUserOtpRequest.Username,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info(constants.ErrExpiredOtp)

			ctx.IndentedJSON(http.StatusUnauthorized, errorExpiredOtpResponse)
			return
		} else if errors.Is(err, constants.TokenCreationFailedError) {
			errorExpiredOtpResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Token,
					ErrorMessage: constants.ErrTokenCreationFailed,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				constants.User:    bffValidateUserOtpRequest.Username,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info(constants.ErrTokenCreationFailed)

			ctx.IndentedJSON(http.StatusInternalServerError, errorExpiredOtpResponse)
			return
		}

		logger.WithFields(logrus.Fields{
			constants.User:    bffValidateUserOtpRequest.Username,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info(constants.ErrAuthenticationFailed)

		ctx.IndentedJSON(http.StatusInternalServerError, genericModels.ErrorAPIResponse{
			Error: constants.ErrSignInFailed,
		})
		return
	}

	logger.WithFields(logrus.Fields{
			constants.User:    bffValidateUserOtpRequest.Username,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info(constants.OtpValidatedSuccessMsg)

	
	userResponse.AccessToken=token
	userResponse.Message=constants.OtpValidatedSuccessMsg

	ctx.IndentedJSON(http.StatusOK, userResponse)
}
