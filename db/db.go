package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"home-counter/src/config"
	"log"
	"os"
)

func MakeConnect() *pgxpool.Pool {
	cfg, err := pgx.ParseConfig(config.Config.DbDSN)
	if err != nil {
		log.Printf("Unexpected dsn for database")
		panic(err)
	}
	cfg.PreferSimpleProtocol = true
	cfg.RuntimeParams["standard_conforming_strings"] = "on"

	db, err := pgxpool.Connect(context.Background(), cfg.ConnString())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return db
}
