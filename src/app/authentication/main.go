package main

import (
	ServiceConstants "authentication/commons/constants"
	"authentication/router"
	"fmt"
	"log"
	"stock_broker_application/src/constants"
	"stock_broker_application/src/utils"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// @title omnenest-backend
// @version 1.0
// @description Omnenest backend for watchlist micro-service (Middleware layer).
// @query.collection.format multi
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @x-extension-openapi {"example": "value on a json format"}
func main() {

	if err := utils.InitPostgresConfg(constants.BaseConfig); err != nil {
		log.Fatalf(constants.ErrDBConnectionFailed, err)
	}

	if err := utils.InitJWTConfig(constants.BaseConfig); err != nil {
		log.Fatalf(constants.ErrJWTConfigReadFailed, err)
	}

	postgresClient := utils.GetPostgresClient().GormDB
	redisClient, _, redisError := utils.GetRedisClient(false)
	if redisError != nil {
		log.Fatalf(constants.RedisConnectionFailedError, redisError.Error())
	}

	startRouter(postgresClient, redisClient)
}

func startRouter(postgresClient *gorm.DB, redisClient *redis.Client) {

	logger := logrus.New()
	router := router.GetRouter(postgresClient, redisClient)
	logger.Info(fmt.Sprintf(constants.RunningServerPort, ServiceConstants.PortDefaultValude))
	router.Run(fmt.Sprintf(":%d", ServiceConstants.PortDefaultValude))
}
