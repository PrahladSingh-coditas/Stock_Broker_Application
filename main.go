package main

import (
	"log"
	"stock_broker_application/src/constants"
	"stock_broker_application/src/models"
	"stock_broker_application/src/utils"
)

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

	if err := dbClient.GormDB.AutoMigrate(&models.Watchlist{}); err != nil {
		log.Fatalf(constants.ErrDBMigrationFailed, err)
		return
	}

	if err := dbClient.GormDB.AutoMigrate(&models.ScripMaster{}); err != nil {
		log.Fatalf(constants.ErrDBMigrationFailed, err)
		return
	}

	if err := dbClient.GormDB.AutoMigrate(&models.WatchlistScrip{}); err != nil {
		log.Fatalf(constants.ErrDBMigrationFailed, err)
		return
	}

	log.Println(constants.MsgDBMigrationSuccess)
}

//try doing automigrate in single functioncall and handle the error only once
