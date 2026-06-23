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
	"github.com/redis/go-redis/v9"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func GetRouter(gdb *gorm.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	docs.SwaggerInfo.Title = constants.SwaggerTitle

	router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(files.Handler))

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{genericConstants.AllowedOrigin},
		AllowMethods: []string{genericConstants.POST, genericConstants.GET},
		AllowHeaders: []string{genericConstants.Origin, genericConstants.ContentType, genericConstants.Authorization},
	}))

	createUserRepository := repository.NewCreateUserRepository(gdb)
	createUserService := business.NewCreateUserService(createUserRepository)
	createUserHandler := handlers.NewCreateUserHandler(createUserService)

	signinUserRepository := repository.NewSigninUserRepository(gdb)
	signinUserService := business.NewSigninUserService(signinUserRepository)
	signinUserHandler := handlers.NewSigninUserHandler(signinUserService)

	forgotPasswordRepository := repository.NewForgotPasswordRepository(gdb)
	forgotPasswordService := business.NewForgotPasswordService(forgotPasswordRepository)
	forgotPasswordHandler := handlers.NewForgotPasswordHandler(forgotPasswordService)

	verifyUserOtpRepository := repository.NewValidateUserOtpRepository(gdb)
	verifyUserOtpService := business.NewValidateUserOtpService(verifyUserOtpRepository)
	verifyUserOtpHandler := handlers.NewValidateUserOtpHandler(verifyUserOtpService)

	changePasswordRepository := repository.NewChangePasswordRepository(gdb)
	changePasswordService := business.NewChangePasswordService(changePasswordRepository)
	changePasswordHandler := handlers.NewChangePasswordHandler(changePasswordService)

	logoutUserService := business.NewLogoutUser(redisClient)
	logoutUserHandler := handlers.LogoutUserHandler(logoutUserService)

	authGroup := router.Group(constants.RoutePrefix)
	{
		authGroup.POST(constants.Signup, createUserHandler.HandleCreaterUser)
		authGroup.POST(constants.Signin, signinUserHandler.HandleSigninUser)
		authGroup.POST(constants.ForgotPassword, forgotPasswordHandler.HandleForgotPassword)
		authGroup.POST(constants.ValidateOtp, verifyUserOtpHandler.HandleValidateUserOtp)
		authGroup.POST(constants.ChangePassword, middleware.AuthMiddleware(), changePasswordHandler.HandleChangePassword)
		authGroup.POST(constants.Logout, middleware.AuthMiddleware(), logoutUserHandler.HandleUserLogout)

	}
	return router
}
