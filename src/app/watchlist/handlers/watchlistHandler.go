package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"stock_broker_application/src/app/watchlist/business"
	"stock_broker_application/src/app/watchlist/commons/constants"
	"stock_broker_application/src/app/watchlist/models"
	genericModels "stock_broker_application/src/models"
	"stock_broker_application/src/utils/validations"
	"strings"
	"time"

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

// HandlerWatchlist handles the add, get, del of watchlist
// @Summary Add, Del, Get of watchlist
// @Description Handles get, add, delete request for watchlists
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BFFAdgToWatchlistRequest true "ADD, GET, DEL request of watchlist"
// @Success 200 {object} models.BFFAdgToWatchlistResponse "action completed successfully"
// @Success 201 {object} models.BFFAdgToWatchlistResponse "Data added successfully"
// @Failure 400 {object} models.ErrorAPIResponse "Invalid input payload"
// @Failure 404 {object} models.ErrorAPIResponse "User does not exist"
// @Failure 500 {object} models.ErrorAPIResponse "Internal Server Error"
// @Router /api/watch/watchlist [post]
func (controller *WatchlistHandler) HandleWatchlist(ctx *gin.Context) {
	start := time.Now()
	logger := logrus.New()

	var bffAdgToWatchlistRequest models.BFFAdgToWatchlistRequest

	username := ctx.GetString(constants.Username)

	if err := ctx.ShouldBind(&bffAdgToWatchlistRequest); err != nil {
		errorMsgs := genericModels.ErrorMessage{Key: err.(*json.UnmarshalTypeError).Field, ErrorMessage: constants.ErrUnexpectedValue}

		logger.WithFields(logrus.Fields{
			constants.Username: username,
			constants.Latency:  time.Since(start).Milliseconds(),
		}).Error(constants.ErrBindingFailed)

		ctx.IndentedJSON(http.StatusBadRequest, genericModels.ErrorAPIResponse{
			Message: errorMsgs,
			Error:   constants.ErrInvalidPayload,
		})
		return
	}

	if err := validations.GetBFFValidator().Struct(&bffAdgToWatchlistRequest); err != nil {
		validationErrors, _ := validations.FormatValidationErrors(err)

		logger.WithFields(logrus.Fields{
			constants.Username: username,
			constants.Latency:  time.Since(start).Milliseconds(),
		}).Error(constants.ErrValidationFailed)

		ctx.IndentedJSON(http.StatusBadRequest, validationErrors)
		return
	}

	bffAdgToWatchlistRequest.Action = models.ActionType(strings.ToUpper(string(bffAdgToWatchlistRequest.Action)))

	watchlistWithId, warnings, err := controller.service.ServiceWatchlist(ctx, ctx.Request.Context(), logger, bffAdgToWatchlistRequest, username)

	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), errors.New(constants.ErrUserNotFoundMsg).Error()) {

			logger.Error(constants.ErrUserNotFoundMsg)

			ctx.IndentedJSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{Key: constants.User, ErrorMessage: constants.ErrUserNotFoundMsg},
				Error:   constants.ErrRequestFailed,
			})
			return
		} else if strings.Contains(err.Error(), errors.New(constants.ErrDatabaseQueryErrorMsg).Error()) {

			logger.Error(constants.ErrDatabaseQueryErrorMsg)

			ctx.IndentedJSON(http.StatusInternalServerError, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{Key: constants.Database, ErrorMessage: constants.ErrDatabaseQueryErrorMsg},
				Error:   constants.ErrRequestFailed,
			})
			return
		} else if strings.Contains(err.Error(), errors.New(constants.ErrScripNotFoundMsg).Error()) {

			logger.Error(constants.ErrScripNotFoundMsg)
			ctx.IndentedJSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{Key: constants.FieldScripId, ErrorMessage: constants.ErrScripNotFoundMsg},
				Error:   constants.ErrRequestFailed,
			})
			return
		} else if strings.Contains(err.Error(), errors.New(constants.ErrNoWatchlistForScripMsg).Error()) {

			logger.Error(constants.ErrNoWatchlistForScripMsg)

			ctx.IndentedJSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{Key: constants.FieldWatchlistId, ErrorMessage: constants.ErrNoWatchlistForScripMsg},
				Error:   constants.ErrRequestFailed,
			})
			return
		} else if strings.Contains(err.Error(), errors.New(constants.EmptyWatchlistIdsError).Error()) {

			logger.Error(constants.EmptyWatchlistIdsError)

			ctx.IndentedJSON(http.StatusNotFound, genericModels.ErrorAPIResponse{
				Message: genericModels.ErrorMessage{Key: constants.FieldWatchlistId, ErrorMessage: constants.EmptyWatchlistIdsError},
				Error:   constants.ErrRequestFailed,
			})
			return
		}

		logger.Error(constants.ErrInternalServer)

		ctx.IndentedJSON(http.StatusInternalServerError, genericModels.ErrorAPIResponse{
			Message: genericModels.ErrorMessage{Key: constants.Server, ErrorMessage: constants.ErrInternalServer},
			Error:   constants.ErrRequestFailed,
		})
		return
	}

	bffAdgToWatchlistResponse := models.BFFAdgToWatchlistResponse{
		Status:          constants.Successful,
		Action:          bffAdgToWatchlistRequest.Action,
		WatchlistWithId: watchlistWithId,
		Warnings:        warnings,
	}

	logger.Info(constants.ActionCompletedMsg)

	if bffAdgToWatchlistRequest.Action == models.ADD {
		ctx.IndentedJSON(http.StatusCreated, bffAdgToWatchlistResponse)
	} else {
		ctx.IndentedJSON(http.StatusOK, bffAdgToWatchlistResponse)
	}
}
