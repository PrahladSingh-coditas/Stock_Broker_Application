package router

import (
	"watchlist/business"
	"watchlist/commons/constants"
	"watchlist/docs"
	"watchlist/handlers"
	"watchlist/middleware"
	"watchlist/repository"

	genericConstants "stock_broker_application/src/constants"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter() *gin.Engine { // it basically gives a gin engine
	router := gin.New()
	//router.Use(middleware.AuthMiddleware())
	router.Use(gin.Recovery())

	docs.SwaggerInfo.Title = constants.SwaggerTitle //prints title

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(files.Handler))

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{genericConstants.AllowedOrigin},                                                        // it has * therefore from any port we cna call;i.e.terminal,swagger,postman etc etc
		AllowMethods: []string{genericConstants.POST, genericConstants.GET},                                           // Only these HTTP methods allowed.
		AllowHeaders: []string{genericConstants.Origin, genericConstants.ContentType, genericConstants.Authorization}, // This allows frontend to send these headers.
	}))

	watchlistsRepository := repository.NewWatchlistRepository() // remove new
	watchlistsService := business.NewWatchlistService(watchlistsRepository)
	watchlistsHandler := handlers.NewWatchlistHandler(watchlistsService)

	authGroup := router.Group(constants.RoutePrefix)
	{
		authGroup.POST(constants.WatchlistADG, middleware.AuthMiddleware(), watchlistsHandler.HandleWatchlist)
	}

	return router
}
