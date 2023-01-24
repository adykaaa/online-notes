package migrations

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
)

const (
	defaultAttempts = 5
	defaultTimeout  = 2 * time.Second
)

func MigrateDB(dbURL string, source string, logger *zerolog.Logger) error {
	var (
		attempts = defaultAttempts
		m        *migrate.Migrate
	)

	if dbURL == "" || source == "" {
		logger.Error().Msg("DB URL and migration source cannot be empty!")
		return errors.New("migration error: empty param")
	}

	_, err := os.Stat(strings.TrimPrefix(source, "file://"))
	if err != nil {
		logger.Error().Msgf("migration folder cannot be reached! %v", err)
		return err
	}

	for attempts > 0 {
		m, err = migrate.New(source, dbURL)
		if err == nil {
			break
		}

		logger.Error().Msgf("db migration failed, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		logger.Error().Msgf("db migration failed. %v", err)
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error().Msgf("migrate up error. %v", err)
		return err
	}

	defer m.Close()

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info().Msg("there were no changes since the last migration")
		return nil
	}

	logger.Info().Msg("migrating the DB was successful")
	return nil
}
