package protocol

import (
	"encoding/json"
	"fmt"
)

const (
	OperationSubscribe   = "subscribe"
	OperationUnsubscribe = "unsubscribe"
	TopicBBO             = "snapshotL1"
	TopicOrderBook       = "update"
	TopicSpotTrades      = "tradeHistoryApi"
	TopicFuturesTrades   = "tradeHistoryApiV3"
	TopicTrades          = TopicFuturesTrades
)

// SubscribeRequest is the common BTSE public subscription request.
type SubscribeRequest struct {
	Operation string   `json:"op"`
	Arguments []string `json:"args"`
}

// Envelope is the common outer format for BTSE market messages.
type Envelope struct {
	Topic string          `json:"topic,omitempty"`
	Event string          `json:"event,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
}

// BuildKey returns the topic used as the transport subscription key.
//
// Version:
//   - 2026-08-16: Added.
func (e *Envelope) BuildKey() (string, error) {
	if e == nil {
		return "", fmt.Errorf("failed to build websocket message key: envelope=null")
	}
	if e.Topic == "" {
		return "", fmt.Errorf("failed to build websocket message key: topic=empty")
	}
	return e.Topic, nil
}

// BBO is a best-bid/best-ask snapshot message.
type BBO struct {
	Topic string  `json:"topic"`
	Data  BBOData `json:"data"`
}

// BBOData contains a level-one order book snapshot.
type BBOData struct {
	Bids      [][2]string `json:"bids"`
	Asks      [][2]string `json:"asks"`
	Type      string      `json:"type"`
	Symbol    string      `json:"symbol"`
	Timestamp int64       `json:"timestamp"`
}

// OrderBook is a level-two order book snapshot or delta message.
type OrderBook struct {
	Topic string        `json:"topic"`
	Data  OrderBookData `json:"data"`
}

// OrderBookData contains absolute price-level sizes for an order book update.
type OrderBookData struct {
	Bids       [][2]string `json:"bids"`
	Asks       [][2]string `json:"asks"`
	Type       string      `json:"type"`
	Symbol     string      `json:"symbol"`
	Timestamp  int64       `json:"timestamp"`
	SeqNum     int64       `json:"seqNum"`
	PrevSeqNum int64       `json:"prevSeqNum"`
}

// Trades is a public trade message.
type Trades struct {
	Topic string  `json:"topic"`
	Data  []Trade `json:"data"`
}

// Trade contains one public fill.
type Trade struct {
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Size      float64 `json:"size"`
	Price     float64 `json:"price"`
	TradeID   int64   `json:"tradeId"`
	Timestamp int64   `json:"timestamp"`
}
