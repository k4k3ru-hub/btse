package tickerindices

import (
	"context"
	"testing"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

type recordingExecutor struct{ request transport.Request }

func (e *recordingExecutor) Do(_ context.Context, request transport.Request) ([]byte, error) {
	e.request = request
	return []byte(`{"success":true,"code":1,"msg":"Success","time":1760325052630,"data":[{"symbol":"BTC-PERP-USDT","indexPrice":"117995.0","markPrice":"118001.0"}]}`), nil
}

// TestSendDecodesMarkAndIndexPrices verifies complete index-ticker decoding.
//
// Version:
//   - 2026-08-26: Added.
func TestSendDecodesMarkAndIndexPrices(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	result, err := client.Send(context.Background(), Params{Symbol: "BTC-PERP-USDT"})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Path != endpoint.MarketsTickerIndicesPath || executor.request.Query.Get("symbol") != "BTC-PERP-USDT" {
		t.Fatalf("request = %#v", executor.request)
	}
	if len(result) != 1 || result[0].IndexPrice != "117995.0" || result[0].MarkPrice != "118001.0" {
		t.Fatalf("result = %#v", result)
	}
}
