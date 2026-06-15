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
	router.Use(gin.Recovery())

	docs.SwaggerInfo.Title = constants.SwaggerTitle

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(files.Handler))
	//router.Use(middleware.AuthMiddleware())

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{genericConstants.AllowedOrigin},
		AllowMethods: []string{genericConstants.POST, genericConstants.GET},
		AllowHeaders: []string{genericConstants.Origin, genericConstants.ContentType, genericConstants.Authorization},
	}))

	createUserRepository := repository.NewCreateUserRepository()
	createUserService := business.NewCreateUserService(createUserRepository)
	createUserHandler := handlers.NewCreateUserHandler(createUserService)

	signinUserRepository := repository.NewSigninUserRepository()
	signinUserService := business.NewSigninUserService(signinUserRepository)
	signinUserHandler := handlers.NewSigninUserHandler(signinUserService)

	forgotPasswordRepository := repository.NewForgotPasswordRepository()
	forgotPasswordService := business.NewForgotPasswordService(forgotPasswordRepository)
	forgotPasswordHandler := handlers.NewForgotPasswordHandler(forgotPasswordService)

	verifyUserOtpRepository := repository.NewValidateUserOtpRepository()
	verifyUserOtpService := business.NewValidateUserOtpService(verifyUserOtpRepository)
	verifyUserOtpHandler := handlers.NewValidateUserOtpHandler(verifyUserOtpService)

	changePasswordRepository := repository.NewChangePasswordRepository()
	changePasswordService := business.NewChangePasswordService(changePasswordRepository)
	changePasswordHandler := handlers.NewChangePasswordHandler(changePasswordService)

	/*
		authGroup := router.Group(constants.AuthRoutePrefix)
		{
			authGroup.Use(middleware.AuthMiddleware())
			authGroup.POST(constants.ChangePassword, changePasswordHandler.HandleChangePassword)
		}
	*/

	v1Group := router.Group(constants.V1RoutePrefix)
	{
		v1Group.POST(constants.Signup, createUserHandler.HandleCreaterUser)
		v1Group.POST(constants.Signin, signinUserHandler.HandleSigninUser)
		v1Group.POST(constants.ForgotPassword, forgotPasswordHandler.HandleForgotPassword)
		v1Group.POST(constants.ValidateOtp, verifyUserOtpHandler.HandleValidateUserOtp)
		v1Group.POST(constants.ChangePassword, middleware.AuthMiddleware() , changePasswordHandler.HandleChangePassword)

	}
	return router
}
