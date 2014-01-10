package etsys

type ExchangeConnector interface {
	Connect()
	GetTradelog() chan *Trade
	GetOrderlog() chan *OrderState
	SendOrder(*Order)
}

// Root container for all trading infrastructure
type Infrastructure struct {
	Connectors      map[string]ExchangeConnector
	TradeSource     <-chan *Trade
	OrderSource     <-chan *OrderState
	TradeRelay      TradeRelay      // multiplexer
	OrderStateRelay OrderStateRelay // multiplexer
}

type TradeRelay struct{}
type OrderStateRelay struct{}
