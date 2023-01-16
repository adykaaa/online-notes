package main

import (
	"fmt"
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
		fmt.Errorf("Could not load config. %v", err)
	}

	pg, err := db.NewSqlDB("postgres", config.DBConnString)
	if err != nil {
		fmt.Errorf("error trying to connect to postgres. %v", err)
	}

	db = db.New(pg.GetDB())

	migrations.MigrateDB(config.DBConnString)
	if err != nil {
		fmt.Errorf("Database migration failure! %v", err)
	}

	r := http.NewChiRouter()
	httpServer := http.NewServer(r, ":8080")

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
