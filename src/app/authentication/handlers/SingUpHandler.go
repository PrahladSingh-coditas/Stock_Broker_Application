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

	"github.com/gin-gonic/gin"
)

type CreaterUserHandler struct {
	service *business.CreateUserService
}

func NewCreateUserHandler(service *business.CreateUserService) *CreaterUserHandler {
	return &CreaterUserHandler{
		service: service,
	}
}

// HandlerCreaterUser handles the user creation request.
// @Summary Create a new user
// @Description Handles user registration by validating input and storing user details
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.BFFCreateUserRequest true "User Registration Request"
// @Success 201 {string} string "User created successfully"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 409 {object} models.ErrorAPIResponse "User already exists"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /api/auth/signup [post]
func (controller *CreaterUserHandler) HandleCreaterUser(ctx *gin.Context) {

	var bffCreateUserRequest models.BFFCreateUserRequest

	if err := ctx.ShouldBind(&bffCreateUserRequest); err != nil {
		errorMsgs := genericModels.ErrorMessage{Key: err.(*json.UnmarshalTypeError).Field, ErrorMessage: constants.ErrUnexpectedValue}
		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	if err := validations.GetBFFValidator().Struct(&bffCreateUserRequest); err != nil {
		validationErros, _ := validations.FormatValidationErrors(err)
		ctx.IndentedJSON(http.StatusBadRequest, validationErros)
		return
	}

	err := controller.service.CreateNewUser(ctx, ctx.Request.Context(), bffCreateUserRequest)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrDuplicateEntry) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          strings.Split(err.Error(), constants.ErrDuplicateEntry)[0],
					ErrorMessage: constants.ErrUserAlreadyExists,
				},
				Error: constants.ErrConflict,
			}
			ctx.IndentedJSON(http.StatusConflict, errorResponse)
			return
		}
		errorResponse := genericModels.ErrorAPIResponse{
			Error: constants.ErrUserCreationFailed,
		}
		ctx.IndentedJSON(http.StatusInternalServerError, errorResponse)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, constants.UserCreationSuccessMsg)
}
