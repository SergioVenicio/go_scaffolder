package publisher

import "context"

type Publisher interface {
	Publish(ctx context.Context, location string, data []byte) error
}
