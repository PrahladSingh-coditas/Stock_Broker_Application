package router

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/docs"
	"authentication/handlers"
	"authentication/middleware"
	"authentication/repository"

	genericConstants "stock_broker_application/src/constants"
	"stock_broker_application/src/utils"

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

	createUserRepository := repository.NewCreateUserRepository()
	createUserService := business.NewCreateUserService(createUserRepository)
	createUserHandler := handlers.NewCreateUserHandler(createUserService)

	signInRepository := repository.NewSignInRepository()
	signInService := business.NewSignInService(signInRepository)
	signInHandler := handlers.NewSignInHandler(signInService)

	postgresClient := utils.GetPostgresClient().GormDB
	verifyUserOtpRepository := repository.NewValidateUserOtpRepository()
	verifyUserOtpService := business.NewValidateUserOtpService(verifyUserOtpRepository, postgresClient)
	verifyUserOtpHandler := handlers.NewValidateUserOtpHandler(verifyUserOtpService)

	changePasswordRepository := repository.NewChangePasswordRepository()
	changePasswordService := business.NewChangePasswordService(changePasswordRepository, postgresClient)
	changePasswordHandler := handlers.NewChangePasswordHandler(changePasswordService)

	authGroup := router.Group(constants.AuthRoutePrefix)
	{
		authGroup.POST(constants.Signup, createUserHandler.HandleCreaterUser)
		authGroup.POST(constants.Signin, signInHandler.HandleSignIn)
		authGroup.POST(constants.Validateotp, verifyUserOtpHandler.HandleValidateUserOtp)
		authGroup.POST(constants.Changepassword, middleware.AuthMiddleware(), changePasswordHandler.HandleChangePassword)
	}

	return router
}
