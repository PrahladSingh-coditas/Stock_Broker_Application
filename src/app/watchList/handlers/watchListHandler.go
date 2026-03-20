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
// @Success 200 {object} models.BFFAdgToWatchlistResponse "GET or DELETE success"
// @Success 201 {object} models.BFFAdgToWatchlistResponse "ADD success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 409 {object} map[string]string "Conflict"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/watchlist/watchlist-api [post]
func (controller *WatchListHandler) HandleWatchList(c *gin.Context) {
	var req models.BFFAdgToWatchListRequest

	username := c.MustGet("username").(string)

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	response, err := controller.service.HandleWatchlist(c.Request.Context(), username, req)

	if err != nil {

		switch err.Error() {

		case "scrip not found":
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})

		case "watchlist limit exceeded":
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})

		case "watchlist not owned by user":
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}

		return
	}

	if req.Action == models.ADD {
		c.JSON(http.StatusCreated, response)
		return
	}

	c.JSON(200, response)

}
