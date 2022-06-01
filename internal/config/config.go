package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort       string `env:"PETSERA_HTTP_PORT" envDefault:"8080"`

	
	GCPBucket           string `env:"PETSERA_BUCKET,required"`

	DBHost     string `env:"PETSERA_DB_HOST,required"`
	DBPort     string `env:"PETSERA_DB_PORT" envDefault:"5432"`
	DBUser     string `env:"PETSERA_DB_USER,required"`
	DBPassword string `env:"PETSERA_DB_PASSWORD"`
	DBName     string `env:"PETSERA_DB_NAME" envDefault:"petsera"`
}

func MustLoad() Config {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("failed to load .env file: %v", err)
		}
		log.Println("ignoring missing .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}
