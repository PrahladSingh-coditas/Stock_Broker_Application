package handlers

import (
	"encoding/json"
	"net/http"
	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils/validations"
	"strings"
	"watchlists/business"
	"watchlists/commons"
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

	var req models.BFFAdgToWatchlistRequest

	if err := ctx.ShouldBind(&req); err != nil {

		var field string
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			field = ute.Field
		}
		commons.ErrorResponse(field, constants.ErrUnexpectedValue, constants.ErrInvalidPayload, string(req.Action), ctx)
		return
	}

	req.Action = models.ActionType(strings.ToUpper(string(req.Action)))
	req.ScripId = strings.ToUpper(req.ScripId)

	if err := validations.GetBFFValidator().Struct(&req); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)
		ctx.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	usernameInterface, exists := ctx.Get(constants.FieldUsername)
	if !exists {
		commons.ErrorResponse(constants.FieldUsername, constants.UserNotFoundError, constants.UserNotFoundError, string(req.Action), ctx)
		return
	}
	username := usernameInterface.(string)

	watchlists, warnings, err := controller.service.Watchlists(ctx, ctx.Request.Context(), req, username)

	if err != nil {

		switch err.Error() {

		case constants.NoValidWatchlistIdsError:
			commons.ErrorResponse(constants.FieldWatchlistId, constants.NoValidWatchlistIdsError, constants.RequestFailedError, string(req.Action), ctx)

		case constants.UserNotFoundError:
			commons.ErrorResponse(constants.FieldUsername, constants.UserNotFoundError, constants.RequestFailedError, string(req.Action), ctx)

		case constants.ScripIdNotFoundError:
			commons.ErrorResponse(constants.FieldScripId, constants.ScripIdNotFoundError, constants.RequestFailedError, string(req.Action), ctx)

		case constants.WatchlistNotFoundError:
			commons.ErrorResponse(constants.FieldWatchlistId, constants.WatchlistNotFoundError, constants.RequestFailedError, string(req.Action), ctx)

		case constants.QueryError:
			commons.ErrorResponse(constants.Database, constants.QueryError, constants.RequestFailedError, string(req.Action), ctx)

		default:
			ctx.JSON(http.StatusInternalServerError, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{},
				Error:   err.Error(),
			})
		}
		return
	}

	response := models.BFFAdgToWatchlistResponse{
		Status:          constants.WatchlistSuccessMsg,
		Action:          req.Action,
		WatchlistWithId: watchlists,
	}

	if len(watchlists) > 0 {
		response.WatchlistWithId = watchlists
	}

	if len(warnings) > 0 {
		response.Warnings = warnings
	}

	ctx.JSON(http.StatusOK, response)
}
