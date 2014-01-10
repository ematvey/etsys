package sim

import (
	. "etsys"
	"fmt"
	// "time"
)

// Order Book
type OrderBook struct {
	Ticker     string
	Asks       map[float64]*OrderStack
	Bids       map[float64]*OrderStack
	pipeAdd    <-chan *Order
	pipeCancel <-chan int64
}

func (ob *OrderBook) StartProcessing() {
	ob.Asks = make(map[float64]*OrderStack)
	ob.Bids = make(map[float64]*OrderStack)
	go ob.process()
}

func (ob *OrderBook) process() {
	for {
		select {
		case order := <-ob.pipeAdd:
			ob.addOrder(order)
		case id := <-ob.pipeCancel:
			ob.cancelOrder(id)
		}
	}
}

func (ob *OrderBook) BestBid() (max float64) {
	max = -1
	if len(ob.Bids) == 0 {
		return
	}
	for k := range ob.Bids {
		if k > max {
			max = k
		}
	}
	return
}

func (ob *OrderBook) BestAsk() (min float64) {
	min = -1
	if len(ob.Asks) == 0 {
		return
	}
	for k := range ob.Asks {
		if min == -1 || k < min {
			min = k
		}
	}
	return
}

func (ob *OrderBook) addOrder(order *Order) {
	ba := ob.BestAsk()
	bb := ob.BestBid()
	// if bb > ba {
	// 	panic("bb > ba")
	// }
	if order.IsBuy {
		if ba > 0 && order.Price >= ba {
			ob.executeBuy(order)
		} else {
			ob.putBid(order)
		}
	} else {
		if bb > 0 && order.Price <= bb {
			ob.executeSell(order)
		} else {
			ob.putAsk(order)
		}
	}
}

// tmp plug; rewrite
func (ob *OrderBook) cancelOrder(id int64) {
	panic("not implemented")
}

func (ob *OrderBook) putAsk(order *Order) {
	if ob.Asks[order.Price] == nil {
		ob.Asks[order.Price] = &OrderStack{}
	}
	ob.Asks[order.Price].Add(order)
}

func (ob *OrderBook) putBid(order *Order) {
	if ob.Bids[order.Price] == nil {
		ob.Bids[order.Price] = &OrderStack{}
	}
	ob.Bids[order.Price].Add(order)
}

func (ob *OrderBook) executeSell(order *Order) {
	for order.GetBalance() > 0 {
		bb := ob.BestBid()
		if bb == -1 {
			ob.putAsk(order)
			return
		}
		stack := ob.Bids[bb]
		stack.Execute(order)
		if stack.IsEmpty() {
			delete(ob.Bids, bb)
		}
	}
}

func (ob *OrderBook) executeBuy(order *Order) {
	for order.GetBalance() > 0 {
		ba := ob.BestAsk()
		if ba == -1 {
			ob.putBid(order)
			return
		}
		stack := ob.Asks[ba]
		stack.Execute(order)
		if stack.IsEmpty() {
			delete(ob.Asks, ba)
		}
	}
}

func (ob *OrderBook) String() string {
	return fmt.Sprintf("OrderBook[\n  Ticker: %s\n  BestAsk: %+v\n  BestBid: %+v\n  Asks: %+v\n  Bids: %+v\n]", ob.Ticker, ob.BestAsk(), ob.BestBid(), ob.Asks, ob.Bids)
}
