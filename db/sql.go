package db

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	defaultConnAttempts = 5
	defaultConnTimeout  = 2 * time.Second
)

// SQLdb struct provides all functions to execute SQL queries
type SQLdb struct {
	db *sql.DB
	*Queries
	connAttempts int
	connTimeout  time.Duration
}

func NewSqlDB(url string, driver string) (*SQLdb, error) {
	SQLdb := &SQLdb{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}
	var err error
	for SQLdb.connAttempts > 0 {
		SQLdb.db, err = sql.Open(driver, string(url))
		if err != nil {
			fmt.Printf("error trying to connect to the database on %s. %v.  Attempts left: %v", url, err, SQLdb.connAttempts)
			break
		}
		time.Sleep(SQLdb.connTimeout)
		SQLdb.connAttempts--
	}

	fmt.Print("connection to DB was successful")
	return SQLdb, nil
}

func (db *SQLdb) GetDB() *sql.DB {
	return db.db
}

func (db *SQLdb) Close() {
	if db.db != nil {
		err := db.db.Close()
		if err != nil {
			fmt.Printf("error closing the Postgre DB connection")
		}
	}
}
