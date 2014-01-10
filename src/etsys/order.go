package etsys

import (
	"time"
)

const (
	OrderStateCreated = iota
	OrderStateSent
	OrderStateActive
	OrderStateFilled
	OrderStateCancelled
)

// Order represents Exchange order
type Order struct {
	Ticker        string
	TickerId      int64
	TickerIdExt   int64
	ExtId         int64
	Id            int64
	Price         float64
	Volume        float64
	IsBuy         bool
	StateSequence []*OrderState
	StatePipe     chan *OrderState
}

// --- Setters ---

func (o *Order) changeState(state int, balance float64, extTime time.Time, trade *Trade) {
	newState := &OrderState{
		Order:   o,
		Balance: balance,
		Time:    time.Now(),
		TimeExt: extTime,
		State:   state,
		Trade:   trade,
	}
	o.StateSequence = append(
		o.StateSequence,
		newState,
	)
	o.StatePipe <- newState
}

func (o *Order) SetActive(extTime time.Time) {
	o.changeState(OrderStateActive, o.Volume, extTime, nil)
}

func (o *Order) Init() {
	o.StateSequence = make([]*OrderState, 0, 10)
	o.changeState(OrderStateCreated, o.Volume, time.Time{}, nil)
}

func (o *Order) RecordTrade(t *Trade) {
	os := o.GetState()
	if t.Volume > os.Balance {
		panic("trade volume > order balance")
	}
	var ns int
	if t.Volume == os.Balance {
		ns = OrderStateFilled
	} else {
		ns = OrderStateActive
	}
	o.changeState(ns, os.Balance-t.Volume, t.TimeExt, t)
	state := o.GetState()
	if t != state.Trade {
		panic("very severe error")
	}
	if o.IsBuy {
		t.BuyMatch = state
	} else {
		t.SellMatch = state
	}
}

// --- Getters ---
func (o *Order) GetState() *OrderState {
	return o.StateSequence[len(o.StateSequence)-1]
}
func (o *Order) GetBalance() float64 {
	return o.GetState().Balance
}

// --- State Getters ---
func (o *Order) IsPartiallyFilled() bool {
	return o.GetState().Balance < o.Volume
}
func (o *Order) IsActive() bool {
	return o.GetState().IsActive()
}
func (o *Order) IsDone() bool {
	return o.GetState().IsDone()
}
func (o *Order) IsFilled() bool {
	return o.GetState().IsFilled()
}
func (o *Order) IsCancelled() bool {
	return o.GetState().IsCancelled()
}
