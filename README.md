# btse

Go SDK for BTSE public Futures market data.

## Supported data

- BBO: `snapshotL1:<symbol>` on the Futures OSS WebSocket endpoint
- Trades: `tradeHistoryApiV3:<symbol>` on the Futures WebSocket endpoint
- Candles: `GET /api/v2.3/ohlcv` on the Futures REST endpoint

The WebSocket feeds use separate BTSE endpoints, so create one client for BBO
and another for trades.

```go
restClient, err := btse.NewRESTClient(nil)
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

Import path: `github.com/k4k3ru-hub/btse/go`.
