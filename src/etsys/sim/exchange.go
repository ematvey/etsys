package sim

import (
	. "etsys"
)

type SimulatedExchange struct {
	markets  map[string]*SimulatedMarket
	Tradelog chan *Trade
	Orderlog chan *OrderState
}

func (se *SimulatedExchange) Run() {
	for m := range se.markets {
		se.markets[m].Run()
	}
}

func (se *SimulatedExchange) AttachMarket(m *SimulatedMarket) {
	se.markets[m.Ticker] = m
}

func (se *SimulatedExchange) GetMarket(ticker string) *SimulatedMarket {
	return se.markets[ticker]
}
func (se *SimulatedExchange) GetTickers() []string {
	ts := make([]string, 0)
	for t := range se.markets {
		ts = append(ts, t)
	}
	return ts
}

func MakeSimulatedExchange() *SimulatedExchange {
	se := &SimulatedExchange{
		markets:  make(map[string]*SimulatedMarket),
		Tradelog: make(chan *Trade),
		Orderlog: make(chan *OrderState),
	}
	return se
}

func MakeSomeSimulatedExchange() *SimulatedExchange {
	se := MakeSimulatedExchange()
	se.AttachMarket(MakeSimulatedMarket("A", se.Orderlog, se.Tradelog))
	se.AttachMarket(MakeSimulatedMarket("B", se.Orderlog, se.Tradelog))
	return se
}
