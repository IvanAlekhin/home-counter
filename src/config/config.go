package config

import (
	"errors"
	"fmt"
)
import "github.com/caarlos0/env/v6"

var Config = config{}
var ErrorConfig = errors.New("config parsing error")

type config struct {
	Port         string    `env:"PORT,notEmpty" envDefault:"8080"`
	DbDSN        string `env:"DbDSN,notEmpty" envDefault:"postgres://postgres:1@localhost:5432/home_counter"`
	IsProduction bool   `env:"PRODUCTION,notEmpty"`
	CookieSecret string `env:"COOKIE_SECRET,notEmpty"`
	AppUrl string `env:"APP_URL,notEmpty"`
	AuthSecret string `env:"AUTH_SECRET,notEmpty"`
	AuthUrl string `env:"AUTH_URL,notEmpty"`
	AuthId string `env:"AUTH_ID,notEmpty"`
}

func init () {
	err := env.Parse(&Config)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(ErrorConfig)
	}
}
