package main

import (
	ServiceConstants "authentication/commons/constants"
	"authentication/router"
	"fmt"
	"log"
	"stock_broker_application/src/constants"
	"stock_broker_application/src/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// @title backend
// @version 1.0
// @description backend for Auth micro-service (Middleware layer).
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

	startRouter(postgresClient)
}

func startRouter(gdb *gorm.DB) {
	logger := logrus.New()
	router := router.GetRouter(gdb)
	logger.Info(fmt.Sprintf(constants.RunningServerPort, ServiceConstants.PortDefaultValude))
	router.Run(fmt.Sprintf(":%d", ServiceConstants.PortDefaultValude))
}
