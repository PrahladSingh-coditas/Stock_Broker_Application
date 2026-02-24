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

	if err := ctx.ShouldBind(&bffSigninUserRequest); err != nil {
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

	if err := validations.GetBFFValidator().Struct(&bffSigninUserRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)
		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}

	err := controller.service.SigninUser(ctx, ctx.Request.Context(), bffSigninUserRequest)
	if err != nil {

		if err.Error() == constants.ErrUserNotFound {
			ctx.JSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Error: constants.ErrUserNotFound,
			})
			return
		}

		ctx.JSON(http.StatusUnauthorized, genericModels.ErrorAPIResponse{
			Error: constants.ErrAuthenticationFailed,
		})
		return
	}
	ctx.IndentedJSON(http.StatusOK, constants.UserLoggedInSuccessMsg)
}
