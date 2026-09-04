package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/k4k3ru-hub/btse/go/websocket/protocol"
)

const rawOrderBookGrouping = 0

// OrderBookClient manages level-two order book subscriptions.
type OrderBookClient struct{ executor Executor }

// OrderBookParams configures a level-two order book subscription.
type OrderBookParams struct{ Symbol string }

// Key returns the BTSE raw order book topic.
//
// Version:
//   - 2026-09-04: Added.
func (p OrderBookParams) Key() (string, error) {
	if p.Symbol == "" {
		return "", fmt.Errorf("failed to build order book subscription key: symbol=empty")
	}
	return fmt.Sprintf("%s:%s_%d", protocol.TopicOrderBook, p.Symbol, rawOrderBookGrouping), nil
}

// NewOrderBookClient creates a level-two order book subscription client.
//
// Version:
//   - 2026-09-04: Added.
func NewOrderBookClient(executor Executor) (*OrderBookClient, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create order book client: executor=null")
	}
	return &OrderBookClient{executor: executor}, nil
}

// Subscribe subscribes to level-two order book snapshots and deltas.
//
// Version:
//   - 2026-09-04: Added.
func (c *OrderBookClient) Subscribe(ctx context.Context, params OrderBookParams) error {
	return c.execute(ctx, protocol.OperationSubscribe, params)
}

// Unsubscribe unsubscribes from level-two order book snapshots and deltas.
//
// Version:
//   - 2026-09-04: Added.
func (c *OrderBookClient) Unsubscribe(ctx context.Context, params OrderBookParams) error {
	return c.execute(ctx, protocol.OperationUnsubscribe, params)
}

func (c *OrderBookClient) execute(ctx context.Context, operation string, params OrderBookParams) error {
	if c == nil {
		return fmt.Errorf("failed to %s order book: client=null", operation)
	}
	key, err := params.Key()
	if err != nil {
		return fmt.Errorf("failed to %s order book: %w", operation, err)
	}
	payload, err := json.Marshal(protocol.SubscribeRequest{Operation: operation, Arguments: []string{key}})
	if err != nil {
		return fmt.Errorf("failed to %s order book: failed to encode request: %w", operation, err)
	}
	if operation == protocol.OperationSubscribe {
		err = c.executor.Subscribe(ctx, key, payload)
	} else {
		err = c.executor.Unsubscribe(ctx, key, payload)
	}
	if err != nil {
		return fmt.Errorf("failed to %s order book: %w", operation, err)
	}
	return nil
}
