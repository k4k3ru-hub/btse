package btse

import (
	"context"
	"testing"
)

type testSessionHandler struct{}

// HandleMessage handles a test WebSocket message.
//
// Version:
//   - 2026-08-16: Added.
func (h *testSessionHandler) HandleMessage(WebSocketSessionContext, []byte) {}

// HandleClose handles a test WebSocket closure.
//
// Version:
//   - 2026-08-16: Added.
func (h *testSessionHandler) HandleClose(WebSocketSessionContext) {}

// TestNewRESTClientComposesMarketOperations verifies the root REST facade.
//
// Version:
//   - 2026-08-16: Added.
func TestNewRESTClientComposesMarketOperations(t *testing.T) {
	constructors := []struct {
		constructor       func(*RESTClientOption) (*RESTClient, error)
		hasFundingHistory bool
	}{
		{constructor: NewRESTClient, hasFundingHistory: true},
		{constructor: NewFuturesRESTClient, hasFundingHistory: true},
		{constructor: NewSpotRESTClient},
	}
	for _, test := range constructors {
		client, err := test.constructor(nil)
		if err != nil {
			t.Fatalf("REST constructor error = %v", err)
		}
		if client.Market() == nil || client.Market().Candles() == nil || client.Market().MarketSummary() == nil {
			t.Fatal("REST constructor did not compose market clients")
		}
		if test.hasFundingHistory && client.Market().FundingHistory() == nil {
			t.Fatal("Futures REST constructor did not compose funding history")
		}
		if client.Markets() == nil || client.Markets().Ticker24Hr() == nil || client.Markets().TickerIndices() == nil {
			t.Fatal("REST constructor did not compose Markets v1 clients")
		}
	}
}

// TestProductEndpointConstants verifies the public Spot and Futures endpoints.
//
// Version:
//   - 2026-08-16: Added.
func TestProductEndpointConstants(t *testing.T) {
	if DefaultSpotWebSocketEndpointURL != "wss://ws.btse.com/ws/spot" || DefaultSpotOSSWebSocketEndpointURL != "wss://ws.btse.com/ws/oss/spot" {
		t.Fatal("Spot WebSocket endpoints are invalid")
	}
	if DefaultFuturesWebSocketEndpointURL != "wss://ws.btse.com/ws/futures" || DefaultFuturesOSSWebSocketEndpointURL != "wss://ws.btse.com/ws/oss/futures" {
		t.Fatal("Futures WebSocket endpoints are invalid")
	}
}

// TestNewWebSocketClientComposesSubscriptions verifies the root WebSocket facade.
//
// Version:
//   - 2026-08-16: Added.
func TestNewWebSocketClientComposesSubscriptions(t *testing.T) {
	client, err := NewWebSocketClient(context.Background(), DefaultFuturesWebSocketEndpointURL, &testSessionHandler{}, nil)
	if err != nil {
		t.Fatalf("NewWebSocketClient() error = %v", err)
	}
	if client.BBO() == nil || client.Trades() == nil {
		t.Fatal("NewWebSocketClient() did not compose subscription clients")
	}
}
