// Package markets composes BTSE Markets v1 REST operations.
package markets

import (
	"fmt"

	ticker24hr "github.com/k4k3ru-hub/btse/go/rest/markets/ticker_24hr"
	tickerindices "github.com/k4k3ru-hub/btse/go/rest/markets/ticker_indices"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

type Client struct {
	ticker24Hr    *ticker24hr.Client
	tickerIndices *tickerindices.Client
}

// NewClient creates a Markets v1 API client.
//
// Version:
//   - 2026-08-26: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create Markets API client: executor=null")
	}
	ticker24HrClient, err := ticker24hr.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create Markets API client: %w", err)
	}
	tickerIndicesClient, err := tickerindices.NewClient(executor)
	if err != nil {
		return nil, fmt.Errorf("failed to create Markets API client: %w", err)
	}
	return &Client{ticker24Hr: ticker24HrClient, tickerIndices: tickerIndicesClient}, nil
}

// Ticker24Hr returns the 24-hour ticker client.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) Ticker24Hr() *ticker24hr.Client {
	if c == nil {
		return nil
	}
	return c.ticker24Hr
}

// TickerIndices returns the index and mark-price client.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) TickerIndices() *tickerindices.Client {
	if c == nil {
		return nil
	}
	return c.tickerIndices
}
