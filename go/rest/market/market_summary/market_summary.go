package market_summary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

// Params configures a market-summary request.
type Params struct {
	Symbol             string
	ListFullAttributes bool
}

// MarketSummary contains BTSE Futures market metadata and current summary values.
type MarketSummary struct {
	Symbol                 string   `json:"symbol"`
	Last                   float64  `json:"last"`
	LowestAsk              float64  `json:"lowestAsk"`
	HighestBid             float64  `json:"highestBid"`
	PercentageChange       float64  `json:"percentageChange"`
	Volume                 float64  `json:"volume"`
	High24Hr               float64  `json:"high24Hr"`
	Low24Hr                float64  `json:"low24Hr"`
	Base                   string   `json:"base"`
	Quote                  string   `json:"quote"`
	Active                 bool     `json:"active"`
	Size                   float64  `json:"size"`
	MinValidPrice          float64  `json:"minValidPrice"`
	MinPriceIncrement      float64  `json:"minPriceIncrement"`
	MinOrderSize           float64  `json:"minOrderSize"`
	MaxOrderSize           float64  `json:"maxOrderSize"`
	MinSizeIncrement       float64  `json:"minSizeIncrement"`
	OpenInterest           float64  `json:"openInterest"`
	OpenInterestUSD        float64  `json:"openInterestUSD"`
	ContractStart          int64    `json:"contractStart"`
	ContractEnd            int64    `json:"contractEnd"`
	TimeBasedContract      bool     `json:"timeBasedContract"`
	OpenTime               int64    `json:"openTime"`
	CloseTime              int64    `json:"closeTime"`
	StartMatching          int64    `json:"startMatching"`
	InactiveTime           int64    `json:"inactiveTime"`
	FundingRate            float64  `json:"fundingRate"`
	ContractSize           float64  `json:"contractSize"`
	MaxPosition            float64  `json:"maxPosition"`
	MinRiskLimit           float64  `json:"minRiskLimit"`
	MaxRiskLimit           float64  `json:"maxRiskLimit"`
	AvailableSettlement    []string `json:"availableSettlement"`
	Futures                bool     `json:"futures"`
	FundingIntervalMinutes int      `json:"fundingIntervalMinutes"`
	FundingTime            int64    `json:"fundingTime"`
}

// Client requests BTSE market summaries.
type Client struct {
	executor transport.Executor
}

// NewClient creates a market-summary client.
//
// Version:
//   - 2026-08-16: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create market summary client: executor=null")
	}
	return &Client{executor: executor}, nil
}

// Send requests market summaries.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Send(ctx context.Context, params Params) ([]MarketSummary, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request market summary: client=null")
	}

	query := make(url.Values)
	if params.Symbol != "" {
		query.Set("symbol", params.Symbol)
	}
	query.Set("listFullAttributes", strconv.FormatBool(params.ListFullAttributes))

	body, err := c.executor.Do(ctx, transport.Request{
		Method: http.MethodGet,
		Path:   endpoint.FuturesMarketSummaryPath,
		Query:  query,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to request market summary: %w", err)
	}

	var result []MarketSummary
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request market summary: failed to decode response body: %w", err)
	}
	return result, nil
}
