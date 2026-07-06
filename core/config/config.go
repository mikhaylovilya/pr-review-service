package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServerPort int `env:"HTTP_SERVER_PORT"`
	Db             PostgresConfig
}

type PostgresConfig struct {
	User     string `env:"POSTGRES_PORT"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	Name     string `env:"POSTGRES_NAME"`
	URI      string `env:"POSTGRES_URI"`
}

func NewConfig() *Config {
	wd, wd_err := os.Getwd()
	if wd_err != nil {
		log.Fatalf("Could not get wd: %v", wd_err)
		return nil
	}

	defaultEnvPath := wd + "/.env"
	if _, defaultEnvPath_err := os.Stat(defaultEnvPath); defaultEnvPath_err != nil {
		log.Fatalf(".env file does not exist: %v", defaultEnvPath_err)
		return nil
	}

	var cfg Config

	if readenv_err := cleanenv.ReadConfig(defaultEnvPath, &cfg); readenv_err != nil {
		log.Fatalf("Could not parse %s file: %v", defaultEnvPath, readenv_err)
		return nil
	}

	return &cfg
}
