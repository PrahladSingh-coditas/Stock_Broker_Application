package router

import (
	"authentication/middleware"
	genericConstants "stock_broker_application/src/constants"
	"stock_broker_application/src/utils"
	"watchlists/business"
	"watchlists/commons/constants"
	"watchlists/docs"
	"watchlists/handlers"
	"watchlists/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())

	docs.SwaggerInfo.Title = constants.SwaggerTitle

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(files.Handler))

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{genericConstants.AllowedOrigin},
		AllowMethods: []string{genericConstants.POST, genericConstants.GET},
		AllowHeaders: []string{genericConstants.Origin, genericConstants.ContentType, genericConstants.Authorization},
	}))

	db := utils.GetPostgresClient().GormDB
	rdb := utils.GetRedisClient()

	adgScripRepository := repository.NewadgStoWatchlistsRepository(db, rdb)
	adgScripService := business.NewadgStoWatchlistService(adgScripRepository, rdb)
	adgScripHandler := handlers.NewAdgStoWatchlistHandler(adgScripService)

	authGroup := router.Group(constants.AdgRoutePrefix)
	{
		authGroup.POST(constants.AdgScripToWatchlist, middleware.AuthMiddleware(), adgScripHandler.HandleAdgStoWatchlist)
	}

	return router
}
