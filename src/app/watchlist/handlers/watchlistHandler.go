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
// @Failure 401 {object} models.ErrorAPIResponse "Pass mismatch"
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

		// 400 error
		if strings.Contains(errorString, constants.ErrEmptyScripId) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.ScripId,
					ErrorMessage: constants.ErrEmptyScripId,
				},
				Error: constants.ErrInvalidPayload,
			}

			logger.WithFields(logrus.Fields{
				"user":    username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrEmptyScripId)

			ctx.IndentedJSON(http.StatusBadRequest, errorResponse)
			return
		}
		// 400 error if watchlist ids empty
		if strings.Contains(errorString, constants.ErrEmptyWatchlists) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.Watchlist,
					ErrorMessage: constants.ErrEmptyWatchlists,
				},
				Error: constants.ErrInvalidPayload,
			}

			logger.WithFields(logrus.Fields{
				"user":    username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrEmptyWatchlists)

			ctx.IndentedJSON(http.StatusBadRequest, errorResponse)
			return
		}
		//400 error if invalid action
		if strings.Contains(errorString, constants.ErrInvalidAction) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.Action,
					ErrorMessage: constants.ErrInvalidAction,
				},
				Error: constants.ErrInvalidPayload,
			}

			logger.WithFields(logrus.Fields{
				"user":    username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrInvalidAction)

			ctx.IndentedJSON(http.StatusBadRequest, errorResponse)
			return

		}

		//404 if wrong watchlist entered
		if strings.Contains(errorString, constants.ErrWatchlistNotFound) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.ScripId,
					ErrorMessage: constants.ErrWatchlistNotFound,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				"user":    username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrWatchlistNotFound)

			ctx.IndentedJSON(http.StatusNotFound, errorResponse)
			return
		}
		
		//404 if any query error
		if strings.Contains(errorString, constants.ErrNoRowsAffected) {
			errorResponse := genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{
					Key:          commons.Watchlist,
					ErrorMessage: constants.ErrNoRowsAffected,
				},
				Error: constants.ErrAuthenticationFailed,
			}

			logger.WithFields(logrus.Fields{
				"user":    username,
				"latency": time.Since(start).Milliseconds(),
			}).Info(constants.ErrWatchlistNotFound)

			ctx.IndentedJSON(http.StatusNotFound, errorResponse)
			return
		}

		
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := models.BFFAdgToWatchlistResponse{
		Status:          "Used ActionType Successfully",
		Action:          bffAdgToWatchlistRequest.Action,
		WatchlistWithId: watchlistNameWithId,
	}

	if len(watchlistNameWithId) == 0 {
		response.Status = "No watchlists updated"
	}

	if len(watchlistNameWithId) > 0 {
		response.WatchlistWithId = watchlistNameWithId
	}

	if len(warnings) > 0 {
		response.Warnings = warnings
	}

	ctx.JSON(http.StatusOK, response)
}
