package etsys

import (
	_ "log"
)

type ExchangeConnector interface {
	Connect()
	GetTradelog() chan *Trade
	GetOrderlog() chan *OrderState
	SendOrder(*Order)
	GetTickers() []string
	StateDump() string
}

// Root container for all trading infrastructure
type Infrastructure struct {
	Connectors map[string]ExchangeConnector
	Tradelog   chan *Trade
	Orderlog   chan *OrderState
}

func (I *Infrastructure) Run() {
	for c := range I.Connectors {
		I.Connectors[c].Connect()
	}
}

func (I *Infrastructure) Tickers() (tickers []string) {
	tickers = make([]string, 0)
	for t := range I.Connectors {
		tickers = append(tickers, t)
	}
	return
}
