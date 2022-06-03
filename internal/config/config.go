package config

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DevMode bool `env:"PETSERA_DEV_MODE" envDefault:"false"`

	HTTPPort string `env:"PETSERA_HTTP_PORT" envDefault:"8080"`

	GCSBucket    string `env:"PETSERA_BUCKET,required"`
	GoogleAPIKey string `env:"PETSERA_GOOGLE_API_KEY,required"`

	DBHost               string `env:"PETSERA_DB_HOST,required"`
	DBPort               string `env:"PETSERA_DB_PORT" envDefault:"5432"`
	DBUser               string `env:"PETSERA_DB_USER,required"`
	DBPassword           string `env:"PETSERA_DB_PASSWORD,required"`
	DBName               string `env:"PETSERA_DB_NAME" envDefault:"petsera"`
	DBMaxOpenConnections int    `env:"PETSERA_MAX_OPEN_CONNECTIONS" envDefault:"0"`
	DBMaxIdleConnections int    `env:"PETSERA_MAX_IDLE_CONNECTIONS" envDefault:"5"`
	DBEnableMigration    bool   `env:"PETSERA_ENABLE_MIGRATION" envDefault:"true"`
	DBMigrations         string `env:"PETSERA_MIGRATIONS" envDefault:"scripts/migrations/postgres"`
}

func MustLoad() Config {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			logrus.WithError(err).Fatal("failed to load .env file")
		}
		logrus.Warn("ignoring missing .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	return cfg
}
