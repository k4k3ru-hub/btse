package subscriptions

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/k4k3ru-hub/btse/go/websocket/protocol"
)

type recordingExecutor struct {
	subscriptions map[string]protocol.SubscribeRequest
	unsubscribes  map[string]protocol.SubscribeRequest
}

// Subscribe records a subscription.
//
// Version:
//   - 2026-08-16: Added.
func (e *recordingExecutor) Subscribe(_ context.Context, key string, payload []byte) error {
	var r protocol.SubscribeRequest
	if err := json.Unmarshal(payload, &r); err != nil {
		return err
	}
	e.subscriptions[key] = r
	return nil
}

// Unsubscribe records an unsubscription.
//
// Version:
//   - 2026-08-16: Added.
func (e *recordingExecutor) Unsubscribe(_ context.Context, key string, payload []byte) error {
	var r protocol.SubscribeRequest
	if err := json.Unmarshal(payload, &r); err != nil {
		return err
	}
	e.unsubscribes[key] = r
	return nil
}

// TestBBOClientBuildsBTSEPayload verifies BBO subscription payloads.
//
// Version:
//   - 2026-08-16: Added.
func TestBBOClientBuildsBTSEPayload(t *testing.T) {
	e := &recordingExecutor{subscriptions: map[string]protocol.SubscribeRequest{}, unsubscribes: map[string]protocol.SubscribeRequest{}}
	c, _ := NewBBOClient(e)
	p := BBOParams{Symbol: "BTC-PERP"}
	if err := c.Subscribe(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	if err := c.Unsubscribe(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	if r := e.subscriptions["snapshotL1:BTC-PERP"]; r.Operation != "subscribe" || len(r.Arguments) != 1 || r.Arguments[0] != "snapshotL1:BTC-PERP" {
		t.Errorf("subscribe = %+v", r)
	}
	if r := e.unsubscribes["snapshotL1:BTC-PERP"]; r.Operation != "unsubscribe" {
		t.Errorf("unsubscribe = %+v", r)
	}
}

// TestOrderBookClientBuildsBTSEPayload verifies order book subscription payloads.
//
// Version:
//   - 2026-09-04: Added.
func TestOrderBookClientBuildsBTSEPayload(t *testing.T) {
	e := &recordingExecutor{subscriptions: map[string]protocol.SubscribeRequest{}, unsubscribes: map[string]protocol.SubscribeRequest{}}
	c, err := NewOrderBookClient(e)
	if err != nil {
		t.Fatal(err)
	}
	p := OrderBookParams{Symbol: "BTC-PERP"}
	if err := c.Subscribe(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	if err := c.Unsubscribe(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	const key = "update:BTC-PERP_0"
	if r := e.subscriptions[key]; r.Operation != "subscribe" || len(r.Arguments) != 1 || r.Arguments[0] != key {
		t.Errorf("subscribe = %+v", r)
	}
	if r := e.unsubscribes[key]; r.Operation != "unsubscribe" || len(r.Arguments) != 1 || r.Arguments[0] != key {
		t.Errorf("unsubscribe = %+v", r)
	}
	if _, err := (OrderBookParams{}).Key(); err == nil {
		t.Fatal("Key() error = nil, want empty symbol error")
	}
}

// TestTradesClientBuildsBTSEPayload verifies trade subscription payloads.
//
// Version:
//   - 2026-08-16: Added.
func TestTradesClientBuildsBTSEPayload(t *testing.T) {
	e := &recordingExecutor{subscriptions: map[string]protocol.SubscribeRequest{}, unsubscribes: map[string]protocol.SubscribeRequest{}}
	c, _ := NewTradesClient(e)
	p := TradesParams{Symbol: "BTC-PERP"}
	if err := c.Subscribe(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	if r := e.subscriptions["tradeHistoryApiV3:BTC-PERP"]; r.Operation != "subscribe" || r.Arguments[0] != "tradeHistoryApiV3:BTC-PERP" {
		t.Errorf("subscribe = %+v", r)
	}
}

// TestTradesParamsSelectsMarketTopic verifies market-specific trade topics.
//
// Version:
//   - 2026-08-23: Added.
func TestTradesParamsSelectsMarketTopic(t *testing.T) {
	tests := []struct {
		name   string
		params TradesParams
		want   string
	}{
		{name: "spot", params: TradesParams{Symbol: "BTC-USDT", Market: TradeMarketSpot}, want: "tradeHistoryApi:BTC-USDT"},
		{name: "futures", params: TradesParams{Symbol: "BTC-PERP", Market: TradeMarketFutures}, want: "tradeHistoryApiV3:BTC-PERP"},
		{name: "compatibility default", params: TradesParams{Symbol: "BTC-PERP"}, want: "tradeHistoryApiV3:BTC-PERP"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.params.Key()
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("Key() = %q, want %q", got, test.want)
			}
		})
	}
	if _, err := (TradesParams{Symbol: "BTC-USDT", Market: "invalid"}).Key(); err == nil {
		t.Fatal("Key() error = nil, want invalid market error")
	}
}
