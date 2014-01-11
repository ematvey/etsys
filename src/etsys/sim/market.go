package sim

import (
	. "etsys"
	"time"
)

// Simulated Market
type SimulatedMarket struct {
	Ticker        string
	FairValueProc *RandomWalkProcess
	OrderSource   *RandomQuoteProcess
	OrderBook     *OrderBook
	OrderReciever chan *Order
	Orderlog      chan<- *OrderState
	Tradelog      chan<- *Trade
	tradeId       int64
}

func (m *SimulatedMarket) procOrderStates() {
	for os := range m.OrderSource.StatePipe {
		m.Orderlog <- os
		if os.Trade != nil && os.Trade.BuyMatch != nil && os.Trade.SellMatch != nil {
			m.tradeId++
			os.Trade.Id = m.tradeId
			m.Tradelog <- os.Trade
		}
	}
}

func (m *SimulatedMarket) Run() {
	// run from end to start
	go m.procOrderStates()
	m.OrderBook.StartProcessing()
	m.OrderSource.StartQuoting()
	m.FairValueProc.Walk()
}

func MakeSimulatedMarket(
	ticker string,
	orderlog chan<- *OrderState,
	tradelog chan<- *Trade,
) *SimulatedMarket {
	// pipes initialized explicitly
	fvPipe := make(chan float64)
	orderReciever := make(chan *Order)
	orderStatePipe := make(chan *OrderState)
	cancelPipe := make(chan int64)

	market := &SimulatedMarket{
		Ticker: ticker,
		FairValueProc: &RandomWalkProcess{
			Value:      1000.0,
			Volatility: 0.05,
			Interval:   50 * time.Millisecond,
			Pipe:       fvPipe,
		},
		OrderSource: &RandomQuoteProcess{
			Ticker:    ticker,
			StatePipe: orderStatePipe,
			OrderPipe: orderReciever,
			FvPipe:    fvPipe,
		},
		OrderBook: &OrderBook{
			Ticker:     ticker,
			pipeAdd:    orderReciever,
			pipeCancel: cancelPipe,
		},
		Orderlog:      orderlog,
		Tradelog:      tradelog,
		OrderReciever: orderReciever,
	}

	return market
}
