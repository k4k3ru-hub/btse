package rest

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	"github.com/k4k3ru-hub/btse/go/rest/market"
	"github.com/k4k3ru-hub/btse/go/rest/market/candles"
	"github.com/k4k3ru-hub/btse/go/rest/market/funding_history"
	"github.com/k4k3ru-hub/btse/go/rest/market/market_summary"
	"github.com/k4k3ru-hub/btse/go/rest/markets"
	ticker24hr "github.com/k4k3ru-hub/btse/go/rest/markets/ticker_24hr"
	tickerindices "github.com/k4k3ru-hub/btse/go/rest/markets/ticker_indices"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

// HTTPClient executes HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ClientOption configures a REST client.
type ClientOption struct {
	BaseURL        string
	ConnectTimeout time.Duration
	HTTPClient     HTTPClient
}

// Client executes REST requests and provides composed API clients.
type Client struct {
	baseURL    string
	httpClient HTTPClient
	market     *market.Client
	markets    *markets.Client
}

// ResponseError represents a non-successful HTTP response.
type ResponseError struct {
	StatusCode int
	Body       string
}

// Error returns the REST response error message.
//
// Version:
//   - 2026-08-16: Added.
func (e *ResponseError) Error() string {
	if e == nil {
		return "failed to execute rest request: response_error=null"
	}
	return fmt.Sprintf("failed to execute rest request: unexpected HTTP status: status_code=%d", e.StatusCode)
}

// DefaultClientOption returns the default Futures REST options.
//
// Version:
//   - 2026-08-16: Added.
func DefaultClientOption() *ClientOption {
	return &ClientOption{BaseURL: endpoint.DefaultFuturesBaseURL, ConnectTimeout: 3 * time.Second}
}

// DefaultSpotClientOption returns the default Spot REST options.
//
// Version:
//   - 2026-08-16: Added.
func DefaultSpotClientOption() *ClientOption {
	return &ClientOption{BaseURL: endpoint.DefaultSpotBaseURL, ConnectTimeout: 3 * time.Second}
}

// NewClient creates a REST client and composes its API clients.
//
// Version:
//   - 2026-08-19: Composed settled Futures funding history.
//   - 2026-08-18: Composed market metadata with the unified Markets v1 API.
//   - 2026-08-16: Added.
func NewClient(option *ClientOption) (*Client, error) {
	return newClient(option, DefaultClientOption(), market.Paths{
		Candles:        endpoint.FuturesOHLCVPath,
		MarketSummary:  endpoint.MarketsPath,
		FundingHistory: endpoint.FuturesFundingHistoryPath,
	})
}

// NewFuturesClient creates a Futures REST client.
//
// Version:
//   - 2026-08-19: Composed settled Futures funding history.
//   - 2026-08-18: Composed market metadata with the unified Markets v1 API.
//   - 2026-08-16: Added.
func NewFuturesClient(option *ClientOption) (*Client, error) {
	return NewClient(option)
}

// NewSpotClient creates a Spot REST client.
//
// Version:
//   - 2026-08-18: Composed market metadata with the unified Markets v1 API.
//   - 2026-08-16: Added.
func NewSpotClient(option *ClientOption) (*Client, error) {
	return newClient(option, DefaultSpotClientOption(), market.Paths{
		Candles:       endpoint.SpotOHLCVPath,
		MarketSummary: endpoint.MarketsPath,
	})
}

func newClient(option *ClientOption, defaults *ClientOption, paths market.Paths) (*Client, error) {
	if option == nil {
		option = defaults
	}
	baseURL := strings.TrimRight(option.BaseURL, "/")
	if baseURL == "" {
		baseURL = defaults.BaseURL
	}
	httpClient := option.HTTPClient
	if httpClient == nil {
		timeout := option.ConnectTimeout
		if timeout == 0 {
			timeout = defaults.ConnectTimeout
		}
		if timeout < 0 {
			return nil, fmt.Errorf("failed to create rest client: connect_timeout=out_of_range")
		}
		httpClient = &http.Client{Transport: &http.Transport{DialContext: (&net.Dialer{Timeout: timeout}).DialContext}}
	}
	client := &Client{baseURL: baseURL, httpClient: httpClient}
	marketClient, err := market.NewClientWithPaths(client, paths)
	if err != nil {
		return nil, fmt.Errorf("failed to create rest client: %w", err)
	}
	client.market = marketClient
	marketsClient, err := markets.NewClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create rest client: %w", err)
	}
	client.markets = marketsClient
	return client, nil
}

// Market returns the composed public market API client.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Market() *market.Client {
	if c == nil {
		return nil
	}
	return c.market
}

// Markets returns the composed Markets v1 API client.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) Markets() *markets.Client {
	if c == nil {
		return nil
	}
	return c.markets
}

// Candles requests candle data.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Candles(ctx context.Context, params candles.Params) ([]candles.Candle, error) {
	if c == nil || c.market == nil || c.market.Candles() == nil {
		return nil, fmt.Errorf("failed to request candles: client=null")
	}
	return c.market.Candles().Send(ctx, params)
}

// MarketSummary requests market metadata and current summary values.
//
// Version:
//   - 2026-08-18: Changed the operation to return Markets v1 metadata.
//   - 2026-08-16: Added.
func (c *Client) MarketSummary(ctx context.Context, params market_summary.Params) ([]market_summary.MarketSummary, error) {
	if c == nil || c.market == nil || c.market.MarketSummary() == nil {
		return nil, fmt.Errorf("failed to request market summary: client=null")
	}
	return c.market.MarketSummary().Send(ctx, params)
}

// MarketSummaryWithResponse requests market metadata with the BTSE response time.
//
// Returns:
//   - Market summaries and the response time in Unix milliseconds.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) MarketSummaryWithResponse(ctx context.Context, params market_summary.Params) (market_summary.Response, error) {
	if c == nil || c.market == nil || c.market.MarketSummary() == nil {
		return market_summary.Response{}, fmt.Errorf("failed to request market summary: client=null")
	}
	return c.market.MarketSummary().SendWithResponse(ctx, params)
}

// FundingHistory requests settled funding rates.
//
// Version:
//   - 2026-08-19: Added.
func (c *Client) FundingHistory(ctx context.Context, params funding_history.Params) ([]funding_history.FundingRate, error) {
	if c == nil || c.market == nil || c.market.FundingHistory() == nil {
		return nil, fmt.Errorf("failed to request funding history: funding_history_client=null")
	}
	return c.market.FundingHistory().Send(ctx, params)
}

// Ticker24Hr requests Markets v1 tickers, including current funding data.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) Ticker24Hr(ctx context.Context, params ticker24hr.Params) ([]ticker24hr.Ticker, error) {
	if c == nil || c.markets == nil || c.markets.Ticker24Hr() == nil {
		return nil, fmt.Errorf("failed to request markets 24-hour tickers: client=null")
	}
	return c.markets.Ticker24Hr().Send(ctx, params)
}

// Ticker24HrWithResponse requests Markets v1 tickers with the BTSE response time.
//
// Returns:
//   - Tickers and the response time in Unix milliseconds.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) Ticker24HrWithResponse(ctx context.Context, params ticker24hr.Params) (ticker24hr.Response, error) {
	if c == nil || c.markets == nil || c.markets.Ticker24Hr() == nil {
		return ticker24hr.Response{}, fmt.Errorf("failed to request markets 24-hour tickers: client=null")
	}
	return c.markets.Ticker24Hr().SendWithResponse(ctx, params)
}

// TickerIndices requests Markets v1 index and mark prices.
//
// Version:
//   - 2026-08-26: Added.
func (c *Client) TickerIndices(ctx context.Context, params tickerindices.Params) ([]tickerindices.Ticker, error) {
	if c == nil || c.markets == nil || c.markets.TickerIndices() == nil {
		return nil, fmt.Errorf("failed to request markets index tickers: client=null")
	}
	return c.markets.TickerIndices().Send(ctx, params)
}

// Do executes a REST request.
//
// Version:
//   - 2026-08-18: Added support for request-specific base URLs.
//   - 2026-08-16: Added.
func (c *Client) Do(ctx context.Context, request transport.Request) ([]byte, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to execute rest request: client=null")
	}
	if request.Method == "" {
		return nil, fmt.Errorf("failed to execute rest request: method=empty")
	}
	if request.Path == "" {
		return nil, fmt.Errorf("failed to execute rest request: path=empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	baseURL := c.baseURL
	if request.BaseURL != "" {
		baseURL = strings.TrimRight(request.BaseURL, "/")
	}
	httpRequest, err := http.NewRequestWithContext(ctx, request.Method, baseURL+"/"+strings.TrimLeft(request.Path, "/"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rest request: failed to create HTTP request: %w", err)
	}
	httpRequest.URL.RawQuery = request.Query.Encode()
	httpRequest.Header = request.Header.Clone()
	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rest request: %w", err)
	}
	if response == nil {
		return nil, fmt.Errorf("failed to execute rest request: response=null")
	}
	if response.Body == nil {
		return nil, fmt.Errorf("failed to execute rest request: response_body=null")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rest request: failed to read response body: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, &ResponseError{StatusCode: response.StatusCode, Body: string(body)}
	}
	return body, nil
}
