package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SergioVenicio/go_scaffolder/config"
	"github.com/SergioVenicio/go_scaffolder/consumer"
	"github.com/SergioVenicio/go_scaffolder/database"
	"github.com/SergioVenicio/go_scaffolder/models"
	"github.com/SergioVenicio/go_scaffolder/repositories"
	"github.com/sirupsen/logrus"
)

var (
	db     database.Database
	cfg    *config.Config
	logger *logrus.Logger
)

func setUpLogger() {
	logger = logrus.New()
	logLevel := logrus.WarnLevel
	switch cfg.LogLevel {
	case "DEBUG":
		logLevel = logrus.DebugLevel
	case "INFO":
		logLevel = logrus.InfoLevel
	case "WARNING":
		logLevel = logrus.WarnLevel
	case "ERROR":
		logLevel = logrus.ErrorLevel
	default:
		logLevel = logrus.ErrorLevel
	}
	logger.SetLevel(logLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
}

var once sync.Once

func setUpDatabase() {
	db = database.NewPostgresql(cfg)
	db.GetDb().AutoMigrate(&models.Product{})
	db.GetDb().AutoMigrate(&models.Image{})
}

func init() {
	cfg = config.NewConfig()
	once.Do(func() {
		setUpLogger()
		setUpDatabase()
	})
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
}

func main() {
	ctx := context.Background()
	products := repositories.NewProductRepository(db, logger)
	productCreateConsumer := consumer.NewCreateProductConsumer(ctx, cfg, logger, products)
	imageDownloadConsumer := consumer.NewImageConsumer(ctx, cfg, logger, products)
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for range cfg.MaxWorkers {
		go productCreateConsumer.Consume()
		go imageDownloadConsumer.Consume()
	}

	go func() {
		sig := <-sigChan
		switch sig {
		default:
			productCreateConsumer.Shutdown()
			imageDownloadConsumer.Shutdown()
		}
	}()
	<-sigChan
}
