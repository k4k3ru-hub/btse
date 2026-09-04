package protocol

import (
	"encoding/json"
	"testing"
)

// TestOrderBookDecodingPreservesSynchronizationFields verifies order book wire decoding.
//
// Version:
//   - 2026-09-04: Added.
func TestOrderBookDecodingPreservesSynchronizationFields(t *testing.T) {
	message := []byte(`{"topic":"update:BTC-PERP_0","data":{"bids":[["59000.5","2.1"]],"asks":[["59001","0"]],"type":"delta","symbol":"BTC-PERP","timestamp":1788451200000,"seqNum":102,"prevSeqNum":101}}`)

	var orderBook OrderBook
	if err := json.Unmarshal(message, &orderBook); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if orderBook.Topic != "update:BTC-PERP_0" {
		t.Fatalf("Topic = %q", orderBook.Topic)
	}
	if orderBook.Data.Type != "delta" || orderBook.Data.Symbol != "BTC-PERP" {
		t.Fatalf("Data = %+v", orderBook.Data)
	}
	if orderBook.Data.SeqNum != 102 || orderBook.Data.PrevSeqNum != 101 {
		t.Fatalf("sequence = (%d, %d)", orderBook.Data.SeqNum, orderBook.Data.PrevSeqNum)
	}
	if len(orderBook.Data.Bids) != 1 || orderBook.Data.Bids[0] != [2]string{"59000.5", "2.1"} {
		t.Fatalf("Bids = %+v", orderBook.Data.Bids)
	}
	if len(orderBook.Data.Asks) != 1 || orderBook.Data.Asks[0] != [2]string{"59001", "0"} {
		t.Fatalf("Asks = %+v", orderBook.Data.Asks)
	}
}
