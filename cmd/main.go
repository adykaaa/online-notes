package main

import (
	"fmt"
	"log"

	"github.com/adykaaa/online-notes/db"
	"github.com/adykaaa/online-notes/db/migrations"
	"github.com/adykaaa/online-notes/http"
	"github.com/adykaaa/online-notes/utils"
)

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load configuration. %v", err)
	}

	logger := utils.NewLogger(config.LogLevel)

	sqldb, err := db.NewSQLdb("postgres", config.DBConnString)
	if err != nil {
		log.Fatalf("error initializing SQL db. %v", err)
	}

	err = migrations.MigrateDB(config.DBConnString)
	if err != nil {
		log.Fatalf("database migration failure! %v", err)
	}

	r := http.NewChiRouter(sqldb)

	httpServer, err := http.NewServer(r, config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("error during server initialization! %v", err)
	}
	logger.Info().Msg("szevasz")

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		fmt.Errorf("app - Run - httpServer.Shutdown: %v", err)
	}
}
