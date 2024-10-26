package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mvp-mogila/ozon-test-task/pkg/postgres"
)

type Config struct {
	Host               string `env:"HOST" env-default:"localhost"`
	Port               string `env:"PORT" env-default:"8080"`
	RequestTimeout     int    `env:"TIMEOUT" env-default:"3"`
	UseInMemoryStorage bool   `env:"INMEM" env-default:"true"`
	PostgresHost       string `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresPort       string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser       string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresPassword   string `env:"POSTGRES_PASSWORD"`
	DatabaseName       string `env:"POSTGRES_DB" env-default:"postgres"`
}

func LoadConfig() *Config {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	return &cfg
}

func InitPostgresConfig(cfg *Config) *postgres.PostgresConfig {
	return &postgres.PostgresConfig{
		User:     cfg.PostgresUser,
		Password: cfg.PostgresPassword,
		Host:     cfg.PostgresHost,
		Port:     cfg.PostgresPort,
		Database: cfg.DatabaseName,
	}
}
