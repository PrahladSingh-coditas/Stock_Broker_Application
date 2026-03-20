package router

import (
	"stock_broker_application/src/app/watchlist/business"
	"stock_broker_application/src/app/watchlist/commons/constants"
	"stock_broker_application/src/app/watchlist/handlers"
	"stock_broker_application/src/app/watchlist/middleware"
	"stock_broker_application/src/app/watchlist/repository"
	genericConstants "stock_broker_application/src/constants"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"stock_broker_application/src/app/watchlist/docs"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	//router.Use(middleware.AuthMiddleware
	router.Use(gin.Recovery())

	docs.SwaggerInfo.Title = constants.SwaggerTitle

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(files.Handler))

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{genericConstants.AllowedOrigin},
		AllowMethods: []string{genericConstants.POST, genericConstants.GET},
		AllowHeaders: []string{genericConstants.Origin, genericConstants.ContentType, genericConstants.Authorization},
	}))

	watchlistRepository := repository.NewWatchlistRepository()
	watchlistService := business.NewWatchlistService(watchlistRepository)
	watchlistHandler := handlers.NewWatchlistHandler(watchlistService)

	authGroup := router.Group(constants.WatchRoutePrefix)
	{
		authGroup.POST(constants.Watchlist, middleware.WatchlistMiddleware(), watchlistHandler.HandleWatchlist)
	}

	return router
}