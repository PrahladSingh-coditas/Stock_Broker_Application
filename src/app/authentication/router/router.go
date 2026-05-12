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
	//router.Use(middleware.AuthMiddleware())
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

	signInUserRepository := repository.NewSignInUserRepository(gdb)
	signInUserService := business.NewSignInUserService(signInUserRepository)
	signInUserHandler := handlers.NewSignInUserHandler(signInUserService)

	validateUserOtpRepository := repository.NewValidateUserOtpRepository(gdb)
	validateUserOtpService := business.NewValidateUserOtpService(validateUserOtpRepository)
	validateUserOtpHandler := handlers.NewValidateUserOtpHandler(validateUserOtpService)

	changePasswordRepository := repository.NewChangePasswordRepository(gdb)
	changePasswordService := business.NewChangePasswordService(changePasswordRepository)
	changePasswordHandler := handlers.NewChangePasswordHandler(changePasswordService)

	logoutUserService := business.NewLogoutUserService(redisClient)
	logoutUserHandler := handlers.NewLogoutUserHandler(logoutUserService)

	authGroup := router.Group(constants.AuthRoutePrefix)
	{
		authGroup.POST(constants.Signup, createUserHandler.HandleCreaterUser)
		authGroup.POST(constants.Signin, signInUserHandler.HandleSignInUser)
		authGroup.POST(constants.ValidateOtp, validateUserOtpHandler.HandleValidateUserOtp)
		authGroup.PATCH(constants.ChangePassword, middleware.AuthMiddleware(redisClient), changePasswordHandler.HandleChangePassword)
		authGroup.POST(constants.Logout, middleware.AuthMiddleware(redisClient), logoutUserHandler.HandleLogoutUser)
	}

	return router
}
