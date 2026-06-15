package main

import (
	"log"
	"stock_broker_application/src/constants"
	"stock_broker_application/src/models"
	"stock_broker_application/src/utils"
)

// @title stock_broker_application-backend
// @version 1.0
// @description stock_broker_application-backend
// @BasePath /v1
// @query.collection.format multi
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @x-extension-openapi {"example": "value on a json format"}
func main() {
	err := utils.InitPostgresConfg(constants.RootConfig)
	if err != nil {
		log.Fatalf(constants.ErrDBInitFailed, err)
		return
	}

	// Perform Migrations
	dbClient := utils.GetPostgresClient()
	if err := dbClient.GormDB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf(constants.ErrDBMigrationFailed, err)
		return
	}

	log.Println(constants.MsgDBMigrationSuccess)
}
