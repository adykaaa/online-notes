package migrations

import (
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 5
	defaultTimeout  = 2 * time.Second
)

func MigrateDB(db_url string) error {
	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://db/migrations/", db_url)
		if err == nil {
			break
		}

		log.Printf("Migration: DB is trying to connect, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: DB connect error: %s", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	defer m.Close()

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return nil
	}
	log.Printf("Migrate: up success")
	return nil
}
