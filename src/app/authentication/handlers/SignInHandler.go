package handlers

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/models"
	"encoding/json"
	"net/http"
	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils/validations"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SignInUserHandler struct {
	service *business.SignInUserService
}

func NewSignInUserHandler(service *business.SignInUserService) *SignInUserHandler {
	return &SignInUserHandler{
		service: service,
	}
}

// HandleSignInUser handles the user sign in request.
// @Summary Sign In User
// @Description Handles the user sign and checks data from database
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.BFFSignInUserRequest true "User Sign In Request"
// @Success 200 {string} string "User Sign In successfully"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 401 {object} models.ErrorAPIResponse "Invalid username or password"
// @Failure 404 {object} models.ErrorAPIResponse "User record not found"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /api/auth/signin [post]
func (controller *SignInUserHandler) HandleSignInUser(ctx *gin.Context) {

	start := time.Now()
	logger := logrus.New()

	var bffSignInRequest models.BFFSignInUserRequest
	if bindingError := ctx.ShouldBind(&bffSignInRequest); bindingError != nil {
		errorMsgs := genericModels.ErrorMessage{Key: bindingError.(*json.UnmarshalTypeError).Field, ErrorMessage: constants.ErrUnexpectedValue}

		logger.WithFields(logrus.Fields{
			constants.User:    bffSignInRequest.Username,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info(constants.ErrBindingFailed)

		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	if validatorError := validations.GetBFFValidator().Struct(&bffSignInRequest); validatorError != nil {
		validationErrors, _ := validations.FormatValidationErrors(validatorError)

		logger.WithFields(logrus.Fields{
			constants.User:    bffSignInRequest.Username,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info(constants.ErrValidationFailed)

		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}

	errorFromService := controller.service.SignInUser(ctx, ctx.Request.Context(), bffSignInRequest)

	if errorFromService != nil {
		if errorFromService == gorm.ErrRecordNotFound {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.User,
					ErrorMessage: constants.ErrRecordNotFOund,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				constants.User:    bffSignInRequest.Username,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info(constants.ErrAuthenticationFailed)

			ctx.IndentedJSON(http.StatusNotFound, errorResponse)
			return
		} else if strings.Contains(errorFromService.Error(), constants.ErrPasswordNotMatch) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Password,
					ErrorMessage: constants.ErrPasswordNotMatch,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				constants.User:    bffSignInRequest.Username,
				constants.Latency: time.Since(start).Milliseconds(),
			}).Info(constants.ErrPasswordNotMatch)


			ctx.IndentedJSON(http.StatusUnauthorized, errorResponse)
			return
		}

		errorResponse := genericModels.ErrorAPIResponse{
			Error: constants.ErrAuthenticationFailed,
		}

		logger.WithFields(logrus.Fields{
			constants.User:    bffSignInRequest.Username,
			constants.Latency: time.Since(start).Milliseconds(),
		}).Info(constants.ErrAuthenticationFailed)

		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	logger.WithFields(logrus.Fields{
		constants.User:    bffSignInRequest.Username,
		constants.Latency: time.Since(start).Milliseconds(),
	}).Info(constants.UserLoggedInSuccessMsg)

	ctx.IndentedJSON(http.StatusOK, constants.UserLoggedInSuccessMsg)
}
