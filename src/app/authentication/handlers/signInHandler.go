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
type SigninUserHandler struct {
	service *business.SigninUserService
}

func NewSigninUserHandler(service *business.SigninUserService) *SigninUserHandler {
	return &SigninUserHandler{
		service: service,
	}
}

// HandleCreaterUser handles the user signin request.
// @Summary Sign in a user
// @Description Authenticates user credentials and returns success if valid
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.BFFSigninUserRequest true "User Signin Request"
// @Success 200 {string} string "User signed in successfully"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 401 {object} models.ErrorAPIResponse "Invalid username or password"
// @Failure 404 {object} models.ErrorAPIResponse "User Not Found"
// @Failure 500 {object} models.ErrorAPIResponse "Authentication failed"
// @Router /api/auth/signin [post]
func (controller *SigninUserHandler) HandleSigninUser(ctx *gin.Context) {
	var bffSigninUserRequest models.BFFSigninUserRequest
	//check for api binding errors
	if err := ctx.ShouldBind(&bffSigninUserRequest); err != nil {
		errorMsgs := genericModels.ErrorMessage{
			Key:          "request",
			ErrorMessage: constants.InvalidPayloadError,
		}
		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.InvalidPayloadError,
		})
		return
	}
	// valiation errors
	if err := validations.GetBFFValidator().Struct(&bffSigninUserRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)
		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}
	// then we call the function that is in service
	err := controller.service.SigninUser(ctx, ctx.Request.Context(), bffSigninUserRequest)
	if err != nil {
		if err.Error() == constants.UserNotFoundError {
			ctx.JSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Error: constants.UserNotFoundError,
			})
			return
		}
		ctx.JSON(http.StatusUnauthorized, genericModels.ErrorAPIResponse{
			Error: constants.AuthenticationFailedError,
		})
		return
	}
	ctx.IndentedJSON(http.StatusOK, constants.UserLoggedInSuccessMsg)
}
