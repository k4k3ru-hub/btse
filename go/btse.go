// Package btse provides the root facade for the BTSE SDK.
package btse

import (
	"context"

	"github.com/k4k3ru-hub/btse/go/rest"
	"github.com/k4k3ru-hub/btse/go/rest/market/candles"
	"github.com/k4k3ru-hub/btse/go/rest/market/market_summary"
	btseWebsocket "github.com/k4k3ru-hub/btse/go/websocket"
	"github.com/k4k3ru-hub/btse/go/websocket/protocol"
	"github.com/k4k3ru-hub/btse/go/websocket/subscriptions"
)

const (
	DefaultFuturesWebSocketEndpointURL    = btseWebsocket.DefaultFuturesEndpointURL
	DefaultFuturesOSSWebSocketEndpointURL = btseWebsocket.DefaultFuturesOSSEndpointURL
)

type RESTClient = rest.Client
type RESTClientOption = rest.ClientOption
type RESTResponseError = rest.ResponseError
type Candle = candles.Candle
type CandleParams = candles.Params
type CandleResolution = candles.Resolution
type MarketSummary = market_summary.MarketSummary
type MarketSummaryParams = market_summary.Params

const (
	CandleResolution1Minute  = candles.Resolution1Minute
	CandleResolution5Minute  = candles.Resolution5Minute
	CandleResolution15Minute = candles.Resolution15Minute
	CandleResolution30Minute = candles.Resolution30Minute
	CandleResolution1Hour    = candles.Resolution1Hour
	CandleResolution4Hour    = candles.Resolution4Hour
	CandleResolution6Hour    = candles.Resolution6Hour
	CandleResolution1Day     = candles.Resolution1Day
	CandleResolution1Week    = candles.Resolution1Week
	CandleResolution1Month   = candles.Resolution1Month
)

type WebSocketClient = btseWebsocket.Client
type WebSocketClientOption = btseWebsocket.ClientOption
type WebSocketSessionHandler = btseWebsocket.SessionHandler
type WebSocketSessionContext = btseWebsocket.SessionContext
type BBOParams = subscriptions.BBOParams
type TradesParams = subscriptions.TradesParams
type WebSocketEnvelope = protocol.Envelope
type BBO = protocol.BBO
type BBOData = protocol.BBOData
type Trades = protocol.Trades
type Trade = protocol.Trade

// DefaultRESTClientOption returns the default Futures REST options.
//
// Version:
//   - 2026-08-16: Added.
func DefaultRESTClientOption() *RESTClientOption { return rest.DefaultClientOption() }

// NewRESTClient creates a Futures REST client.
//
// Version:
//   - 2026-08-16: Added.
func NewRESTClient(option *RESTClientOption) (*RESTClient, error) { return rest.NewClient(option) }

// DefaultWebSocketClientOption returns the WebSocket transport defaults.
//
// Version:
//   - 2026-08-16: Added.
func DefaultWebSocketClientOption() *WebSocketClientOption {
	return btseWebsocket.DefaultClientOption()
}

// NewWebSocketClient creates a WebSocket client for a BTSE endpoint.
//
// Version:
//   - 2026-08-16: Added.
func NewWebSocketClient(ctx context.Context, endpointURL string, handler WebSocketSessionHandler, option *WebSocketClientOption) (*WebSocketClient, error) {
	return btseWebsocket.NewClient(ctx, endpointURL, handler, option)
}
