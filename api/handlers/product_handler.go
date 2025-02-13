package handlers

import (
	"context"

	"github.com/SergioVenicio/go_scaffolder/config"
	"github.com/SergioVenicio/go_scaffolder/models"
	"github.com/SergioVenicio/go_scaffolder/publisher"
	"github.com/SergioVenicio/go_scaffolder/repositories"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	cfg       *config.Config
	logger    *logrus.Logger
	publisher publisher.Publisher
	products  repositories.Repository[models.Product]
}

func NewProductHandler(
	cfg *config.Config,
	logger *logrus.Logger,
	productRepository repositories.Repository[models.Product],
	publisher publisher.Publisher,
) *ProductHandler {
	return &ProductHandler{
		cfg:       cfg,
		logger:    logger,
		products:  productRepository,
		publisher: publisher,
	}
}

func (productHandler *ProductHandler) NewProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		productHandler.logger.WithError(err).Warnf("invalid payload %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := product.Validate(); err != nil {
		productHandler.logger.WithError(err).Warnf("invalid payload %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	product.NewID()
	for idx := range product.Images {
		product.Images[idx].ID = uuid.New()
		product.Images[idx].Product = product.ID
	}
	ctx := context.Background()
	topic := productHandler.cfg.ProductCreationTopic
	msg, err := json.Marshal(&product)
	if err != nil {
		productHandler.logger.WithError(err).Errorf("error on marshal %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := productHandler.publisher.Publish(ctx, topic, msg); err != nil {
		productHandler.logger.WithError(err).Errorf("error on publish %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}
