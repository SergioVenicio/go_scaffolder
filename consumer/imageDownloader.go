package consumer

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/SergioVenicio/go_scaffolder/config"
	"github.com/SergioVenicio/go_scaffolder/models"
	"github.com/SergioVenicio/go_scaffolder/repositories"
	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
)

type ImageConsumer struct {
	cfg      *config.Config
	logger   *logrus.Logger
	products repositories.Repository[models.Product]
	client   *pubsub.Client
}

func NewImageConsumer(
	ctx context.Context,
	cfg *config.Config,
	logger *logrus.Logger,
	productRepository repositories.Repository[models.Product],
) Consumer {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()
	client, err := pubsub.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		panic(err)
	}
	return &ImageConsumer{
		cfg:      cfg,
		logger:   logger,
		client:   client,
		products: productRepository,
	}
}

func (c *ImageConsumer) Consume() error {
	ctx := context.Background()
	topic := c.client.Topic(c.cfg.ProductCreationTopic)
	ok, err := topic.Exists(ctx)
	if !ok || err != nil {
		topic, err = c.client.CreateTopic(ctx, c.cfg.ProductCreationTopic)
		if err != nil {
			return err
		}
	}

	sub_id := c.cfg.ImageSubID
	topicSub := c.client.Subscription(sub_id)
	ok, err = topicSub.Exists(ctx)
	if !ok || err != nil {
		topicSub, err = c.client.CreateSubscription(ctx, sub_id, pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			c.logger.WithError(err).Error("error on subscription create")
			return err
		}
	}
	return topicSub.Receive(ctx, c.OnMessage)
}

func (c *ImageConsumer) Shutdown() {
	c.client.Close()
}

func getFileType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpeg"
	default:
		return ".png"
	}
}

func (c *ImageConsumer) OnMessage(ctx context.Context, msg *pubsub.Message) {
	defer msg.Ack()
	var product models.Product
	if err := json.Unmarshal(msg.Data, &product); err != nil {
		c.logger.WithError(err).Error("invalid message received")
		return
	}

	c.logger.WithField("product", product).Debug("product received")
	for _, img := range product.Images {
		c.logger.WithField("url", img.URL).Debug("image downloading")
		res, err := http.Get(img.URL)
		if err != nil {
			c.logger.WithError(err).WithField("url", img.URL).Error("download error")
			return
		}
		fileType := getFileType(res.Header.Get("Content-Type"))
		file, err := os.Create(fmt.Sprintf("images/%s.%s", img.ID, fileType))
		if err != nil {
			c.logger.WithError(err).WithField("url", img.URL).Error("new image file")
			return
		}

		_, err = io.Copy(file, res.Body)
		if err != nil {
			log.Fatal(err)
		}

		c.logger.WithField("url", img.URL).Info("new image file saved")
		file.Close()
		res.Body.Close()
	}
}
