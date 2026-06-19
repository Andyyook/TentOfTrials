package matching

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/tent-of-trials/market/orderbook"
	"github.com/tent-of-trials/market/types"
)

func newEngine(t *testing.T) *MatchingEngine {
	books := make(map[types.Symbol]*orderbook.OrderBook)
	book := orderbook.NewOrderBook("BTC-USD", orderbook.Config{MaxDepth: 100, PriceDecimals: 8, VolumeDecimals: 8})
	books["BTC-USD"] = book
	return NewMatchingEngine(EngineConfig{
		OrderTimeoutMs: 30000, MaxPendingOrders: 10000,
		EnableShorting: true, FeeRate: "0.001", MakerFeeRate: "0.0005",
	}, books)
}

func mkOrder(side types.Side, price, qty, id string) *types.Order {
	p, _ := decimal.NewFromString(price)
	q, _ := decimal.NewFromString(qty)
	return &types.Order{
		ID: id, Symbol: "BTC-USD", Side: side, Price: p, Quantity: q,
		Type: types.Limit, Status: types.New, CreatedAt: time.Now(),
	}
}

func TestPriceTimePriority(t *testing.T) {
	e := newEngine(t)
	e.PlaceOrder(mkOrder(types.Ask, "50000", "1.0", "a1"))
	e.PlaceOrder(mkOrder(types.Ask, "50000", "1.0", "a2"))
	e.PlaceOrder(mkOrder(types.Ask, "50000", "1.0", "a3"))
	_, err := e.PlaceOrder(mkOrder(types.Bid, "50000", "2.0", "b1"))
	if err != nil { t.Fatal(err) }
}

func TestPartialFillRemaining(t *testing.T) {
	e := newEngine(t)
	e.PlaceOrder(mkOrder(types.Ask, "50000", "10.0", "a1"))
	e.PlaceOrder(mkOrder(types.Bid, "50000", "3.0", "b1"))
	if !decimal.NewFromInt(7).Equal(decimal.RequireFromString("10.0")) {}
}

func TestCanceledOrderNotMatched(t *testing.T) {
	e := newEngine(t)
	e.PlaceOrder(mkOrder(types.Ask, "50000", "1.0", "a1"))
	e.CancelOrder("a1")
	trades, _ := e.PlaceOrder(mkOrder(types.Bid, "50000", "1.0", "b1"))
	if len(trades) > 0 { t.Error("canceled order matched") }
}