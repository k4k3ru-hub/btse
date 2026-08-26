// Package ticker24hr implements the BTSE Markets v1 24-hour ticker endpoint.
package ticker24hr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	marketsprotocol "github.com/k4k3ru-hub/btse/go/rest/markets/protocol"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

type Params struct {
	Symbol string
	Full   bool
}

type Ticker struct {
	Symbol                 string `json:"symbol"`
	LastPrice              string `json:"lastPrice"`
	OpenPrice              string `json:"openPrice"`
	HighPrice              string `json:"highPrice"`
	LowPrice               string `json:"lowPrice"`
	Amount                 string `json:"amount"`
	Volume                 string `json:"volume"`
	OpenTime               int64  `json:"openTime"`
	CloseTime              int64  `json:"closeTime"`
	PriceChange            string `json:"priceChange"`
	PriceChangePercent     string `json:"priceChangePercent"`
	PrevClosePrice         string `json:"prevClosePrice"`
	BidPrice               string `json:"bidPrice"`
	BidQuantity            string `json:"bidQty"`
	AskPrice               string `json:"askPrice"`
	AskQuantity            string `json:"askQty"`
	OpenInterest           string `json:"openInterest"`
	FundingRate            string `json:"fundingRate"`
	NextFundingTime        int64  `json:"nextFundingTime"`
	FundingIntervalMinutes int64  `json:"fundingIntervalMinutes"`
}

type Client struct{ executor transport.Executor }

// NewClient creates a Markets 24-hour ticker client.
//
// Version:
//   - 2026-08-26: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create markets 24-hour ticker client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests Markets 24-hour tickers.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) Send(ctx context.Context, params Params) ([]Ticker, error) {
	if c == nil || c.executor == nil {
		return nil, fmt.Errorf("failed to request markets 24-hour tickers: client=null")
	}
	query := url.Values{"isFull": {strconv.FormatBool(params.Full)}}
	if params.Symbol != "" {
		query.Set("symbol", params.Symbol)
	}
	body, err := c.executor.Do(ctx, transport.Request{BaseURL: endpoint.DefaultMarketsBaseURL, Method: http.MethodGet, Path: endpoint.MarketsTicker24HrPath, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request markets 24-hour tickers: %w", err)
	}
	var envelope marketsprotocol.Envelope[[]Ticker]
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("failed to request markets 24-hour tickers: failed to decode response body: %w", err)
	}
	if err := envelope.Validate(); err != nil {
		return nil, fmt.Errorf("failed to request markets 24-hour tickers: %w", err)
	}
	return envelope.Data, nil
}
