# btse

Go SDK for BTSE public Spot and Futures market data.

## Supported data

- BBO: `snapshotL1:<symbol>` on the Spot or Futures OSS WebSocket endpoint
- Spot Trades: `tradeHistoryApi:<symbol>` on the Spot WebSocket endpoint
- Futures Trades: `tradeHistoryApiV3:<symbol>` on the Futures WebSocket endpoint
- Candles: Spot `GET /api/v3.2/ohlcv`; Futures `GET /api/v2.3/ohlcv`
- Market metadata: unified Markets v1 `GET /public-api/market/v1/markets`
- Settled Futures funding rates: `GET /api/v2.1/funding_history`
- Current funding rate, next funding time, and open interest: Markets v1 `GET /public-api/market/v1/ticker/24hr`
- Mark price and index price: Markets v1 `GET /public-api/market/v1/ticker/indices`

BTSE Futures Open Interest is currently obtained by REST polling; the public
WebSocket feeds do not document an Open Interest stream. The Markets v1
24-hour ticker provides `openInterest`, which is the venue-reported number of
open positions and is treated as contracts. The Markets v1 market summary
separately provides `contractSize` and `baseCurrency`, so a base asset quantity
can be derived as `openInterest * contractSize`. The unified Markets v1
responses do not currently provide an Open Interest notional value; consumers
that calculate one should preserve the conversion price and its timestamp.

Use `MarketSummaryWithResponse` or `Ticker24HrWithResponse` when the BTSE
response time is required as the observation timestamp. `ResponseTime` is a
Unix timestamp in milliseconds. The existing methods without `WithResponse`
remain available for compatibility.

The WebSocket feeds use separate BTSE endpoints, so create one client for BBO
and another for trades.
Funding-history parameters are owned by
`github.com/k4k3ru-hub/btse/go/rest/market/funding_history`.

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
	Symbol: "BTC-PERP-USDT",
})

fundingRates, err := restClient.FundingHistory(ctx, funding_history.Params{
	Symbol: "BTC-PERP",
	From:   startTime.UnixMilli(),
	To:     endTime.UnixMilli(),
})

currentFunding, err := restClient.Ticker24Hr(ctx, btse.Ticker24HrParams{
	Symbol: "BTC-PERP-USDT",
	Full:   true,
})
if err != nil {
	return err
}

openInterest, err := restClient.MarketSummaryWithResponse(ctx, btse.MarketSummaryParams{
	Symbol: "BTC-PERP-USDT",
})
if err != nil {
	return err
}
_ = openInterest.ResponseTime
_ = openInterest.MarketSummaries
indices, err := restClient.TickerIndices(ctx, btse.TickerIndicesParams{
	Symbol: "BTC-PERP-USDT",
})
if err != nil {
	return err
}

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

Set `TradesParams.Market` so the SDK selects the protocol topic. The zero value
retains the historical Futures `tradeHistoryApiV3` behavior for compatibility.
Use `MarketSummary.WebSocketSymbol` to resolve unified Markets v1 symbols such
as `BTC-PERP-USDT` to the Futures WebSocket symbol `BTC-PERP`.

`NewRESTClient` remains an alias for the Futures client for compatibility. Use
the Spot WebSocket constants with `NewWebSocketClient` for Spot subscriptions.

Import path: `github.com/k4k3ru-hub/btse/go`.
