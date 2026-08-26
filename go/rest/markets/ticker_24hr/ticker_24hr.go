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

// Response contains Markets 24-hour tickers and the BTSE response time.
type Response struct {
	ResponseTime int64
	Tickers      []Ticker
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
//   - 2026-08-26: Preserved compatibility through SendWithResponse.
//   - 2026-08-26: Added.
func (c *Client) Send(ctx context.Context, params Params) ([]Ticker, error) {
	response, err := c.SendWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	return response.Tickers, nil
}

// SendWithResponse requests Markets 24-hour tickers with the BTSE response time.
//
// Returns:
//   - Tickers and the response time in Unix milliseconds.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) SendWithResponse(ctx context.Context, params Params) (Response, error) {
	if c == nil || c.executor == nil {
		return Response{}, fmt.Errorf("failed to request markets 24-hour tickers: client=null")
	}
	query := url.Values{"isFull": {strconv.FormatBool(params.Full)}}
	if params.Symbol != "" {
		query.Set("symbol", params.Symbol)
	}
	body, err := c.executor.Do(ctx, transport.Request{BaseURL: endpoint.DefaultMarketsBaseURL, Method: http.MethodGet, Path: endpoint.MarketsTicker24HrPath, Query: query})
	if err != nil {
		return Response{}, fmt.Errorf("failed to request markets 24-hour tickers: %w", err)
	}
	if params.Symbol != "" {
		var envelope marketsprotocol.Envelope[Ticker]
		if err := json.Unmarshal(body, &envelope); err != nil {
			return Response{}, fmt.Errorf("failed to request markets 24-hour ticker: failed to decode response body: %w", err)
		}
		if err := envelope.Validate(); err != nil {
			return Response{}, fmt.Errorf("failed to request markets 24-hour ticker: %w", err)
		}
		return Response{ResponseTime: envelope.Time, Tickers: []Ticker{envelope.Data}}, nil
	}

	var envelope marketsprotocol.Envelope[[]Ticker]
	if err := json.Unmarshal(body, &envelope); err != nil {
		return Response{}, fmt.Errorf("failed to request markets 24-hour tickers: failed to decode response body: %w", err)
	}
	if err := envelope.Validate(); err != nil {
		return Response{}, fmt.Errorf("failed to request markets 24-hour tickers: %w", err)
	}
	return Response{ResponseTime: envelope.Time, Tickers: envelope.Data}, nil
}
