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

type sqlDB struct {
	logger *zerolog.Logger
	*sqlc.Queries
	db           *sql.DB
	connAttempts int
	connTimeout  time.Duration
}

func NewSQL(driver string, url string, l *zerolog.Logger) (*sqlDB, error) {
	sqlDB := &sqlDB{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
		logger:       l,
	}

	var err error

	for sqlDB.connAttempts > 0 {
		sqlDB.db, err = sql.Open(driver, url)
		if err != nil {
			l.Error().Msgf("error trying to open DB. %v  Attempt: %d", err, sqlDB.connAttempts)
		}

		err = sqlDB.db.Ping()
		if err != nil {
			l.Error().Msgf("error trying to connect to the DB: %v. Attempt: %d", err, sqlDB.connAttempts)
		}

		//if we could create the DB object and connect to the DB instance, we break from the loop
		if err == nil {
			break
		}
		time.Sleep(sqlDB.connTimeout)
		sqlDB.connAttempts--

		if sqlDB.connAttempts == 0 {
			l.Error().Msgf("could not establish connection to the database")
			return nil, err
		}
	}

	sqlDB.Queries = sqlc.New(sqlDB.db)
	l.Info().Msg("db connection is successful.")
	return sqlDB, nil
}
