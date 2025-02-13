package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	Env string `envconfig:"ENV" env-default:"local"`
	HTTPServer
	DBConfig
}

type HTTPServer struct {
	Address string `envconfig:"SERVER_ADDRESS" env-default:":8080"`
}

type DBConfig struct {
	ConnectionString string `envconfig:"POSTGRES_CONN"`
	Username         string `envconfig:"POSTGRES_USERNAME"`
	Password         string `envconfig:"POSTGRES_PASSWORD"`
	Host             string `envconfig:"POSTGRES_HOST"`
	Port             int    `envconfig:"POSTGRES_PORT" env-default:"5432"`
	Database         string `envconfig:"POSTGRES_DATABASE"`
}

func New() *Config {
	const op = "config.new"
	var cfg Config

	err := envconfig.Process(op, &cfg)
	if err != nil {
		log.Fatal(op, err)
	}
	return &cfg
}
