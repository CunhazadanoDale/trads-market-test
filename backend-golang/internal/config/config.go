package config

import "os"

type Config struct {
	Port        string
	AppEnv string
	DatabaseUrl string
	BaseUrlIBGE string
}

func LoadConfig() *Config {
	return &Config {
		Port:        os.Getenv("PORT"),
		AppEnv: os.Getenv("APP_ENV"),
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		BaseUrlIBGE: os.Getenv("BASE_URL_IBGE"),
	}
}