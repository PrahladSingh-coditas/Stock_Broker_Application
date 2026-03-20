package main

import (
	"fmt"
	"log"
	"stock_broker_application/src/app/watchlist/commons/constants"
	"stock_broker_application/src/app/watchlist/router"
	"stock_broker_application/src/utils"

	"github.com/sirupsen/logrus"
)

// @title backend
// @version 1.0
// @description backend for watchlist micro-service 
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

	startRouter()
}

func startRouter() {
	logger := logrus.New()
	router := router.GetRouter()
	logger.Info(fmt.Sprintf(constants.RunningServerPort, constants.PortDefaultValude))
	router.Run(fmt.Sprintf(":%d", constants.PortDefaultValude))
}
