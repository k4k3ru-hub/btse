package market_summary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

// Params configures a market-summary request.
type Params struct {
	Symbol  string
	Symbols []string
	Types   []string
	// ListFullAttributes is retained for source compatibility and is ignored by Markets v1.
	ListFullAttributes bool
}

// MarketSummary contains BTSE market metadata.
type MarketSummary struct {
	Symbol                 string   `json:"symbol"`
	Type                   string   `json:"type"`
	Category               string   `json:"category"`
	TradeCurrency          string   `json:"tradeCurrency"`
	Last                   float64  `json:"last"`
	LowestAsk              float64  `json:"lowestAsk"`
	HighestBid             float64  `json:"highestBid"`
	PercentageChange       float64  `json:"percentageChange"`
	Volume                 float64  `json:"volume"`
	High24Hr               float64  `json:"high24Hr"`
	Low24Hr                float64  `json:"low24Hr"`
	Base                   string   `json:"baseCurrency"`
	Quote                  string   `json:"quoteCurrency"`
	Active                 bool     `json:"active"`
	Size                   float64  `json:"size"`
	MinValidPrice          float64  `json:"-"`
	MinOrderPrice          float64  `json:"minOrderPrice,string"`
	MinPriceIncrement      float64  `json:"minPriceIncrement,string"`
	PricePrecision         int      `json:"pricePrecision"`
	MinOrderSize           float64  `json:"minOrderSize,string"`
	MaxOrderSize           float64  `json:"maxOrderSize,string"`
	MinSizeIncrement       float64  `json:"minSizeIncrement,string"`
	SizePrecision          int      `json:"sizePrecision"`
	OpenInterest           float64  `json:"openInterest"`
	OpenInterestUSD        float64  `json:"openInterestUSD"`
	ContractStart          int64    `json:"contractStartTime"`
	ContractEnd            int64    `json:"contractEndTime"`
	TimeBasedContract      bool     `json:"timeBasedContract"`
	OpenTime               int64    `json:"openTime"`
	CloseTime              int64    `json:"closeTime"`
	StartMatching          int64    `json:"matchingStartTime"`
	InactiveTime           int64    `json:"becomeInactiveTime"`
	FundingRate            float64  `json:"fundingRate"`
	ContractSize           float64  `json:"contractSize,string"`
	MaxPosition            float64  `json:"maxPosition"`
	MinRiskLimit           float64  `json:"minRiskLimit"`
	MaxRiskLimit           float64  `json:"maxRiskLimit"`
	AvailableSettlement    []string `json:"availableSettlement"`
	Futures                bool     `json:"futures"`
	FundingIntervalMinutes int      `json:"fundingIntervalMinutes"`
	FundingTime            int64    `json:"fundingTime"`
}

// WebSocketSymbol resolves the venue symbol used by BTSE WebSocket feeds.
//
// Returns:
//   - WebSocket venue symbol.
//   - Resolution error.
//
// Version:
//   - 2026-08-23: Added.
func (m MarketSummary) WebSocketSymbol() (string, error) {
	var symbol string
	switch strings.TrimSpace(m.Type) {
	case "Spot":
		symbol = strings.TrimSpace(m.Symbol)
	case "FuturesPerpetual", "FuturesTimeBased":
		symbol = strings.TrimSpace(m.TradeCurrency)
	default:
		return "", fmt.Errorf("failed to resolve market summary websocket symbol: market_type=invalid")
	}
	if symbol == "" {
		return "", fmt.Errorf("failed to resolve market summary websocket symbol: symbol=empty")
	}
	return symbol, nil
}

// Client requests BTSE market summaries.
type Client struct {
	executor transport.Executor
	baseURL  string
	path     string
}

// NewClient creates a market-summary client.
//
// Version:
//   - 2026-08-18: Changed the default endpoint to the unified Markets v1 API.
//   - 2026-08-16: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	return NewClientWithEndpoint(executor, endpoint.DefaultMarketsBaseURL, endpoint.MarketsPath)
}

// NewClientWithPath creates a market-summary client for a product-specific API path.
//
// Version:
//   - 2026-08-18: Added Markets v1 filters and response decoding.
//   - 2026-08-16: Added.
func NewClientWithPath(executor transport.Executor, path string) (*Client, error) {
	return NewClientWithEndpoint(executor, "", path)
}

// NewClientWithEndpoint creates a market-information client for an API endpoint.
//
// Version:
//   - 2026-08-18: Added.
func NewClientWithEndpoint(executor transport.Executor, baseURL, path string) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create market summary client: executor=null")
	}
	if path == "" {
		return nil, fmt.Errorf("failed to create market summary client: path=empty")
	}
	return &Client{executor: executor, baseURL: baseURL, path: path}, nil
}

// Send requests market summaries.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Send(ctx context.Context, params Params) ([]MarketSummary, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request market summary: client=null")
	}

	parameterCount := 0
	if params.Symbol != "" {
		parameterCount++
	}
	if len(params.Symbols) > 0 {
		parameterCount++
	}
	if len(params.Types) > 0 {
		parameterCount++
	}
	if parameterCount > 1 {
		return nil, fmt.Errorf("failed to request market summary: filters=invalid")
	}

	query := make(url.Values)
	if params.Symbol != "" {
		query.Set("symbol", params.Symbol)
	}
	if len(params.Symbols) > 0 {
		value, err := json.Marshal(params.Symbols)
		if err != nil {
			return nil, fmt.Errorf("failed to request market summary: failed to encode symbols: %w", err)
		}
		query.Set("symbols", string(value))
	}
	if len(params.Types) > 0 {
		value, err := json.Marshal(params.Types)
		if err != nil {
			return nil, fmt.Errorf("failed to request market summary: failed to encode types: %w", err)
		}
		query.Set("types", string(value))
	}

	body, err := c.executor.Do(ctx, transport.Request{
		BaseURL: c.baseURL,
		Method:  http.MethodGet,
		Path:    c.path,
		Query:   query,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to request market summary: %w", err)
	}

	var response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"msg"`
		Data    struct {
			Symbols []MarketSummary `json:"symbols"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to request market summary: failed to decode response body: %w", err)
	}
	if !response.Success {
		return nil, fmt.Errorf("failed to request market summary: api request unsuccessful: code=%d message=%q", response.Code, response.Message)
	}
	return response.Data.Symbols, nil
}
