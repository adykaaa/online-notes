package main

import (
	"log"

	"github.com/adykaaa/online-notes/db"
	"github.com/adykaaa/online-notes/db/migrations"
	"github.com/adykaaa/online-notes/http"
	"github.com/adykaaa/online-notes/utils"
)

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config. %v", err)
	}

	logger := utils.NewLogger(config.LogLevel)

	sqldb, err := db.NewSQLdb("postgres", config.DBConnString, &logger)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	err = migrations.MigrateDB(config.DBConnString, "file://db/migrations/", &logger)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	router := http.NewChiRouter(sqldb, config.PASETOSecret, &logger)

	httpServer, err := http.NewServer(router, config.HTTPServerAddress, &logger)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	err = httpServer.Shutdown()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
}
