package candles

import (
	"context"
	"testing"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

type recordingExecutor struct{ request transport.Request }

// Do records a candle request.
//
// Version:
//   - 2026-08-16: Added.
func (e *recordingExecutor) Do(_ context.Context, request transport.Request) ([]byte, error) {
	e.request = request
	return []byte(`[[1624987380,36477,36478,36470,36473.5,693.049]]`), nil
}

// TestClientSendBuildsQueryAndDecodesCandles verifies the public candle operation.
//
// Version:
//   - 2026-08-16: Added.
func TestClientSendBuildsQueryAndDecodesCandles(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	result, err := client.Send(context.Background(), Params{Symbol: "BTC-PERP", Resolution: Resolution1Minute, Start: 1000, End: 2000})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Query.Get("symbol") != "BTC-PERP" || executor.request.Query.Get("resolution") != "1" || executor.request.Query.Get("start") != "1000" || executor.request.Query.Get("end") != "2000" {
		t.Errorf("query = %v", executor.request.Query)
	}
	if len(result) != 1 || result[0].Timestamp != 1624987380 || result[0].Close != 36473.5 || result[0].Volume != 693.049 {
		t.Errorf("result = %+v", result)
	}
}

// TestClientSendUsesSpotPath verifies Spot candle path selection.
//
// Version:
//   - 2026-08-16: Added.
func TestClientSendUsesSpotPath(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClientWithPath(executor, endpoint.SpotOHLCVPath)
	if err != nil {
		t.Fatalf("NewClientWithPath() error = %v", err)
	}
	if _, err := client.Send(context.Background(), Params{Symbol: "BTC-USDT", Resolution: Resolution1Minute}); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Path != endpoint.SpotOHLCVPath {
		t.Errorf("request path = %q", executor.request.Path)
	}
}

// TestClientSendRejectsInvalidParams verifies candle request validation.
//
// Version:
//   - 2026-08-16: Added.
func TestClientSendRejectsInvalidParams(t *testing.T) {
	client, _ := NewClient(&recordingExecutor{})
	for _, params := range []Params{{Resolution: Resolution1Minute}, {Symbol: "BTC-PERP"}, {Symbol: "BTC-PERP", Resolution: Resolution1Minute, Start: 2, End: 1}} {
		if _, err := client.Send(context.Background(), params); err == nil {
			t.Errorf("Send(%+v) error = nil", params)
		}
	}
}
