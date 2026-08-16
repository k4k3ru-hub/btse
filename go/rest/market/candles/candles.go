package candles

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

// Resolution identifies a supported candle interval in minutes.
type Resolution string

const (
	Resolution1Minute  Resolution = "1"
	Resolution5Minute  Resolution = "5"
	Resolution15Minute Resolution = "15"
	Resolution30Minute Resolution = "30"
	Resolution1Hour    Resolution = "60"
	Resolution4Hour    Resolution = "240"
	Resolution6Hour    Resolution = "360"
	Resolution1Day     Resolution = "1440"
	Resolution1Week    Resolution = "10080"
	Resolution1Month   Resolution = "43200"
)

// Params configures a candle request.
type Params struct {
	Symbol     string
	Resolution Resolution
	Start      int64
	End        int64
}

// Candle represents one BTSE OHLCV row.
type Candle struct {
	Timestamp int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// UnmarshalJSON decodes the positional BTSE OHLCV representation.
//
// Version:
//   - 2026-08-16: Added.
func (c *Candle) UnmarshalJSON(data []byte) error {
	if c == nil {
		return fmt.Errorf("failed to decode candle: candle=null")
	}

	var values []json.Number
	if err := json.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("failed to decode candle: %w", err)
	}
	if len(values) != 6 {
		return fmt.Errorf("failed to decode candle: values=invalid actual_length=%d expected_length=6", len(values))
	}

	timestamp, err := values[0].Int64()
	if err != nil {
		return fmt.Errorf("failed to decode candle: failed to decode timestamp: %w", err)
	}
	numbers := make([]float64, 5)
	for i := range numbers {
		numbers[i], err = values[i+1].Float64()
		if err != nil {
			return fmt.Errorf("failed to decode candle: failed to decode value: %w: index=%d", err, i+1)
		}
	}

	c.Timestamp = timestamp
	c.Open = numbers[0]
	c.High = numbers[1]
	c.Low = numbers[2]
	c.Close = numbers[3]
	c.Volume = numbers[4]
	return nil
}

// Client requests BTSE candle data.
type Client struct {
	executor transport.Executor
	path     string
}

// NewClient creates a candle client.
//
// Version:
//   - 2026-08-16: Added.
func NewClient(executor transport.Executor) (*Client, error) {
	return NewClientWithPath(executor, endpoint.FuturesOHLCVPath)
}

// NewClientWithPath creates a candle client for a product-specific API path.
//
// Version:
//   - 2026-08-16: Added.
func NewClientWithPath(executor transport.Executor, path string) (*Client, error) {
	if executor == nil {
		return nil, fmt.Errorf("failed to create candles client: executor=null")
	}
	if path == "" {
		return nil, fmt.Errorf("failed to create candles client: path=empty")
	}
	return &Client{executor: executor, path: path}, nil
}

// Send requests candle data.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Send(ctx context.Context, params Params) ([]Candle, error) {
	if c == nil {
		return nil, fmt.Errorf("failed to request candles: client=null")
	}
	if params.Symbol == "" {
		return nil, fmt.Errorf("failed to request candles: symbol=empty")
	}
	if !params.Resolution.valid() {
		return nil, fmt.Errorf("failed to request candles: resolution=invalid")
	}
	if params.Start < 0 {
		return nil, fmt.Errorf("failed to request candles: start=out_of_range min_value=0")
	}
	if params.End < 0 {
		return nil, fmt.Errorf("failed to request candles: end=out_of_range min_value=0")
	}
	if params.Start > 0 && params.End > 0 && params.Start > params.End {
		return nil, fmt.Errorf("failed to request candles: time_range=invalid")
	}

	query := url.Values{
		"symbol":     {params.Symbol},
		"resolution": {string(params.Resolution)},
	}
	if params.Start > 0 {
		query.Set("start", strconv.FormatInt(params.Start, 10))
	}
	if params.End > 0 {
		query.Set("end", strconv.FormatInt(params.End, 10))
	}

	body, err := c.executor.Do(ctx, transport.Request{Method: http.MethodGet, Path: c.path, Query: query})
	if err != nil {
		return nil, fmt.Errorf("failed to request candles: %w", err)
	}
	var result []Candle
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to request candles: failed to decode response body: %w", err)
	}
	return result, nil
}

func (r Resolution) valid() bool {
	switch r {
	case Resolution1Minute, Resolution5Minute, Resolution15Minute, Resolution30Minute,
		Resolution1Hour, Resolution4Hour, Resolution6Hour, Resolution1Day,
		Resolution1Week, Resolution1Month:
		return true
	default:
		return false
	}
}
