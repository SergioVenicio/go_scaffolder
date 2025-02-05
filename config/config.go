package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel    string
	DatabaseURI string
	ServerPort  string
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("load on load env file .env: %s", err)
	}

	cfg := &Config{
		LogLevel:    viper.GetString("LOG_LEVER"),
		DatabaseURI: viper.GetString("DATABASE_URI"),
		ServerPort:  viper.GetString("SERVER_PORT"),
	}
	return cfg
}
