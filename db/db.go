package db

import (
	"database/sql"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"home-counter/src/config"
	"log"
)

func MakeConnect() *sql.DB {
	cfg, err := pgx.ParseConfig(config.Config.DbDSN)
	if err != nil {
		log.Printf("Unexpected dsn for database")
		panic(err)
	}
	cfg.PreferSimpleProtocol = false

	db, err := sql.Open("pgx", cfg.ConnString())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(1)
	db.SetMaxOpenConns(1)
	return db
}
