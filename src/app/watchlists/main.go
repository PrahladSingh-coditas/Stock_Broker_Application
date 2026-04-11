package main

import (
	"fmt"
	"log"
	"stock_broker_application/src/constants"
	"stock_broker_application/src/utils"
	"watchlists/router"

	ServiceConstants "watchlists/commons/constants"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// @title Watchlists ADG Service API
// @version 1.0
// @description Watchlists ADD DEL GET APIs for Stock Broker Application
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

	if _, err := utils.InitRedis(); err != nil {
		log.Fatalf(constants.ErrRedisInitFailed, err)
	}

	redisClient := utils.GetRedisClient()
	db := utils.GetPostgresClient().GormDB

	startRouter(db, redisClient)
}

func startRouter(db *gorm.DB, redisClient *redis.Client) {
	logger := logrus.New()
	router := router.GetRouter(db, redisClient)
	logger.Info(fmt.Sprintf(constants.RunningServerPort, ServiceConstants.PortDefaultValude))
	router.Run(fmt.Sprintf(":%d", ServiceConstants.PortDefaultValude))
}
