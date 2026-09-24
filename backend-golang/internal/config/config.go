package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	AppEnv             string
	DatabaseUrl        string
	BaseUrlIBGE        string
	BaseUrlLocalidades string
	BaseUrlANS         string
	AnsFonte           string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:               os.Getenv("PORT"),
		AppEnv:             os.Getenv("APP_ENV"),
		DatabaseUrl:        os.Getenv("DATABASE_URL"),
		BaseUrlIBGE:        os.Getenv("BASE_URL_IBGE"),
		BaseUrlLocalidades: os.Getenv("BASE_URL_LOCALIDADES"),
		BaseUrlANS:         os.Getenv("BASE_URL_ANS"),
		AnsFonte:           os.Getenv("ANS_FONTE"),
	}
}
