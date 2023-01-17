package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	sqldb, err := db.NewSQLdb("postgres", config.DBConnString)
	if err != nil {
		log.Fatalf("error initializing SQL db. %v", err)
	}

	err = migrations.MigrateDB(config.DBConnString)
	if err != nil {
		log.Fatalf("database migration failure! %v", err)
	}

	repo := db.NewRepository(sqldb)
	r := http.NewChiRouter(repo)

	httpServer, err := http.NewServer(r, ":8080")
	if err != nil {
		log.Fatalf("error during server initialization! %v", err)
	}

	log.Printf("HTTP server is now listening on %v")
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		fmt.Printf("Server run interrupted by signal %s", s.String())
	case err := <-httpServer.Notify():
		fmt.Printf("Server connection error %v", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		fmt.Errorf("app - Run - httpServer.Shutdown: %v", err)
	}
}
