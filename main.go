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
	dbClient := utils.GetPostgresClient().GormDB
	err = dbClient.AutoMigrate(
		&models.User{},
		&models.Watchlist{},
		&models.ScripMaster{},
		&models.WatchlistScrip{},
	)
	if err != nil {
		log.Fatalf(constants.ErrDBMigrationFailed, err)
	}

	log.Println(constants.MsgDBMigrationSuccess)
}

//try doing automigrate in single functioncall and handle the error only once
