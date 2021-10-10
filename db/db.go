package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"home-counter/src/config"
	"log"
)

func MakeConnect () *sql.DB {
	db, err := sql.Open("pgx", config.Config.DbDSN)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(1)
	return db
}
