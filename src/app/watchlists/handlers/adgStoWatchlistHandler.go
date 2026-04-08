package handlers

import (
	"errors"
	"fmt"
	"strings"
	"watchlists/commons/constants"

	"encoding/json"
	"net/http"
	"stock_broker_application/src/utils/validations"
	"watchlists/business"
	"watchlists/commons"
	"watchlists/models"

	genericModels "stock_broker_application/src/models"

	"github.com/gin-gonic/gin"
)

type AdgStoWatchlistHandler struct {
	AdgStoWatchlistService *business.AdgStoWatchlistService
}

func NewAdgStoWatchlistHandler(adgStoWatchlistService *business.AdgStoWatchlistService) *AdgStoWatchlistHandler {
	return &AdgStoWatchlistHandler{
		AdgStoWatchlistService: adgStoWatchlistService,
	}
}

// this fucntion handles user requests for add, delete and get scrips to/ from a list of watchlists
// Handles ADG scrip to watchlist functionality
// @Summary Perform ADG for watchlist
// @Description verifies the JWT and adds/ deletes/ gets the scrip in the request to/ from the list of watchlists
// @Tags Watchlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BffAdgStoWatchlistRequest true "ADG Watchlist Request"
// @Success 200 {object} models.BffAdgStoWatchlistResponse "ADG Performed successfully"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 401 {object} models.ErrorAPIResponse "Case A: JWT expired / Case B: JWT invalid / Case C: Unauthorized request"
// @Failure 404 {object} models.ErrorAPIResponse "User does not exist"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /adg/scrip [post]
func (controller *AdgStoWatchlistHandler) HandleAdgStoWatchlist(ctx *gin.Context) {
	var bffAdgStoWatchlistRequest models.BffAdgStoWatchlistRequest

	if err := ctx.ShouldBind(&bffAdgStoWatchlistRequest); err != nil {
		errorMsg := genericModels.ErrorMessage{
			Key:          err.(*json.UnmarshalTypeError).Field,
			ErrorMessage: constants.ErrUnexpectedValue,
		}
		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsg,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	err := validations.GetBFFValidator().Struct(&bffAdgStoWatchlistRequest)
	if err != nil {
		validationError, _ := validations.FormatValidationErrors(err)
		ctx.IndentedJSON(http.StatusBadRequest, validationError)
		return
	}

	username := ctx.GetString(commons.Username)
	ReqAction := models.Actiontype(strings.ToUpper(string(bffAdgStoWatchlistRequest.Action)))

	warnings, respWatchlistsWithIds, err := controller.AdgStoWatchlistService.AdgStoWatchlist(ctx, username, bffAdgStoWatchlistRequest)
	if err != nil {
		fmt.Println("ERROR:", err)

		switch {
		case errors.Is(err, constants.UserNotFoundError):
			ctx.IndentedJSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.UserID,
					ErrorMessage: constants.ErrUserNotFound,
				},
				Error: fmt.Sprintf(constants.ErrRequestFailed, ReqAction),
			})
			return

		case errors.Is(err, constants.DatabaseQueryError):
			ctx.IndentedJSON(http.StatusInternalServerError, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Database,
					ErrorMessage: constants.ErrDatabaseQuery,
				},
				Error: fmt.Sprintf(constants.ErrRequestFailed, ReqAction),
			})
			return

		case errors.Is(err, constants.ScripDoesNotExistError):
			ctx.IndentedJSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.ScripID,
					ErrorMessage: constants.ErrScripDoesnotExist,
				},
				Error: fmt.Sprintf(constants.ErrRequestFailed, ReqAction),
			})
			return

		case errors.Is(err, constants.ScripNotInWatchlistsError):
			ctx.IndentedJSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.ScripID,
					ErrorMessage: constants.ErrScripNotInWatchlists,
				},
				Error: fmt.Sprintf(constants.ErrRequestFailed, ReqAction),
			})
			return

		case errors.Is(err, constants.InvalidWatchlistsError):
			ctx.IndentedJSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.WatchlistID,
					ErrorMessage: constants.ErrInvalidWatchists,
				},
				Error: fmt.Sprintf(constants.ErrRequestFailed, ReqAction),
			})
			return

		case errors.Is(err, constants.InvalidActionTypeError):
			ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Action,
					ErrorMessage: constants.ErrInvalidActiontype,
				},
				Error: fmt.Sprintf(constants.ErrRequestFailed, ReqAction),
			})
			return

		// case errors.Is(err, constants.ScripNotAddedToAnyWatchlistsError):
		// 	ctx.IndentedJSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
		// 		Message: genericModels.ErrorMessage{
		// 			Key:          constants.ScripID,
		// 			ErrorMessage: constants.ErrScripNotAddedToAnyWatchlists,
		// 		},
		// 		Error: fmt.Sprintf(constants.ErrRequestFailed, ReqAction),
		// 	})
		// 	return

		default:
			ctx.IndentedJSON(http.StatusInternalServerError, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          constants.Server,
					ErrorMessage: constants.InternalServerError,
				},
				Error: fmt.Sprintf(constants.ErrRequestFailed, ReqAction),
			})
			return
		}
	}
	ctx.IndentedJSON(http.StatusOK, models.BffAdgStoWatchlistResponse{
		Status:          "successful",
		Action:          ReqAction,
		WatchlistWithId: respWatchlistsWithIds,
		Warnings:        warnings,
	})
}
