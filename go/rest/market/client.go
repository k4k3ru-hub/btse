package market

import (
	"fmt"

	"github.com/k4k3ru-hub/btse/go/rest/market/candles"
	"github.com/k4k3ru-hub/btse/go/rest/market/market_summary"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

// Client composes BTSE public market API clients.
type Client struct {
	candles       *candles.Client
	marketSummary *market_summary.Client
}

// NewClient creates a public market API client.
//
// Version:
//   - 2026-08-16: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create market API client: executor=null")
	}
	candlesClient, err := candles.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create market API client: %w", err)
	}
	marketSummaryClient, err := market_summary.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create market API client: %w", err)
	}
	return &Client{candles: candlesClient, marketSummary: marketSummaryClient}, nil
}

// MarketSummary returns the market-summary API client.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) MarketSummary() *market_summary.Client {
	if c == nil {
		return nil
	}
	return c.marketSummary
}

// Candles returns the candle API client.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Candles() *candles.Client {
	if c == nil {
		return nil
	}
	return c.candles
}
