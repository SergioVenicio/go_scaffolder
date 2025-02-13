package publisher

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
)

type PubsubPublisher struct {
	client *pubsub.Client
}

func NewPubsub(ctx context.Context, projectID string) Publisher {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	return &PubsubPublisher{
		client: client,
	}
}

func (p *PubsubPublisher) Publish(ctx context.Context, location string, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()
	t := p.client.Topic(location)
	ok, err := t.Exists(ctx)
	if !ok || err != nil {
		t, err = p.client.CreateTopic(ctx, location)
		if err != nil {
			return err
		}
	}
	msg := &pubsub.Message{
		Data: data,
	}
	result := t.Publish(ctx, msg)
	_, err = result.Get(ctx)
	if err != nil {
		return err
	}

	return nil
}
