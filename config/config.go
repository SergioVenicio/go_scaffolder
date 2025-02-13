package config

import (
	"cmp"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel             string
	DatabaseURI          string
	ServerPort           string
	ProductCreationTopic string
	ProjectID            string
	ProductCreationSubID string
	ImageSubID           string
	MaxWorkers           int
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("load on load env file .env: %s", err)
	}

	cfg := &Config{
		LogLevel:             viper.GetString("LOG_LEVEL"),
		DatabaseURI:          viper.GetString("DATABASE_URI"),
		ServerPort:           viper.GetString("SERVER_PORT"),
		ProjectID:            viper.GetString("PUBSUB_PROJECT_ID"),
		ProductCreationTopic: viper.GetString("PUBSUB_PRODUCTS_CREATE_TOPIC_ID"),
		ProductCreationSubID: viper.GetString("PUBSUB_PRODUCTS_CREATE_SUB_ID"),
		ImageSubID:           viper.GetString("PUBSUB_IMAGE_SUB_ID"),
		MaxWorkers:           cmp.Or(viper.GetInt("MAX_WORKERS"), 5),
	}
	return cfg
}
