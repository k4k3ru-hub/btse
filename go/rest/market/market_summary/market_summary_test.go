package market_summary

import (
	"context"
	"testing"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

type recordingExecutor struct {
	request transport.Request
}

// TestClientSendUsesSpotPath verifies Spot market-summary path selection.
//
// Version:
//   - 2026-08-16: Added.
func TestClientSendUsesSpotPath(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClientWithPath(executor, endpoint.SpotMarketSummaryPath)
	if err != nil {
		t.Fatalf("NewClientWithPath() error = %v", err)
	}
	if _, err := client.Send(context.Background(), Params{Symbol: "BTC-USDT"}); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Path != endpoint.SpotMarketSummaryPath {
		t.Errorf("request path = %q", executor.request.Path)
	}
}

// Do records a market-summary request.
//
// Version:
//   - 2026-08-16: Added.
func (e *recordingExecutor) Do(_ context.Context, request transport.Request) ([]byte, error) {
	e.request = request
	return []byte(`{
		"success":true,
		"code":1,
		"msg":"Success",
		"data":{"symbols":[{
		"symbol":"BTC-PERP",
		"type":"FuturesPerpetual",
		"baseCurrency":"BTC",
		"quoteCurrency":"USDT",
		"active":true,
		"minPriceIncrement":"0.1",
		"minOrderSize":"1",
		"maxOrderSize":"500000",
		"minSizeIncrement":"1",
		"contractSize":"0.001"
	}]}}`), nil
}

// TestClientSendBuildsQueryAndDecodesSummary verifies market-summary requests.
//
// Version:
//   - 2026-08-16: Added.
func TestClientSendBuildsQueryAndDecodesSummary(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.Send(context.Background(), Params{
		Symbol: "BTC-PERP",
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.BaseURL != endpoint.DefaultMarketsBaseURL || executor.request.Path != endpoint.MarketsPath {
		t.Errorf("request path = %q", executor.request.Path)
	}
	if executor.request.Query.Get("symbol") != "BTC-PERP" {
		t.Errorf("request query = %v", executor.request.Query)
	}
	if len(result) != 1 || result[0].Symbol != "BTC-PERP" || result[0].Type != "FuturesPerpetual" || result[0].MinPriceIncrement != 0.1 || result[0].ContractSize != 0.001 {
		t.Errorf("result = %+v", result)
	}
}

// TestClientSendOmitsEmptySymbol verifies requesting summaries for all markets.
//
// Version:
//   - 2026-08-16: Added.
func TestClientSendOmitsEmptySymbol(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if _, err := client.Send(context.Background(), Params{}); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Query.Has("symbol") {
		t.Errorf("request query = %v", executor.request.Query)
	}
}

// TestClientSendRejectsCombinedFilters verifies Markets v1 filter exclusivity.
//
// Version:
//   - 2026-08-18: Added.
func TestClientSendRejectsCombinedFilters(t *testing.T) {
	client, _ := NewClient(&recordingExecutor{})
	if _, err := client.Send(context.Background(), Params{Symbol: "BTC-USDT", Types: []string{"Spot"}}); err == nil {
		t.Fatal("Send() error = nil")
	}
}

// TestMarketSummaryWebSocketSymbol resolves product-specific WebSocket symbols.
//
// Version:
//   - 2026-08-23: Added.
func TestMarketSummaryWebSocketSymbol(t *testing.T) {
	tests := []struct {
		name    string
		summary MarketSummary
		want    string
		wantErr bool
	}{
		{name: "spot", summary: MarketSummary{Type: "Spot", Symbol: "BTC-USDT"}, want: "BTC-USDT"},
		{name: "perpetual", summary: MarketSummary{Type: "FuturesPerpetual", Symbol: "BTC-PERP-USDT", TradeCurrency: "BTC-PERP"}, want: "BTC-PERP"},
		{name: "time based", summary: MarketSummary{Type: "FuturesTimeBased", Symbol: "BTC-260925-USDT", TradeCurrency: "BTC-260925"}, want: "BTC-260925"},
		{name: "missing trade currency", summary: MarketSummary{Type: "FuturesPerpetual", Symbol: "BTC-PERP-USDT"}, wantErr: true},
		{name: "unknown type", summary: MarketSummary{Type: "unknown", Symbol: "BTC-USDT"}, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.summary.WebSocketSymbol()
			if (err != nil) != test.wantErr {
				t.Fatalf("WebSocketSymbol() error = %v, wantErr %t", err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("WebSocketSymbol() = %q, want %q", got, test.want)
			}
		})
	}
}
