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

type SignInUserHandler struct {
	service *business.SignInService
}

func NewSignInUserHandler(service *business.SignInService) *SignInUserHandler {
	return &SignInUserHandler{
		service: service,
	}
}

// HandleSignInUser handles the signin request.
// @Summary User Sign-in
// @Description Authenticates user and returns message
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.BFFSignInRequest true "User Sign-in Request"
// @Success 200 {object} models.BFFSignInResponse "Signin successful"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 401 {object} models.ErrorAPIResponse "Invalid credentials"
// @Failure 404 {object} models.ErrorAPIResponse "User Not Found"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /api/auth/signin [post]
func (controller *SignInUserHandler) HandleSignInUser(ctx *gin.Context) {
	var bffSignInRequest models.BFFSignInRequest

	// Bind JSON payload
	if err := ctx.ShouldBind(&bffSignInRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Error: constants.ErrInvalidPayload,
		})
		return
	}

	if err := validations.GetBFFValidator().Struct(&bffSignInRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	response, err := controller.service.SignIn(ctx.Request.Context(), bffSignInRequest.Username, bffSignInRequest.Password)
	if err != nil {
		if err.Error() == constants.ErrUsernameNotFound {
			ctx.JSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Error: constants.ErrUsernameNotFound,
			})
			return
		}

		if err.Error() == constants.ErrPasswordMismatch {
			ctx.JSON(http.StatusUnauthorized, genericModels.ErrorAPIResponse{
				Error: constants.ErrInvalidCredentials,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, genericModels.ErrorAPIResponse{
			Error: constants.ErrLoginFailed,
		})
		return
	}
	ctx.JSON(http.StatusOK, response)

}
