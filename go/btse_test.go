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
	client, err := NewRESTClient(nil)
	if err != nil {
		t.Fatalf("NewRESTClient() error = %v", err)
	}
	if client.Market() == nil || client.Market().Candles() == nil || client.Market().MarketSummary() == nil {
		t.Fatal("NewRESTClient() did not compose market clients")
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
