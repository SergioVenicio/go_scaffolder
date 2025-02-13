package server

import (
	"fmt"

	"github.com/SergioVenicio/go_scaffolder/config"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	app    *fiber.App
	cfg    *config.Config
	logger *logrus.Logger
}

func NewHttpServer(cfg *config.Config, logger *logrus.Logger) *HttpServer {
	app := fiber.New()

	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 60,
	}))

	return &HttpServer{
		app: fiber.New(fiber.Config{
			Prefork:     true,
			JSONEncoder: json.Marshal,
			JSONDecoder: json.Unmarshal,
		}),
		cfg:    cfg,
		logger: logger,
	}
}

func (httpServer *HttpServer) AddHandler(method string, path string, handler fiber.Handler) {
	switch method {
	case "GET":
		httpServer.app.Get(path, handler)
	case "POST":
		httpServer.app.Post(path, handler)
	}
}
func (httpServer *HttpServer) Start() {
	httpServer.logger.WithFields(logrus.Fields{"port": httpServer.cfg.ServerPort}).Info("server running")
	httpServer.logger.Fatal(httpServer.app.Listen(fmt.Sprintf(":%s", httpServer.cfg.ServerPort)))
}
