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
		log.Fatal(err)
	}

	logger := utils.NewLogger(config.LogLevel)

	sqldb, err := db.NewSQLdb("postgres", config.DBConnString, &logger)
	if err != nil {
		logger.Fatal().Err(err).Stack().Send()
	}

	err = migrations.MigrateDB(config.DBConnString, &logger)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	router := http.NewChiRouter(sqldb, &logger)

	httpServer, err := http.NewServer(router, config.HTTPServerAddress, &logger)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	err = httpServer.Shutdown()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
}
