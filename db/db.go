package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"home-counter/src/config"
	"log"
	"os"
)

// MakeConnect не используется, потому что не дружит с pg bouncer
//func MakeConnect() *pgxpool.Pool {
//	cfg, err := pgx.ParseConfig(config.Config.DbDSN)
//	if err != nil {
//		log.Printf("Unexpected dsn for database")
//		panic(err)
//	}
//	cfg.PreferSimpleProtocol = true
//	cfg.RuntimeParams["standard_conforming_strings"] = "on"
//
//	db, err := pgxpool.Connect(context.Background(), cfg.ConnString())
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
//		os.Exit(1)
//	}
//
//	return db
//}

func MakeSingleConnect() *pgx.Conn {

	// кажется, что работа с конфигом лишняя, но на всякий случай оставлю
	cfg, err2 := pgx.ParseConfig(config.Config.DbDSN)
	if err2 != nil {
		log.Printf("Unexpected dsn for database")
		panic(err2)
	}
	cfg.PreferSimpleProtocol = true
	cfg.RuntimeParams["standard_conforming_strings"] = "on"

	conn, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
