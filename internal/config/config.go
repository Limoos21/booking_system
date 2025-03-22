package config

import (
	"os"
	"strconv"
)

type Config struct {
	HttpPort         string
	DsnDatabase      string
	KafkaServer      string
	LogLevel         string
	NameServiceKafka string
	TokenBot         string
}

func NewConfig() *Config {
	return &Config{
		HttpPort:         getEnv("HTTP_PORT", "8080"),
		DsnDatabase:      getEnv("DSN_DATABASE", "host=localhost port=5432 user=classnay_namy_y_admina228 password=classnay_password_sdelann1dmin dbname=basic_db sslmode=disable"),
		KafkaServer:      getEnv("KAFKA_SERVER", "localhost:9092"),
		LogLevel:         getEnv("LOG_LEVEL", "debug"),
		NameServiceKafka: getEnv("NAME_SERVICE_KAFKA", ""),
		TokenBot:         getEnv("TOKEN_BOT", "7617376673:AAHLqRlZN21_FeIxduDLDvV0-Z6XQnCmeBw"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (c *Config) GetHttpPort() int {
	port, err := strconv.Atoi(c.HttpPort)
	if err != nil {
		panic(err)
	}
	return port
}
