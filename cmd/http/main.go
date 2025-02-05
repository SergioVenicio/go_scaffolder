package main

import (
	"github.com/Solana-Listener/payo/api/handlers"
	"github.com/Solana-Listener/payo/config"
	"github.com/Solana-Listener/payo/server"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.NewConfig()
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logLevel := logrus.WarnLevel
	switch cfg.LogLevel {
	case "DEBUG":
		logLevel = 5
	case "INFO":
		logLevel = 4
	case "WARNING":
		logLevel = 3
	case "ERROR":
		logLevel = 2
	default:
		logLevel = 2
	}
	logger.SetLevel(logLevel)
	httpServer := server.NewHttpServer(cfg, logger)
	userHandler := handlers.NewUserHandler(cfg, logger)
	httpServer.AddHandler("GET", "/", userHandler.Root)
	httpServer.Start()
}
