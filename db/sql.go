package db

import (
	"database/sql"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/rs/zerolog"
)

const (
	defaultConnAttempts = 5
	defaultConnTimeout  = 2 * time.Second
)

// sqlDB struct provides all functions to execute SQL queries using composition with the sqlc.Queries struct.
type sqlDB struct {
	logger *zerolog.Logger
	*sqlc.Queries
	db           *sql.DB
	connAttempts int
	connTimeout  time.Duration
}

func NewSQLdb(driver string, url string, logger *zerolog.Logger) (*sqlDB, error) {
	sqlDB := &sqlDB{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
		logger:       logger,
	}
	var err error

	for sqlDB.connAttempts > 0 {

		sqlDB.db, err = sql.Open(driver, url)
		if err != nil {
			logger.Error().Msgf("error trying to open DB. %v  Attempt: %d", err, sqlDB.connAttempts)
		}

		err = sqlDB.db.Ping()
		if err != nil {
			logger.Error().Msgf("error trying to connect to the DB: %v. Attempt: %d", err, sqlDB.connAttempts)
		}

		//if we could create the DB object and connect to the DB, we exit the loop
		if err == nil {
			break
		}

		time.Sleep(sqlDB.connTimeout)
		sqlDB.connAttempts--

		if sqlDB.connAttempts == 0 {
			logger.Error().Msgf("Could not establish connection to the database")
			return nil, err
		}
	}

	sqlDB.Queries = sqlc.New(sqlDB.db)
	logger.Info().Msg("DB connection is successful.")

	return sqlDB, nil
}
