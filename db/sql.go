package db

import (
	"database/sql"
	"fmt"
	"log"
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
			log.Printf("error trying to open DB. %v  Attempt: %v", err, sqlDB.connAttempts)
		}

		err = sqlDB.db.Ping()
		if err != nil {
			log.Printf("error trying to connect to the DB: %v. Attempt: %d", err, sqlDB.connAttempts)
		}

		//if we could create the DB object and connect to the DB, we exit the loop
		if err == nil {
			break
		}

		time.Sleep(sqlDB.connTimeout)
		sqlDB.connAttempts--

		if sqlDB.connAttempts == 0 {
			log.Fatalf("Could not establish DB connection, exiting...")
		}
	}

	sqlDB.Queries = sqlc.New(sqlDB.db)
	log.Println("DB connection is successful.")
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
