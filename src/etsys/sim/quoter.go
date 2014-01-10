package sim

import (
	. "etsys"
	"math"
	"math/rand"
)

// Quote Process
type RandomQuoteProcess struct {
	Ticker    string
	currentId int64
	FvPipe    <-chan float64
	OrderPipe chan<- *Order
	StatePipe chan *OrderState
}

func (qp *RandomQuoteProcess) quoter() {
	for fv := range qp.FvPipe {
		for i := 0; i < rand.Intn(5); i++ {
			vol := float64(rand.Intn(5) + 1)
			qp.currentId++
			order := &Order{
				Ticker:    qp.Ticker,
				Id:        qp.currentId,
				Price:     math.Floor(RndGauss(fv, fv*0.025)),
				Volume:    vol,
				IsBuy:     RndBool(),
				StatePipe: qp.StatePipe,
			}
			order.Init()
			qp.OrderPipe <- order
		}
	}
}

func (qp *RandomQuoteProcess) StartQuoting() {
	go qp.quoter()
}
