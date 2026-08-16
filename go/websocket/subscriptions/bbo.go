package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/k4k3ru-hub/btse/go/websocket/protocol"
)

// BBOClient manages best-bid/best-ask subscriptions.
type BBOClient struct{ executor Executor }

// BBOParams configures a best-bid/best-ask subscription.
type BBOParams struct{ Symbol string }

// Key returns the BTSE BBO topic.
//
// Version:
//   - 2026-08-16: Added.
func (p BBOParams) Key() (string, error) {
	if p.Symbol == "" {
		return "", fmt.Errorf("failed to build BBO subscription key: symbol=empty")
	}
	return protocol.TopicBBO + ":" + p.Symbol, nil
}

// NewBBOClient creates a BBO subscription client.
//
// Version:
//   - 2026-08-16: Added.
func NewBBOClient(executor Executor) (*BBOClient, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create BBO client: executor=null")
	}
	return &BBOClient{executor: executor}, nil
}

// Subscribe subscribes to BBO snapshots.
//
// Version:
//   - 2026-08-16: Added.
func (c *BBOClient) Subscribe(ctx context.Context, params BBOParams) error {
	return c.execute(ctx, protocol.OperationSubscribe, params)
}

// Unsubscribe unsubscribes from BBO snapshots.
//
// Version:
//   - 2026-08-16: Added.
func (c *BBOClient) Unsubscribe(ctx context.Context, params BBOParams) error {
	return c.execute(ctx, protocol.OperationUnsubscribe, params)
}

func (c *BBOClient) execute(ctx context.Context, operation string, params BBOParams) error {
	if c == nil {
		return fmt.Errorf("failed to %s BBO: client=null", operation)
	}
	key, err := params.Key()
	if err != nil {
		return fmt.Errorf("failed to %s BBO: %w", operation, err)
	}
	payload, err := json.Marshal(protocol.SubscribeRequest{Operation: operation, Arguments: []string{key}})
	if err != nil {
		return fmt.Errorf("failed to %s BBO: failed to encode request: %w", operation, err)
	}
	if operation == protocol.OperationSubscribe {
		err = c.executor.Subscribe(ctx, key, payload)
	} else {
		err = c.executor.Unsubscribe(ctx, key, payload)
	}
	if err != nil {
		return fmt.Errorf("failed to %s BBO: %w", operation, err)
	}
	return nil
}
