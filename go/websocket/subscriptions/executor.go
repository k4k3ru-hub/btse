package subscriptions

import "context"

// Executor executes WebSocket subscription operations.
type Executor interface {
	Subscribe(ctx context.Context, key string, payload []byte) error
	Unsubscribe(ctx context.Context, key string, payload []byte) error
}
