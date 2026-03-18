package handlers

import (
	"net/http"
	"watchList/buisness"
	"watchList/models"

	"github.com/gin-gonic/gin"
)

type WatchListHandler struct {
	service buisness.WatchlistService
}

func NewWatchListHandler(service buisness.WatchlistService) *WatchListHandler {
	return &WatchListHandler{
		service: service,
	}
}

// HandleWatchlist godoc
// @Summary Manage watchlist scrips
// @Description Add, Delete or Get scrips in watchlists
// @Tags Watchlist
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body models.BFFAdgToWatchListRequest true "Watchlist Request"
// @Success 200 {object} models.BFFAdgToWatchlistResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/watchlist/watchlist-api [post]
func (controller *WatchListHandler) HandleWatchList(c *gin.Context) {
	var req models.BFFAdgToWatchListRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	usernameInterface, exists := c.Get("username")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "user not authorized",
		})
		return
	}

	username := usernameInterface.(string)

	response, err := controller.service.HandleWatchlist(username, req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, response)

}
