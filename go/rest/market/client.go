package market

import (
	"fmt"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	"github.com/k4k3ru-hub/btse/go/rest/market/candles"
	"github.com/k4k3ru-hub/btse/go/rest/market/funding_history"
	"github.com/k4k3ru-hub/btse/go/rest/market/market_summary"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

// Client composes BTSE public market API clients.
type Client struct {
	candles        *candles.Client
	marketSummary  *market_summary.Client
	fundingHistory *funding_history.Client
}

// Paths configures product-specific public market API paths.
type Paths struct {
	Candles        string
	MarketSummary  string
	FundingHistory string
}

// NewClient creates a public market API client.
//
// Version:
//   - 2026-08-19: Composed settled Futures funding history.
//   - 2026-08-16: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	return NewClientWithPaths(executor, Paths{
		Candles:        endpoint.FuturesOHLCVPath,
		MarketSummary:  endpoint.MarketsPath,
		FundingHistory: endpoint.FuturesFundingHistoryPath,
	})
}

// NewClientWithPaths creates a public market API client for product-specific paths.
//
// Version:
//   - 2026-08-19: Composed settled Futures funding history when configured.
//   - 2026-08-18: Composed market metadata with the unified Markets v1 API.
//   - 2026-08-16: Added.
func NewClientWithPaths(executor transport.Executor, paths Paths) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create market API client: executor=null")
	}
	candlesClient, err := candles.NewClientWithPath(executor, paths.Candles)
	if err != nil {
		return nil, fmt.Errorf("failed to create market API client: %w", err)
	}
	marketSummaryClient, err := market_summary.NewClientWithEndpoint(executor, endpoint.DefaultMarketsBaseURL, paths.MarketSummary)
	if err != nil {
		return nil, fmt.Errorf("failed to create market API client: %w", err)
	}
	var fundingHistoryClient *funding_history.Client
	if paths.FundingHistory != "" {
		fundingHistoryClient, err = funding_history.NewClientWithPath(executor, paths.FundingHistory)
		if err != nil {
			return nil, fmt.Errorf("failed to create market API client: %w", err)
		}
	}
	return &Client{candles: candlesClient, marketSummary: marketSummaryClient, fundingHistory: fundingHistoryClient}, nil
}

// FundingHistory returns the settled funding-history API client.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) FundingHistory() *funding_history.Client {
	if c == nil {
		return nil
	}
	return c.fundingHistory
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
