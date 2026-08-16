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
