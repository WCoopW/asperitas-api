package config

import (
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port       int    `env:"PORT"`
	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	AppEnv     string `env:"APP_ENV"`
	JWT        JWT    `env:"JWT"`
}
type JWT struct {
	Secret string        `env:"JWT_SECRET"`
	Expire time.Duration `env:"JWT_EXPIRE"`
}

func Load() (Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error load .env file: %s", err)
		}
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
