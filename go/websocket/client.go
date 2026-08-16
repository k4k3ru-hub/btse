package websocket

import (
	"context"
	"fmt"

	"github.com/k4k3ru-hub/btse/go/websocket/subscriptions"
	k4k3ruWebsocket "github.com/k4k3ru-hub/websocket/go"
)

const (
	DefaultFuturesEndpointURL    = "wss://ws.btse.com/ws/futures"
	DefaultFuturesOSSEndpointURL = "wss://ws.btse.com/ws/oss/futures"
	DefaultSpotEndpointURL       = "wss://ws.btse.com/ws/spot"
	DefaultSpotOSSEndpointURL    = "wss://ws.btse.com/ws/oss/spot"
)

// ClientOption configures the shared WebSocket transport.
type ClientOption = k4k3ruWebsocket.ClientOption

// SessionHandler handles messages and connection closure.
type SessionHandler = k4k3ruWebsocket.SessionHandler

// SessionContext represents a live WebSocket session.
type SessionContext = k4k3ruWebsocket.SessionContext

// Client manages a BTSE WebSocket connection and its subscriptions.
type Client struct {
	wsClient *k4k3ruWebsocket.Client
	bbo      *subscriptions.BBOClient
	trades   *subscriptions.TradesClient
}

// DefaultClientOption returns the shared WebSocket transport defaults.
//
// Version:
//   - 2026-08-16: Added.
func DefaultClientOption() *ClientOption { return k4k3ruWebsocket.DefaultClientOption() }

// NewClient creates a BTSE WebSocket client for the supplied endpoint.
//
// Version:
//   - 2026-08-16: Added.
func NewClient(ctx context.Context, endpointURL string, handler SessionHandler, option *ClientOption) (*Client, error) {
	if endpointURL == "" {
		return nil, fmt.Errorf("failed to create websocket client: endpoint_url=empty")
	}
	if handler == nil {
		return nil, fmt.Errorf("failed to create websocket client: session_handler=null")
	}
	wsClient, err := k4k3ruWebsocket.NewClient(ctx, endpointURL, handler, cloneClientOption(option))
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket client: %w", err)
	}
	client := &Client{wsClient: wsClient}
	client.bbo, err = subscriptions.NewBBOClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket client: %w", err)
	}
	client.trades, err = subscriptions.NewTradesClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create websocket client: %w", err)
	}
	return client, nil
}

func cloneClientOption(option *ClientOption) *ClientOption {
	if option == nil {
		return DefaultClientOption()
	}
	cloned := &ClientOption{
		EndpointURL:      option.EndpointURL,
		ConnectTimeout:   option.ConnectTimeout,
		HandshakeTimeout: option.HandshakeTimeout,
		SessionOption:    option.SessionOption.Clone(),
	}
	if option.HTTPHeader != nil {
		cloned.HTTPHeader = option.HTTPHeader.Clone()
	}
	if cloned.SessionOption == nil {
		cloned.SessionOption = k4k3ruWebsocket.DefaultSessionOption()
	}
	return cloned
}

// Connect opens the WebSocket connection.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Connect(ctx context.Context) error {
	if c == nil || c.wsClient == nil {
		return fmt.Errorf("failed to connect websocket: client=null")
	}
	if err := c.wsClient.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect websocket: %w", err)
	}
	return nil
}

// Close closes the WebSocket client.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Close() error {
	if c == nil || c.wsClient == nil {
		return fmt.Errorf("failed to close websocket: client=null")
	}
	if err := c.wsClient.Close(); err != nil {
		return fmt.Errorf("failed to close websocket: %w", err)
	}
	return nil
}

// BBO returns the BBO subscription client.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) BBO() *subscriptions.BBOClient {
	if c == nil {
		return nil
	}
	return c.bbo
}

// Trades returns the public trade subscription client.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Trades() *subscriptions.TradesClient {
	if c == nil {
		return nil
	}
	return c.trades
}

// SessionContext gets the current WebSocket session context.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) SessionContext() (SessionContext, error) {
	if c == nil || c.wsClient == nil {
		return nil, fmt.Errorf("failed to get websocket session context: client=null")
	}
	return c.wsClient.SessionContext()
}

// Subscribe registers and sends a WebSocket subscription.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Subscribe(ctx context.Context, key string, payload []byte) error {
	return c.wsClient.Subscribe(ctx, key, payload)
}

// Unsubscribe removes and sends a WebSocket subscription.
//
// Version:
//   - 2026-08-16: Added.
func (c *Client) Unsubscribe(ctx context.Context, key string, payload []byte) error {
	return c.wsClient.Unsubscribe(ctx, key, payload)
}
