// Package tickerindices implements the BTSE Markets v1 index ticker endpoint.
package tickerindices

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	marketsprotocol "github.com/k4k3ru-hub/btse/go/rest/markets/protocol"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

type Params struct{ Symbol string }

type Ticker struct {
	Symbol     string `json:"symbol"`
	IndexPrice string `json:"indexPrice"`
	MarkPrice  string `json:"markPrice"`
}

type Client struct{ executor transport.Executor }

// NewClient creates a Markets index-ticker client.
//
// Version:
//   - 2026-08-26: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create markets index ticker client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests Markets index and mark prices.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) Send(ctx context.Context, params Params) ([]Ticker, error) {
	if c == nil || c.executor == nil {
		return nil, fmt.Errorf("failed to request markets index tickers: client=null")
	}
	query := url.Values{}
	if params.Symbol != "" {
		query.Set("symbol", params.Symbol)
	}
	body, err := c.executor.Do(ctx, transport.Request{BaseURL: endpoint.DefaultMarketsBaseURL, Method: http.MethodGet, Path: endpoint.MarketsTickerIndicesPath, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request markets index tickers: %w", err)
	}
	var envelope marketsprotocol.Envelope[[]Ticker]
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("failed to request markets index tickers: failed to decode response body: %w", err)
	}
	if err := envelope.Validate(); err != nil {
		return nil, fmt.Errorf("failed to request markets index tickers: %w", err)
	}
	return envelope.Data, nil
}
