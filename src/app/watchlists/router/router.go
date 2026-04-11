package router

import (
	genericConstants "stock_broker_application/src/constants"
	"watchlists/business"
	"watchlists/commons/constants"
	"watchlists/docs"
	"watchlists/handlers"
	"watchlists/middleware"
	"watchlists/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func GetRouter(db *gorm.DB, redisClient *redis.Client) *gin.Engine {
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

	adgScripRepository := repository.NewadgStoWatchlistsRepository(db, redisClient)
	adgScripService := business.NewadgStoWatchlistService(adgScripRepository, redisClient)
	adgScripHandler := handlers.NewAdgStoWatchlistHandler(adgScripService)

	authGroup := router.Group(constants.AdgRoutePrefix)
	{
		authGroup.POST(constants.AdgScripToWatchlist, middleware.AuthMiddleware(redisClient), adgScripHandler.HandleAdgStoWatchlist)
	}

	return router
}
