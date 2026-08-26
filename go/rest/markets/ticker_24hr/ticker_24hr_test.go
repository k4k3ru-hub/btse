package ticker24hr

import (
	"context"
	"testing"

	"github.com/k4k3ru-hub/btse/go/rest/endpoint"
	"github.com/k4k3ru-hub/btse/go/rest/transport"
)

type recordingExecutor struct{ request transport.Request }

func (e *recordingExecutor) Do(_ context.Context, request transport.Request) ([]byte, error) {
	e.request = request
	return []byte(`{"success":true,"code":1,"msg":"Success","time":1760325052630,"data":{"symbol":"BTC-PERP-USDT","lastPrice":"118000.1","openPrice":"116000.0","highPrice":"119000.0","lowPrice":"115500.0","amount":"25000000.0","volume":"212.5","openTime":1760238652630,"closeTime":1760325052630,"priceChange":"2000.1","priceChangePercent":"1.7242","prevClosePrice":"115999.9","bidPrice":"118000.0","bidQty":"1.2","askPrice":"118000.2","askQty":"1.1","openInterest":"12345.6","fundingRate":"0.0001","nextFundingTime":1760342400000,"fundingIntervalMinutes":480}}`), nil
}

func TestSendDecodesFullFundingPayload(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	result, err := client.Send(context.Background(), Params{Symbol: "BTC-PERP-USDT", Full: true})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if executor.request.Path != endpoint.MarketsTicker24HrPath || executor.request.Query.Get("symbol") != "BTC-PERP-USDT" || executor.request.Query.Get("isFull") != "true" {
		t.Fatalf("request = %#v", executor.request)
	}
	if len(result) != 1 || result[0].Symbol != "BTC-PERP-USDT" || result[0].LastPrice != "118000.1" || result[0].OpenPrice != "116000.0" || result[0].HighPrice != "119000.0" || result[0].LowPrice != "115500.0" || result[0].Amount != "25000000.0" || result[0].Volume != "212.5" || result[0].OpenTime != 1760238652630 || result[0].CloseTime != 1760325052630 || result[0].PriceChange != "2000.1" || result[0].PriceChangePercent != "1.7242" || result[0].PrevClosePrice != "115999.9" || result[0].BidPrice != "118000.0" || result[0].BidQuantity != "1.2" || result[0].AskPrice != "118000.2" || result[0].AskQuantity != "1.1" || result[0].FundingRate != "0.0001" || result[0].NextFundingTime != 1760342400000 || result[0].FundingIntervalMinutes != 480 || result[0].OpenInterest != "12345.6" {
		t.Fatalf("result = %#v", result)
	}
}

type allTickersExecutor struct{}

func (e *allTickersExecutor) Do(_ context.Context, _ transport.Request) ([]byte, error) {
	return []byte(`{"success":true,"code":1,"msg":"Success","time":1760325052630,"data":[{"symbol":"BTC-PERP-USDT","openInterest":"12345.6"}]}`), nil
}

// TestSendWithResponseDecodesTickerArray verifies the all-symbol response shape.
//
// Version:
//   - 2026-08-26: Added.
func TestSendWithResponseDecodesTickerArray(t *testing.T) {
	client, err := NewClient(&allTickersExecutor{})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	response, err := client.SendWithResponse(context.Background(), Params{Full: true})
	if err != nil {
		t.Fatalf("SendWithResponse() error = %v", err)
	}
	if response.ResponseTime != 1760325052630 || len(response.Tickers) != 1 || response.Tickers[0].OpenInterest != "12345.6" {
		t.Fatalf("response = %#v", response)
	}
}

// TestSendWithResponsePreservesResponseTime verifies the BTSE observation time is exposed.
//
// Version:
//   - 2026-08-26: Added.
func TestSendWithResponsePreservesResponseTime(t *testing.T) {
	executor := &recordingExecutor{}
	client, err := NewClient(executor)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	response, err := client.SendWithResponse(context.Background(), Params{Symbol: "BTC-PERP-USDT", Full: true})
	if err != nil {
		t.Fatalf("SendWithResponse() error = %v", err)
	}
	if response.ResponseTime != 1760325052630 || len(response.Tickers) != 1 || response.Tickers[0].OpenInterest != "12345.6" {
		t.Fatalf("response = %#v", response)
	}
}
