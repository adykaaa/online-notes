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

type postgres struct {
	db           *sql.DB
	connAttempts int
	connTimeout  time.Duration
}

func NewPostgresDB(url string) (*postgres, error) {
	pg := &postgres{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}
	var err error
	for pg.connAttempts > 0 {
		pg.db, err = sql.Open("postgres", string(url))
		if err != nil {
			fmt.Printf("error trying to connect to the postgres DB on %s. %v.  Attempts left: %v", url, err, pg.connAttempts)
			break
		}
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	fmt.Print("connection to Postgres was successful")
	return pg, nil
}

func (p *postgres) GetDB() *sql.DB {
	return p.db
}

func (p *postgres) Close() {
	if p.db != nil {
		err := p.db.Close()
		if err != nil {
			fmt.Printf("error closing the Postgre DB connection")
		}
	}
}
