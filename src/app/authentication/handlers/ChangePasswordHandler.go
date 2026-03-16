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

type ChangePasswordHandler struct {
	changePasswordService *business.ChangePasswordService
}

func NewChangePasswordHandler(changePasswordService *business.ChangePasswordService) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		changePasswordService: changePasswordService,
	}
}

// HandleChangePassword handles the change password
// @Summary Change user password
// @Description Handles the user's change password request and stores the new password in database
// @Tags User
// @Accept json
// @produce json
// @Security BearerAuth
// @Param request body models.BFFChangePasswordRequest true "Change Password Request"
// @Success 200 {object} models.BFFChangePasswordResponse "Password changed successfully"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 404 {object} models.ErrorAPIResponse "User does not exist"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /api/auth/change-password [patch]
func (handler *ChangePasswordHandler) HandleChangePassword(ctx *gin.Context) {
	start := time.Now()
	logger := logrus.New()
	var bffChangePasswordRequest models.BFFChangePasswordRequest

	username := ctx.GetString("username")

	if err := ctx.ShouldBind(&bffChangePasswordRequest); err != nil {
		errorMsgs := genericModels.ErrorMessage{Key: err.(*json.UnmarshalTypeError).Field, ErrorMessage: constants.ErrUnexpectedValue}

		logger.WithFields(logrus.Fields{
			constants.User:    username,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info(constants.ErrBindingFailed)

		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.ErrInvalidPayload,
		})

		return
	}

	if err := validations.GetBFFValidator().Struct(&bffChangePasswordRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)

		logger.WithFields(logrus.Fields{
			constants.User:    username,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info(constants.ErrValidationFailed)

		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}

	err := handler.changePasswordService.ServiceChangePassword(ctx, ctx.Request.Context(), username, bffChangePasswordRequest.NewPassword, logger)

	if err != nil {
		if errors.Is(err, errors.New(constants.ErrFailedToEncrypt)) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Password,
					ErrorMessage: constants.ErrFailedToEncrypt,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				constants.User:    username,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info(constants.ErrFailedToEncrypt)

			ctx.IndentedJSON(http.StatusInternalServerError, errorResponse)
			return
		} else if errors.Is(err,  errors.New(constants.ErrDatabaseQueryErrorMsg)) {
			errorResponse := genericModels.ErrorAPIResponse{
				Error: constants.ErrDatabaseQueryErrorMsg,
			}

			logger.WithFields(logrus.Fields{
				constants.User:    username,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info(constants.ErrDatabaseQueryErrorMsg)

			ctx.IndentedJSON(http.StatusInternalServerError, errorResponse)
			return
		} else if errors.Is(err, errors.New(constants.ErrUserNotFoundMsg)) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.User,
					ErrorMessage: constants.ErrUserNotFoundMsg,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				constants.User:    username,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info(constants.ErrUserNotFoundMsg)

			ctx.IndentedJSON(http.StatusNotFound, errorResponse)
			return
		} else {

			errorResponse := genericModels.ErrorAPIResponse{
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				constants.User:    username,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info(constants.ErrAuthenticationFailed)

			ctx.IndentedJSON(http.StatusInternalServerError, errorResponse)
			return
		}
	}

	logger.WithFields(logrus.Fields{
		constants.User:    username,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info(constants.PasswordUpdateMsg)

	ctx.IndentedJSON(http.StatusOK, constants.PasswordUpdateMsg)

}
