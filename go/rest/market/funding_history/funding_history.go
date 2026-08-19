// Package funding_history provides access to settled BTSE funding rates.
package funding_history

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

// Params configures a funding-history request.
type Params struct {
	Symbol             string
	Count              int
	From               int64
	To                 int64
	UseNewSymbolNaming bool
}

// FundingRate represents one settled funding-rate record.
type FundingRate struct {
	Symbol string  `json:"symbol"`
	Time   int64   `json:"time"`
	Rate   float64 `json:"rate"`
}

// Client requests settled BTSE funding rates.
type Client struct {
	executor transport.Executor
	path     string
}

// NewClient creates a funding-history client.
//
// Version:
//   - 2026-08-19: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	return NewClientWithPath(executor, endpoint.FuturesFundingHistoryPath)
}

// NewClientWithPath creates a funding-history client for an API path.
//
// Version:
//   - 2026-08-19: Added.
func NewClientWithPath(executor transport.Executor, path string) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create funding history client: executor=null")
	}
	if path == "" {
		return nil, fmt.Errorf("failed to create funding history client: path=empty")
	}
	return &Client{executor: executor, path: path}, nil
}

// Send requests settled funding rates in symbol and time order.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) Send(ctx context.Context, params Params) ([]FundingRate, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request funding history: client=null")
	}
	if params.Count < 0 {
		return nil, fmt.Errorf("failed to request funding history: count=out_of_range min_value=0")
	}
	if params.From < 0 {
		return nil, fmt.Errorf("failed to request funding history: from=out_of_range min_value=0")
	}
	if params.To < 0 {
		return nil, fmt.Errorf("failed to request funding history: to=out_of_range min_value=0")
	}
	if params.Count > 0 && (params.From > 0 || params.To > 0) {
		return nil, fmt.Errorf("failed to request funding history: filters=invalid")
	}
	if params.From > 0 && params.To > 0 && params.From > params.To {
		return nil, fmt.Errorf("failed to request funding history: time_range=invalid")
	}

	query := make(url.Values)
	if params.Symbol != "" {
		query.Set("symbol", params.Symbol)
	}
	if params.Count > 0 {
		query.Set("count", strconv.Itoa(params.Count))
	}
	if params.From > 0 {
		query.Set("from", strconv.FormatInt(params.From, 10))
	}
	if params.To > 0 {
		query.Set("to", strconv.FormatInt(params.To, 10))
	}
	if params.UseNewSymbolNaming {
		query.Set("useNewSymbolNaming", "true")
	}

	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: c.path, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request funding history: %w", err)
	}
	var response map[string][]FundingRate
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to request funding history: failed to decode response body: %w", err)
	}

	result := make([]FundingRate, 0)
	for _, rates := range response {
		result = append(result, rates...)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Symbol == result[j].Symbol {
			return result[i].Time < result[j].Time
		}
		return result[i].Symbol < result[j].Symbol
	})
	return result, nil
}
