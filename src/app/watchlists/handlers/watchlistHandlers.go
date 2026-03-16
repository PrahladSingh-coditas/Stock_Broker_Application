package handlers

import (
	"encoding/json"
	"net/http"
	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils/validations"
	"strings"
	"watchlists/business"
	"watchlists/commons/constants"
	"watchlists/models"

	"github.com/gin-gonic/gin"
)

type WatchlistsHandler struct {
	service *business.WatchlistsService
}

func NewWatchlistsHandler(service *business.WatchlistsService) *WatchlistsHandler {
	return &WatchlistsHandler{
		service: service,
	}
}

// Handles Watchlist Operations
// @Summary Executed ADG Operations
// @Description Performs ADD, GET or DELETE Operation Based on ActionType
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BFFAdgToWatchlistRequest true "Watchlist ADG Request"
// @Success 200 {object} models.BFFAdgToWatchlistResponse "ADG Request Successful"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 401 {object} models.ErrorAPIResponse "DataBase Error"
// @Failure 404 {object} models.ErrorAPIResponse "User does not exist"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /api/watch/watchlists [post]
func (controller *WatchlistsHandler) HandleWatchlistADG(ctx *gin.Context) {

	var warnings []string
	var bffWachlistRequest models.BFFAdgToWatchlistRequest

	if err := ctx.ShouldBind(&bffWachlistRequest); err != nil {

		errorMessage := genericModels.ErrorMessage{
			Key:          err.(*json.UnmarshalTypeError).Field,
			ErrorMessage: constants.ErrUnexpectedValue,
		}

		warnings = append(warnings, constants.ErrUnexpectedValue)

		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMessage,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	bffWachlistRequest.Action = models.ActionType(
		strings.ToUpper(string(bffWachlistRequest.Action)),
	)

	if !bffWachlistRequest.Action.IsValid() {
		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Error: constants.InvalidActionTypeError,
		})
		return
	}

	switch bffWachlistRequest.Action {

	case models.GET:
		if len(bffWachlistRequest.WatchlistIds) > 0 {
			ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
				Error: "watchlistIds should not be provided for GET action",
			})
			return
		}

	case models.ADD, models.DEL:
		if len(bffWachlistRequest.WatchlistIds) == 0 {
			ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
				Error: "watchlistIds required for ADD or DEL action",
			})
			return
		}
	}

	if err := validations.GetBFFValidator().Struct(&bffWachlistRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)
		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}

	usernameInterface, _ := ctx.Get(constants.FieldUsername)
	username := usernameInterface.(string)

	watchlists, err := controller.service.Watchlists(
		ctx,
		ctx.Request.Context(),
		bffWachlistRequest,
		username,
	)

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, genericModels.ErrorAPIResponse{
			Error: constants.ServerError,
		})
		return
	}

	var response models.BFFAdgToWatchlistResponse

	response.Status = constants.WatchlistSuccessMsg
	response.Action = bffWachlistRequest.Action
	response.WatchlistWithId = watchlists
	response.Warnings = []string{"No warnings!"}

	ctx.IndentedJSON(http.StatusOK, response)
}
