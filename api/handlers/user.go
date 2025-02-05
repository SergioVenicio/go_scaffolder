package handlers

import (
	"github.com/Solana-Listener/payo/config"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	cfg    *config.Config
	logger *logrus.Logger
}

func NewUserHandler(cfg *config.Config, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		cfg:    cfg,
		logger: logger,
	}
}

func (userHandler *UserHandler) Root(c *fiber.Ctx) error {
	msg := struct {
		Message string `json:"message"`
	}{Message: "ok"}
	return c.JSON(msg)
}
