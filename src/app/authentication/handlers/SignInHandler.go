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

type SigninUserHandler struct {
	service *business.SigninUserService
}

func NewSigninUserHandler(service *business.SigninUserService) *SigninUserHandler {
	return &SigninUserHandler{
		service: service,
	}
}

// HandleSigninUser handles the user signin request.
// @Summary Sign in user
// @Description Authenticate user and return JWT tokens
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.BFFSigninUserRequest true "Signin Request"
// @Success 200 {object} models.BFFSigninUserResponse
// @Failure 400 {object} models.ErrorAPIResponse
// @Failure 401 {object} models.ErrorAPIResponse
// @Router /api/auth/signin [post]
func (controller *SigninUserHandler) HandleSigninUser(ctx *gin.Context) {
	start := time.Now()
	logger := logrus.New()

	var bffSigninUserRequet models.BFFSigninUserRequest
	

	if err := ctx.ShouldBind(&bffSigninUserRequet); err != nil {
		errorMsgs := genericModels.ErrorMessage{
			Key:          err.(*json.UnmarshalTypeError).Field,
			ErrorMessage: constants.ErrUnexpectedValue,
		}

		logger.WithFields(logrus.Fields{
			"user":    bffSigninUserRequet.Username,
			"latency": time.Since(start).Milliseconds(),
		}).Info(constants.ErrBinding)

		ctx.JSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	// 400 error
	if err := validations.GetBFFValidator().Struct(&bffSigninUserRequet); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)

		logger.WithFields(logrus.Fields{
			"user":    bffSigninUserRequet.Username,
			"latency": time.Since(start).Milliseconds(),
		}).Info(constants.ErrUnexpectedValue)

		ctx.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	err := controller.service.SigninUser(ctx, ctx.Request.Context(), bffSigninUserRequet)
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
				"user":    bffSigninUserRequet.Username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrPasswordMismatch)

			ctx.IndentedJSON(http.StatusUnauthorized, errorResponse)
			return
		}

		//404 error
		if strings.Contains(errorString, constants.ErrInvalidEmailorPassword) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          "username",
					ErrorMessage: constants.ErrUserNotFound,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				"user":    bffSigninUserRequet.Username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrUserNotFound)

			ctx.IndentedJSON(http.StatusNotFound, errorResponse)
			return

		}

		// 500 error
		errorResponse := genericModels.ErrorAPIResponse{
			Message: genericModels.ErrorMessage{
				Key:          "server",
				ErrorMessage: constants.ErrInternalServer,
			},
			Error: constants.ErrInternalServer,
		}

		logger.WithFields(logrus.Fields{
			"user":    bffSigninUserRequet.Username,
			"latency": time.Since(start).Milliseconds(),
		}).Info(constants.ErrInternalServer)

		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return

	}

	// 200 status OK
	logger.WithFields(logrus.Fields{
		"user":    bffSigninUserRequet.Username,
		"latency": time.Since(start).Milliseconds(),
	}).Info(constants.UserLoggedInSuccessMsg, constants.UserOtpGeneratedSuccess)

	ctx.IndentedJSON(http.StatusOK, models.BFFSigninUserResponse{
		Message:      constants.UserLoggedInSuccessMsg,
		OtpSent:      constants.UserOtpGeneratedSuccess,
		OtpExpiresAt: constants.UserOtpExpiryMsg,
	})
}