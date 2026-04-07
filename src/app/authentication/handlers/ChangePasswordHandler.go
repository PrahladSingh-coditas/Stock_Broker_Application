package handlers

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/models"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils/validations"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ChangePasswordHandler struct {
	service *business.ChangePasswordService
}

func NewChangePasswordHandler(service *business.ChangePasswordService) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		service: service,
	}
}
//[SDTAP SPSFF..R]
// HandleSigninUser handles the user signin request.
// @Summary Change password
// @Description Change password and return message
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BFFChangePasswordRequest true "Change Password Request"
// @Success 200 {object} models.BFFChangePasswordResponse "Password changed successfully"
// @Failure 400 {object} models.ErrorAPIResponse  "Invalid Input Payload"
// @Failure 401 {object} models.ErrorAPIResponse "Password mismatch with confirm password"
// @Failure 404 {object} models.ErrorAPIResponse "Invalid Username"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /api/auth/change-password [post]
func (controller *ChangePasswordHandler) HandleChangePassword(ctx *gin.Context) {
	start := time.Now()
	logger := logrus.New()

	var bffChangePasswordRequest models.BFFChangePasswordRequest

	if err := ctx.ShouldBind(&bffChangePasswordRequest); err != nil {
		errorMsgs := genericModels.ErrorMessage{
			Key:          err.(*json.UnmarshalTypeError).Field,
			ErrorMessage: constants.ErrUnexpectedValue,
		}

		logger.WithFields(logrus.Fields{
			"password": bffChangePasswordRequest.NewPassword,
			"latency":  time.Since(start).Milliseconds(),
		}).Info(constants.ErrBinding)

		ctx.JSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	// 400 error
	if err := validations.GetBFFValidator().Struct(&bffChangePasswordRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)

		logger.WithFields(logrus.Fields{
			"password": bffChangePasswordRequest.NewPassword,
			"latency":  time.Since(start).Milliseconds(),
		}).Info(constants.ErrUnexpectedValue)

		ctx.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	username := ctx.GetString("username")
	

	err := controller.service.ChangePassword(ctx, ctx.Request.Context(), username, bffChangePasswordRequest.NewPassword, bffChangePasswordRequest.ConfirmPassword)

	if err != nil {
		errorString := err.Error()

		// 401 error
		if strings.Contains(errorString, constants.ErrPasswordMismatch) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          "password",
					ErrorMessage: constants.ErrPasswordMismatch,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				"user":    username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrPasswordMismatch)

			ctx.IndentedJSON(http.StatusUnauthorized, errorResponse)
			return
		}

		//404 error
		if strings.Contains(errorString, constants.ErrNoRowsAffected) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          "username",
					ErrorMessage: constants.ErrUserNotFound,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				"user":    username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrUserNotFound)

			ctx.IndentedJSON(http.StatusNotFound, errorResponse)
			return

		}

		// 500 errors
		//hashing failed
		if strings.Contains(errorString, constants.ErrFailedToEncrypt) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          "password",
					ErrorMessage: constants.ErrFailedToEncrypt,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				"user":    username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrFailedToEncrypt)

			ctx.IndentedJSON(http.StatusInternalServerError, errorResponse)
			return
		}

		//server error
		errorResponse := genericModels.ErrorAPIResponse{
			Message: genericModels.ErrorMessage{
				Key:          "server",
				ErrorMessage: constants.ErrInternalServer,
			},
			Error: constants.ErrInternalServer,
		}

		logger.WithFields(logrus.Fields{
			"user":    username,
			"latency": time.Since(start).Milliseconds(),
		}).Info(constants.ErrInternalServer)

		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return

	}

	// 200 status OK
	logger.WithFields(logrus.Fields{
		"user":    username,
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.UserLoggedInSuccessMsg, constants.UserOtpGeneratedSuccess)

	ctx.IndentedJSON(http.StatusOK, models.BFFChangePasswordResponse{
		Message: constants.PasswordChangedSuccess,
	})
}
