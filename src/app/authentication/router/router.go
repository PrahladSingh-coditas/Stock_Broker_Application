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

	createUserRepository := repository.NewCreateUserRepository()
	createUserService := business.NewCreateUserService(createUserRepository)
	createUserHandler := handlers.NewCreateUserHandler(createUserService)

	signinRepository := repository.NewSigninUserRepository()
	signinService := business.NewSigninUserService(signinRepository)
	signinHandler := handlers.NewSigninUserHandler(signinService)

	//validate OTP
	verifyUserOtpRepository := repository.NewValidateUserOtpRepository()
	verifyUserOtpService := business.NewValidateUserOtpService(verifyUserOtpRepository)
	verifyUserOtpHandler := handlers.NewValidateUserOtpHandler(verifyUserOtpService)

	//change password
	changePasswordRepository:= repository.NewChangePasswordRepository()
	changePasswordService := business.NewChangePasswordService(changePasswordRepository)
	changePasswordHandler := handlers.NewChangePasswordHandler(changePasswordService)

	authGroup := router.Group(constants.AuthRoutePrefix)
	{
		authGroup.POST(constants.Signup, createUserHandler.HandleCreaterUser)
		authGroup.POST(constants.Signin, signinHandler.HandleSigninUser)

		authGroup.POST(constants.Validateotp, verifyUserOtpHandler.HandleValidateUserOtp)

		authGroup.POST(constants.ChangePassword, middleware.AuthMiddleware(), changePasswordHandler.HandleChangePassword)
	}

	return router
}
