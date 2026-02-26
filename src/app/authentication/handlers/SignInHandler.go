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

type SignInHandler struct {
	service *business.SignInService
}

func NewSignInHandler(service *business.SignInService) *SignInHandler {
	return &SignInHandler{
		service: service,
	}
}

// HandlerSignIn handles the user signin request.
// @Summary Sign in an existing user
// @Description Authenticates user and returns JWT token
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.BFFSignInRequest true "User Sign In Request"
// @Success 200 {object} models.BFFSignInResponse "Signin successful"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 401 {object} models.ErrorAPIResponse "Invalid credentials"
// @Failure 404 {object} models.ErrorAPIResponse "User does not exist"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /api/auth/signin [post]
func (controller *SignInHandler) HandleSignIn(ctx *gin.Context) {
	var bffSignInRequest models.BFFSignInRequest

	if errBindReq := ctx.ShouldBind(&bffSignInRequest); errBindReq != nil {
		errorMsgs := genericModels.ErrorMessage{Key: errBindReq.(*json.UnmarshalTypeError).Field, ErrorMessage: constants.ErrUnexpectedValue}

		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	if errValidation := validations.GetBFFValidator().Struct(&bffSignInRequest); errValidation != nil {
		validationErros, _ := validations.FormatValidationErrors(errValidation)
		ctx.IndentedJSON(http.StatusBadRequest, validationErros)
		return
	}

	errWhileSignIn := controller.service.SignIn(ctx, ctx.Request.Context(), bffSignInRequest)
	if errWhileSignIn != nil {
		if errors.Is(errWhileSignIn, errors.New(constants.ErrUserNotFound)) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.Username,
					ErrorMessage: constants.ErrUserNotFound,
				},
				Error: constants.ErrAuthenticationFailed,
			}
			ctx.IndentedJSON(http.StatusNotFound, errorResponse)
			return
		}

		if errors.Is(errWhileSignIn, errors.New(constants.ErrIncorrectPassword)) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.Password,
					ErrorMessage: constants.ErrIncorrectPassword,
				},
				Error: constants.ErrAuthenticationFailed,
			}
			ctx.IndentedJSON(http.StatusUnauthorized, errorResponse)
			return
		}
		ctx.IndentedJSON(http.StatusUnauthorized, genericModels.ErrorAPIResponse{
			Error: constants.ErrSignInFailed,
		})
		return
	}
	ctx.IndentedJSON(http.StatusOK, constants.UserLoggedInSuccessMsg)
}
