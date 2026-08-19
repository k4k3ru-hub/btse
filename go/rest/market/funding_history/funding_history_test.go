package funding_history

import (
	"context"
	"testing"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

type recordingExecutor struct{ request transport.Request }

// Do records a funding-history request.
//
// Version:
//   - 2026-08-19: Added.
func (e *recordingExecutor) Do(_ context.Context, request transport.Request) ([]byte, error) {
	e.request = request
	return []byte(`{"ETH-PERP":[{"time":20,"rate":0.0002,"symbol":"ETH-PERP"}],"BTC-PERP":[{"time":10,"rate":0.0001,"symbol":"BTC-PERP"}]}`), nil
}

// TestClientSendBuildsQueryAndDecodesHistory verifies funding-history requests.
//
// Version:
//   - 2026-08-19: Added.
func TestClientSendBuildsQueryAndDecodesHistory(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	result, err := client.Send(context.Background(), Params{Symbol: "BTC-PERP", From: 1000, To: 2000, UseNewSymbolNaming: true})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Path != endpoint.FuturesFundingHistoryPath || executor.request.Query.Get("symbol") != "BTC-PERP" || executor.request.Query.Get("from") != "1000" || executor.request.Query.Get("to") != "2000" || executor.request.Query.Get("useNewSymbolNaming") != "true" {
		t.Errorf("request = %+v", executor.request)
	}
	if len(result) != 2 || result[0].Symbol != "BTC-PERP" || result[0].Time != 10 || result[0].Rate != 0.0001 {
		t.Errorf("result = %+v", result)
	}
}

// TestClientSendRejectsInvalidParams verifies funding-history validation.
//
// Version:
//   - 2026-08-19: Added.
func TestClientSendRejectsInvalidParams(t *testing.T) {
	client, _ := NewClient(&recordingExecutor{})
	for _, params := range []Params{{Count: -1}, {From: -1}, {To: -1}, {Count: 1, From: 1}, {From: 2, To: 1}} {
		if _, err := client.Send(context.Background(), params); err == nil {
			t.Errorf("Send(%+v) error = nil", params)
		}
	}
}
