package handlers

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/models"
	"encoding/json"
	"net/http"
	"strings"

	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils/validations"

	"github.com/gin-gonic/gin"
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

	var bffSigninUserRequet models.BFFSigninUserRequest

	if err := ctx.ShouldBind(&bffSigninUserRequet); err != nil {
		errorMsgs := genericModels.ErrorMessage{
			Key:          err.(*json.UnmarshalTypeError).Field,
			ErrorMessage: constants.ErrUnexpectedValue,
		}
		ctx.JSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	// 400 error
	if err := validations.GetBFFValidator().Struct(&bffSigninUserRequet); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)
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

		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return

	}

	ctx.IndentedJSON(http.StatusOK, constants.UserLoggedInSuccessMsg)
}
