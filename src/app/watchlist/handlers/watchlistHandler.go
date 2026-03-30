package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"watchlist/business"
	"watchlist/commons"
	"watchlist/commons/constants"
	"watchlist/models"

	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils/validations"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WatchlistHandler struct {
	service *business.WatchlistService
}

func NewWatchlistHandler(service *business.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{
		service: service,
	}
}

// Handle Watchlist ADG Operations
// @Summary Watchlist ADG
// @Description Performs Add,Delete,Get to watchlists
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BFFAdgToWatchlistRequest true "ADG to Watchlist Request"
// @Success 200 {object} models.BFFAdgToWatchlistResponse "ADG performed successfully"
// @Failure 400 {object} models.ErrorAPIResponse  "Invalid Input Payload"
// @Failure 404 {object} models.ErrorAPIResponse "Watchlists not found"
// @Router /api/watchlist/watchlistADG [post]
func (controller *WatchlistHandler) HandleWatchlist(ctx *gin.Context) {
	start := time.Now()
	logger := logrus.New()

	var bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest

	if err := ctx.ShouldBind(&bffAdgToWatchlistRequest); err != nil {
		errorMsgs := genericModels.ErrorMessage{
			Key:          err.(*json.UnmarshalTypeError).Field,
			ErrorMessage: constants.ErrUnexpectedValue,
		}

		logger.WithFields(logrus.Fields{
			"password": bffAdgToWatchlistRequest.Action,
			"latency":  time.Since(start).Milliseconds(),
		}).Info(constants.ErrBinding)

		ctx.JSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	bffAdgToWatchlistRequest.Action = models.ActionType(strings.ToUpper(string(bffAdgToWatchlistRequest.Action)))
	bffAdgToWatchlistRequest.ScripId = strings.ToUpper(bffAdgToWatchlistRequest.ScripId)

	// 400 error
	if err := validations.GetBFFValidator().Struct(&bffAdgToWatchlistRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)

		logger.WithFields(logrus.Fields{
			"action":  bffAdgToWatchlistRequest.Action,
			"latency": time.Since(start).Milliseconds(),
		}).Info(constants.ErrUnexpectedValue)

		ctx.JSON(http.StatusBadRequest, validationErrors)
		return
	}

	username := ctx.GetString("username")

	watchlistNameWithId, warnings, err := controller.service.Watchlist(ctx, ctx.Request.Context(), username, bffAdgToWatchlistRequest)
	if err != nil {
		errorString := err.Error()

		//400 error	
		// 400 error if watchlist ids empty		
		//400 error if invalid action
		
		if strings.Contains(errorString, constants.ErrEmptyScripId) || strings.Contains(errorString, constants.ErrEmptyWatchlists) || strings.Contains(errorString, constants.ErrInvalidAction) {
			var key, errorMsg string

			switch {
			case strings.Contains(errorString, constants.ErrEmptyScripId):
				key = commons.ScripId
				errorMsg = constants.ErrEmptyScripId
			case strings.Contains(errorString, constants.ErrEmptyWatchlists):
				key = commons.Watchlist
				errorMsg = constants.ErrEmptyWatchlists
			case strings.Contains(errorString, constants.ErrInvalidAction):
				key = commons.Action
				errorMsg = constants.ErrInvalidAction
			}

			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          key,
					ErrorMessage: errorMsg,
				},
				Error: constants.ErrInvalidPayload,
			}

			logger.WithFields(logrus.Fields{
				"user":    username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(errorString)

			ctx.IndentedJSON(http.StatusBadRequest, errorResponse)
			return
		}

		//404 if wrong watchlist entered
		//404 if no valid watchlists
		//404 if all watchlists not of user
		//404 if user not found
		//404 if any query error
		
		if strings.Contains(errorString, constants.ErrWatchlistNotFound) ||
			strings.Contains(errorString, constants.ErrNoValidWatchlists) ||
			strings.Contains(errorString, constants.ErrWatchlistsNotOfUser) ||
			strings.Contains(errorString, constants.ErrUserNotFound) ||
			strings.Contains(errorString, constants.ErrNoRowsAffected) {

			var key, errorMsg string

			switch {
			case strings.Contains(errorString, constants.ErrWatchlistNotFound):
				key = commons.ScripId
				errorMsg = constants.ErrWatchlistNotFound
			case strings.Contains(errorString, constants.ErrNoValidWatchlists):
				key = commons.Watchlist
				errorMsg = constants.ErrNoValidWatchlists
			case strings.Contains(errorString, constants.ErrWatchlistsNotOfUser):
				key = commons.Watchlist
				errorMsg = constants.ErrWatchlistsNotOfUser
			case strings.Contains(errorString, constants.ErrUserNotFound):
				key = commons.Username
				errorMsg = constants.ErrUserNotFound
			case strings.Contains(errorString, constants.ErrNoRowsAffected):
				key = commons.Watchlist
				errorMsg = constants.ErrNoRowsAffected
			}

			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          key,
					ErrorMessage: errorMsg,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				"user":    username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(errorMsg)

			ctx.IndentedJSON(http.StatusNotFound, errorResponse)
			return
		}

		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var response models.BFFAdgToWatchlistResponse

	if len(watchlistNameWithId) == 0 {
		response.Status = constants.ActionTypeFailure
	}

	if len(watchlistNameWithId) > 0 {
		response.Status = constants.ActionTypeSuccess
		response.WatchlistWithId = watchlistNameWithId
	}

	if len(warnings) > 0 {
		response.Warnings = warnings
	}
	response.Action = bffAdgToWatchlistRequest.Action
		

	ctx.JSON(http.StatusOK, response)
}
