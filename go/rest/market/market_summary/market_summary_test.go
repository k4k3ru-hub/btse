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
	return []byte(`[{
		"symbol":"BTC-PERP",
		"base":"BTC",
		"quote":"USDT",
		"active":true,
		"minPriceIncrement":0.1,
		"minSizeIncrement":0.001,
		"fundingIntervalMinutes":480
	}]`), nil
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
		Symbol:             "BTC-PERP",
		ListFullAttributes: true,
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Path != endpoint.FuturesMarketSummaryPath {
		t.Errorf("request path = %q", executor.request.Path)
	}
	if executor.request.Query.Get("symbol") != "BTC-PERP" || executor.request.Query.Get("listFullAttributes") != "true" {
		t.Errorf("request query = %v", executor.request.Query)
	}
	if len(result) != 1 || result[0].Symbol != "BTC-PERP" || result[0].MinPriceIncrement != 0.1 || result[0].FundingIntervalMinutes != 480 {
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
