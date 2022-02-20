package models

import (
	"database/sql"
	"encoding/gob"
	"github.com/gorilla/sessions"
	"home-counter/db"
	"home-counter/src/config"
)

var (
	Store = sessions.NewCookieStore([]byte(config.Config.CookieSecret))
	DB    = db.MakeConnect()
)

func init() {
	Store = sessions.NewCookieStore([]byte("something-very-secret"))
	Store.MaxAge(86400 * 30 * 12)
	gob.Register(map[string]interface{}{})
}

type UserData struct {
	Name string
	Id   string
}

type UserConfig struct {
	//Id string `db:"id"`
	Name              string          `db:"name"`
	ElectricityTariff sql.NullFloat64 `db:"electricity_tariff"`
	HotWaterTariff    sql.NullFloat64 `db:"hot_water_tariff"`
	ColdWaterTariff   sql.NullFloat64 `db:"cold_water_tariff"`
	OutWaterTariff    sql.NullFloat64 `db:"out_water_tariff"`
	InternetTariff    sql.NullFloat64 `db:"internet_tariff"`
	Electricity       sql.NullInt64   `db:"electricity_meter"`
	HotWater          sql.NullInt64   `db:"hot_water_meter"`
	ColdWater         sql.NullInt64   `db:"clod_water_meter"`
}
