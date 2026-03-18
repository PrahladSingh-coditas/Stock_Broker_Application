package router

import (
	"authentication/business"
	"authentication/commons/constants"
	"authentication/middleware"

	"authentication/docs"
	"authentication/handlers"
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
	router.Use(gin.Recovery())

	docs.SwaggerInfo.Title = constants.SwaggerTitle

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(files.Handler))

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{genericConstants.AllowedOrigin},
		AllowMethods: []string{genericConstants.POST, genericConstants.GET},
		AllowHeaders: []string{genericConstants.Origin, genericConstants.ContentType, genericConstants.Authorization},
	}))

	//signup api
	createUserRepository := repository.NewCreateUserRepository()
	createUserService := business.NewCreateUserService(createUserRepository)
	createUserHandler := handlers.NewCreateUserHandler(createUserService)

	// signin api
	SignInRepository := repository.NewSignInRepository()
	SignInService := business.NewSignInService(SignInRepository)
	SignInUserHandle := handlers.NewSignInUserHandler(SignInService)

	//validate-otp
	postgresClientDb := utils.GetPostgresClient().GormDB
	validateUserOtpRepository := repository.NewValidateUserOtpRepository()
	validateUserOtpService := business.NewValidateUserOtpService(validateUserOtpRepository, postgresClientDb)
	validateUserOtpHandler := handlers.NewValidateUserOtpHandler(validateUserOtpService)

	//pasword-reset api
	changePasswordRepository := repository.NewChangePasswordRepository()
	changePasswordService := business.NewChangePasswordService(changePasswordRepository)
	changePasswordHandler := handlers.NewChangePasswordHandler(changePasswordService)

	authGroup := router.Group(constants.AuthRoutePrefix)
	{
		authGroup.POST(constants.Signup, createUserHandler.HandleCreaterUser)
		authGroup.POST(constants.Signin, SignInUserHandle.HandleSignInUser)
		authGroup.POST(constants.ValidateOtp, validateUserOtpHandler.HandleValidateUserOtp)
		authGroup.POST(constants.ChangePassword, middleware.AuthMiddleware(), changePasswordHandler.HandleChangePassword)

	}

	return router
}
