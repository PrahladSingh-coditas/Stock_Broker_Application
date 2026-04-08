package main

import (
	"fmt"
	"log"
	"stock_broker_application/src/constants"
	"stock_broker_application/src/utils"
	"watchlists/router"

	ServiceConstants "watchlists/commons/constants"

	"github.com/sirupsen/logrus"
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

	if err := utils.InitRedisConfig(); err != nil {
		log.Fatalf(constants.ErrRedisConfigFailed, err)
	}

	startRouter()
}

func startRouter() {
	logger := logrus.New()
	router := router.GetRouter()
	logger.Info(fmt.Sprintf(constants.RunningServerPort, ServiceConstants.PortDefaultValude))
	router.Run(fmt.Sprintf(":%d", ServiceConstants.PortDefaultValude))
}
