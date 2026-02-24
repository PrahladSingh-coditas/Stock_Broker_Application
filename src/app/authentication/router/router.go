package router

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/docs"
	"authentication/handlers"
	"authentication/middleware"
	"authentication/repository"

	genericConstants "stock_broker_application/src/constants"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.AuthMiddleware())
	router.Use(gin.Recovery())

	docs.SwaggerInfo.Title = constants.SwaggerTitle

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(files.Handler))

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{genericConstants.AllowedOrigin},
		AllowMethods: []string{genericConstants.POST, genericConstants.GET},
		AllowHeaders: []string{genericConstants.Origin, genericConstants.ContentType, genericConstants.Authorization},
	}))

	createUserRepository := repository.NewCreateUserRepository()
	createUserService := business.NewCreateUserService(createUserRepository)
	createUserHandler := handlers.NewCreateUserHandler(createUserService)

	signInUserRepository:=repository.NewSignInUserRepository()
	signInUserService:=business.NewSignInUserService(signInUserRepository)
	signInUserHandler:=handlers.NewSignInUserHandler(signInUserService)

	
	authGroup := router.Group(constants.AuthRoutePrefix)
	{
		authGroup.POST(constants.Signup, createUserHandler.HandleCreaterUser)
		authGroup.POST(constants.Signin, signInUserHandler.HandleSignInUser)
	}

	return router
}
