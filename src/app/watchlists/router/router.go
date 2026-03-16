package router

import (
	"watchlists/business"
	"watchlists/commons/constants"
	"watchlists/docs"
	"watchlists/handlers"
	"watchlists/middleware"
	"watchlists/repository"

	genericConstants "stock_broker_application/src/constants"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	docs.SwaggerInfo.Title = constants.SwaggerTitle

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(files.Handler))

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{genericConstants.AllowedOrigin},
		AllowMethods: []string{genericConstants.POST, genericConstants.GET},
		AllowHeaders: []string{genericConstants.Origin, genericConstants.ContentType, genericConstants.Authorization},
	}))

	watchlistsRepository := repository.NewWatchlistsRepository()
	watchlistsService := business.NewWatchlistsService(watchlistsRepository)
	watchlistsHandler := handlers.NewWatchlistsHandler(watchlistsService)

	authGroup := router.Group(constants.RoutePrefix)
	{
		authGroup.POST(constants.WatchlistADG, middleware.AuthMiddleware(), watchlistsHandler.HandleWatchlistADG)

	}

	return router
}
