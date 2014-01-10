package etsys

import (
	"time"
)

// OrderState represents single Order state transition,
// happened at a single instance of time
type OrderState struct {
	Order   *Order
	State   int
	Balance float64
	Time    time.Time
	TimeExt time.Time
	Trade   *Trade
}

func (os *OrderState) IsFill() bool {
	return os.Trade != nil
}

func (os *OrderState) IsActive() bool {
	return os.State == OrderStateFilled || os.State == OrderStateCancelled
}

func (os *OrderState) IsDone() bool {
	return os.State == OrderStateFilled || os.State == OrderStateCancelled
}

func (os *OrderState) IsFilled() bool {
	return os.State == OrderStateFilled
}

func (os *OrderState) IsCancelled() bool {
	return os.State == OrderStateCancelled
}
