package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultPort   = "8081"
	defaultAppEnv = "development"
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
		Port:               envOrDefault("PORT", defaultPort),
		AppEnv:             envOrDefault("APP_ENV", defaultAppEnv),
		DatabaseUrl:        os.Getenv("DATABASE_URL"),
		BaseUrlIBGE:        os.Getenv("BASE_URL_IBGE"),
		BaseUrlLocalidades: os.Getenv("BASE_URL_LOCALIDADES"),
		BaseUrlANS:         os.Getenv("BASE_URL_ANS"),
		AnsFonte:           os.Getenv("ANS_FONTE"),
	}
}

func (c *Config) ValidateAPI() error {
	return required(c.DatabaseUrl, "DATABASE_URL")
}

func (c *Config) ValidateImport(target string) error {
	if err := required(c.DatabaseUrl, "DATABASE_URL"); err != nil {
		return err
	}

	if target == "ibge" || target == "all" {
		if err := required(c.BaseUrlIBGE, "BASE_URL_IBGE"); err != nil {
			return err
		}
		if err := required(c.BaseUrlLocalidades, "BASE_URL_LOCALIDADES"); err != nil {
			return err
		}
	}

	if target == "ans" || target == "all" {
		if err := required(c.BaseUrlANS, "BASE_URL_ANS"); err != nil {
			return err
		}
		if err := required(c.AnsFonte, "ANS_FONTE"); err != nil {
			return err
		}
	}

	return nil
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func required(value string, name string) error {
	if value == "" {
		return fmt.Errorf("variável %s não configurada", name)
	}

	return nil
}
