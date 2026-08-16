# btse

Go SDK for BTSE public Spot and Futures market data.

## Supported data

- BBO: `snapshotL1:<symbol>` on the Spot or Futures OSS WebSocket endpoint
- Trades: `tradeHistoryApiV3:<symbol>` on the Spot or Futures WebSocket endpoint
- Candles: Spot `GET /api/v3.2/ohlcv`; Futures `GET /api/v2.3/ohlcv`
- Market summary: Spot `GET /api/v3.2/market_summary`; Futures `GET /api/v2.3/market_summary`

The WebSocket feeds use separate BTSE endpoints, so create one client for BBO
and another for trades.

```go
restClient, err := btse.NewFuturesRESTClient(nil)
if err != nil {
	return err
}

spotClient, err := btse.NewSpotRESTClient(nil)
if err != nil {
	return err
}
spotCandles, err := spotClient.Candles(ctx, btse.CandleParams{
	Symbol:     "BTC-USDT",
	Resolution: btse.CandleResolution1Minute,
})
if err != nil {
	return err
}
candles, err := restClient.Candles(ctx, btse.CandleParams{
	Symbol:     "BTC-PERP",
	Resolution: btse.CandleResolution1Minute,
})
if err != nil {
	return err
}
summaries, err := restClient.MarketSummary(ctx, btse.MarketSummaryParams{
	Symbol:             "BTC-PERP",
	ListFullAttributes: true,
})

bboClient, err := btse.NewWebSocketClient(
	ctx,
	btse.DefaultFuturesOSSWebSocketEndpointURL,
	handler,
	nil,
)
if err != nil {
	return err
}
if err := bboClient.Connect(ctx); err != nil {
	return err
}
if err := bboClient.BBO().Subscribe(ctx, btse.BBOParams{Symbol: "BTC-PERP"}); err != nil {
	return err
}
```

`NewRESTClient` remains an alias for the Futures client for compatibility. Use
the Spot WebSocket constants with `NewWebSocketClient` for Spot subscriptions.

Import path: `github.com/k4k3ru-hub/btse/go`.
