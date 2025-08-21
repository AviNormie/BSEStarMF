package consumer

import "context"

type MessageHandler interface {
	HandleMessage(ctx context.Context, key, value []byte) error
}