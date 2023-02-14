package main

import (
	"log"

	"github.com/adykaaa/online-notes/db"
	"github.com/adykaaa/online-notes/db/migrations"
	"github.com/adykaaa/online-notes/http"
	"github.com/adykaaa/online-notes/lib/config"
	logger "github.com/adykaaa/online-notes/lib/logger"
)

func main() {

	config, err := config.Load(".")
	if err != nil {
		log.Fatalf("Could not load config. %v", err)
	}

	l := logger.New(config.LogLevel)

	sqldb, err := db.NewSQLdb("postgres", config.DBConnString, &l)
	if err != nil {
		l.Fatal().Err(err).Send()
	}

	err = migrations.MigrateDB(config.DBConnString, "file://db/migrations/", &l)
	if err != nil {
		l.Fatal().Err(err).Send()
	}

	router, err := http.NewChiRouter(sqldb, config.PASETOSecret, config.AccessTokenDuration, &l)
	if err != nil {
		l.Fatal().Err(err).Send()
	}

	httpServer, err := http.NewServer(router, config.HTTPServerAddress, &l)
	if err != nil {
		l.Fatal().Err(err).Send()
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Fatal().Err(err).Send()
	}
}
