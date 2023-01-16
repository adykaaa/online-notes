package db

import (
	"database/sql"
	"fmt"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
)

const (
	defaultConnAttempts = 5
	defaultConnTimeout  = 2 * time.Second
)

// sqlDB struct provides all functions to execute SQL queries using composition with the sqlc.Queries struct.
type sqlDB struct {
	*sqlc.Queries
	db           *sql.DB
	connAttempts int
	connTimeout  time.Duration
}

func NewSQLdb(driver string, url string) (*sqlDB, error) {
	sqlDB := &sqlDB{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}
	var err error
	for sqlDB.connAttempts > 0 {
		sqlDB.db, err = sql.Open(driver, url)
		if err != nil {
			fmt.Printf("error trying to connect to the database on %s. %v.  Attempts left: %v", url, err, sqlDB.connAttempts)
			break
		}
		time.Sleep(sqlDB.connTimeout)
		sqlDB.connAttempts--
	}

	sqlDB.Queries = sqlc.New(sqlDB.db)

	fmt.Print("connection to DB was successful")
	return sqlDB, nil
}

func (db *sqlDB) GetDB() *sql.DB {
	return db.db
}

func (db *sqlDB) Close() {
	if db.db != nil {
		err := db.db.Close()
		if err != nil {
			fmt.Printf("error closing the Postgre DB connection")
		}
	}
}
