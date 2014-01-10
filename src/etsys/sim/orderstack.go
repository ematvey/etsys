package sim

import (
	. "etsys"
	"fmt"
	"time"
)

// OrderStack, main purpose - to automate things on the
// OrderBook single price level
type OrderStack struct {
	stack []*Order
}

func (os *OrderStack) Add(o *Order) {
	if os.stack == nil {
		os.stack = make([]*Order, 0)
	}
	os.stack = append(os.stack, o)
	o.SetActive(time.Time{})
}

func (os *OrderStack) Execute(o *Order) {
	if o.IsBuy == os.stack[0].IsBuy {
		panic("same direction on OrderStack and target order")
	}
	if !o.IsActive() {
		o.SetActive(time.Time{})
	}
	for o.GetBalance() > 0 && !os.IsEmpty() {
		var so, buy, sell *Order
		empty := os.clearFilled()
		if empty {
			return
		}
		so = os.stack[0]
		if o.IsBuy {
			buy = o
			sell = so
		} else {
			sell = o
			buy = so
		}
		MatchOrdersInternally(buy, sell)
	}
	os.clearFilled()
}

func (os *OrderStack) clearFilled() (empty bool) {
	for len(os.stack) > 0 && os.stack[0].GetBalance() == 0 {
		if len(os.stack) == 1 {
			os.stack = make([]*Order, 0)
			empty = true
			return
		} else {
			os.stack = os.stack[1:len(os.stack)]
		}
	}
	return
}

func (os *OrderStack) Cancel(id int64) {
	panic("not implemented yet")
}

func (os *OrderStack) IsEmpty() bool {
	return len(os.stack) == 0
}

func (os *OrderStack) Volume() (v float64) {
	for _, o := range os.stack {
		v += o.GetBalance()
	}
	return
}

func (os *OrderStack) String() string {
	return fmt.Sprintf("OrderStack[vol=%v,orders=%v]", os.Volume(), len(os.stack))
}
