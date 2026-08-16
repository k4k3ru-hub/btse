package rest

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/k4k3ru-hub/btse/go/rest/market/candles"
)

type recordingHTTPClient struct {
	request *http.Request
}

// Do records a REST request and returns a candle response.
//
// Version:
//   - 2026-08-16: Added.
func (c *recordingHTTPClient) Do(request *http.Request) (*http.Response, error) {
	c.request = request
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`[[1624987380,1,2,0.5,1.5,10]]`)),
	}, nil
}

// TestProductClientsUseExpectedCandleURLs verifies product-specific REST composition.
//
// Version:
//   - 2026-08-16: Added.
func TestProductClientsUseExpectedCandleURLs(t *testing.T) {
	tests := []struct {
		name        string
		constructor func(*ClientOption) (*Client, error)
		symbol      string
		wantURL     string
	}{
		{name: "futures", constructor: NewFuturesClient, symbol: "BTC-PERP", wantURL: "https://api.btse.com/futures/api/v2.3/ohlcv"},
		{name: "spot", constructor: NewSpotClient, symbol: "BTC-USDT", wantURL: "https://api.btse.com/spot/api/v3.2/ohlcv"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			httpClient := &recordingHTTPClient{}
			client, err := test.constructor(&ClientOption{HTTPClient: httpClient})
			if err != nil {
				t.Fatalf("constructor error = %v", err)
			}
			_, err = client.Candles(context.Background(), candles.Params{
				Symbol:     test.symbol,
				Resolution: candles.Resolution1Minute,
			})
			if err != nil {
				t.Fatalf("Candles() error = %v", err)
			}
			if httpClient.request.URL.Scheme+"://"+httpClient.request.URL.Host+httpClient.request.URL.Path != test.wantURL {
				t.Errorf("request URL = %q", httpClient.request.URL.String())
			}
		})
	}
}
