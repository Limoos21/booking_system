package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		Database string
		Username string
		Password string
	}
	TokenBot string
}

func NewConfig() *Config {
	return &Config{
		Database: struct {
			Host     string
			Port     string
			Database string
			Username string
			Password string
		}{Host: getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			Database: getEnv("POSTGRES_DB", "basic_db"),
			Username: getEnv("POSTGRES_USER", "classnay_namy_y_admina228"),
			Password: getEnv("POSTGRES_PASSWORD", "classnay_password_sdelann1dmin")},
		TokenBot: getEnv("TOKEN_BOT", "7617376673:AAHLqRlZN21_FeIxduDLDvV0-Z6XQnCmeBw"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
	)
}
