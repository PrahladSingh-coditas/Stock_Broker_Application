package router

import (
	"watchList/buisness"
	"watchList/commons/constants"
	"watchList/handlers"
	"watchList/middleware"
	"watchList/repository"

	"watchList/docs"

	genericConstants "stock_broker_application/src/constants"
	"stock_broker_application/src/utils"

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

	db := utils.GetPostgresClient().GormDB
	watchListRepository := repository.NewWatchlistRepository(db)
	watchListService := buisness.NewWatchlistService(watchListRepository)
	watchListHandler := handlers.NewWatchListHandler(watchListService)

	authGroup := router.Group(constants.AuthRoutePrefix)
	{
		authGroup.POST(constants.WatchList, middleware.AuthMiddleware(), watchListHandler.HandleWatchList)
	}

	return router
}
