package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/k4k3ru-hub/btse/go/websocket/protocol"
)

// TradesClient manages public trade subscriptions.
type TradesClient struct{ executor Executor }

// TradeMarket identifies the BTSE market protocol used for trade subscriptions.
type TradeMarket string

const (
	TradeMarketSpot    TradeMarket = "spot"
	TradeMarketFutures TradeMarket = "futures"
)

// TradesParams configures a public trade subscription.
type TradesParams struct {
	Symbol string
	Market TradeMarket
}

// Key returns the BTSE trades topic.
//
// Version:
//   - 2026-08-23: Selected the trade topic by Spot or Futures market.
//   - 2026-08-16: Added.
func (p TradesParams) Key() (string, error) {
	if p.Symbol == "" {
		return "", fmt.Errorf("failed to build trades subscription key: symbol=empty")
	}
	topic := protocol.TopicFuturesTrades
	switch p.Market {
	case "", TradeMarketFutures:
	case TradeMarketSpot:
		topic = protocol.TopicSpotTrades
	default:
		return "", fmt.Errorf("failed to build trades subscription key: market=invalid")
	}
	return topic + ":" + p.Symbol, nil
}

// NewTradesClient creates a public trade subscription client.
//
// Version:
//   - 2026-08-16: Added.
func NewTradesClient(executor Executor) (*TradesClient, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create trades client: executor=null")
	}
	return &TradesClient{executor: executor}, nil
}

// Subscribe subscribes to public trade updates.
//
// Version:
//   - 2026-08-16: Added.
func (c *TradesClient) Subscribe(ctx context.Context, params TradesParams) error {
	return c.execute(ctx, protocol.OperationSubscribe, params)
}

// Unsubscribe unsubscribes from public trade updates.
//
// Version:
//   - 2026-08-16: Added.
func (c *TradesClient) Unsubscribe(ctx context.Context, params TradesParams) error {
	return c.execute(ctx, protocol.OperationUnsubscribe, params)
}

func (c *TradesClient) execute(ctx context.Context, operation string, params TradesParams) error {
	if c == nil {
		return fmt.Errorf("failed to %s trades: client=null", operation)
	}
	key, err := params.Key()
	if err != nil {
		return fmt.Errorf("failed to %s trades: %w", operation, err)
	}
	payload, err := json.Marshal(protocol.SubscribeRequest{Operation: operation, Arguments: []string{key}})
	if err != nil {
		return fmt.Errorf("failed to %s trades: failed to encode request: %w", operation, err)
	}
	if operation == protocol.OperationSubscribe {
		err = c.executor.Subscribe(ctx, key, payload)
	} else {
		err = c.executor.Unsubscribe(ctx, key, payload)
	}
	if err != nil {
		return fmt.Errorf("failed to %s trades: %w", operation, err)
	}
	return nil
}
