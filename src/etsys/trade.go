package etsys

import (
	"fmt"
	"time"
)

type Trade struct {
	Id           int64
	ExtId        int64
	Ticker       string
	Time         time.Time
	TimeExt      time.Time
	Price        float64
	Volume       float64
	BuyMatch     *OrderState
	SellMatch    *OrderState
	BuyInitiated bool
}

func (t *Trade) String() string {
	return fmt.Sprintf("Trade{%s id:%v price:%v vol:%v buy:%v %v}",
		t.Ticker, t.Id, t.Price, t.Volume, t.BuyInitiated, t.Time)
}

// Match two orders and produce a trade.
// Order state updates could be extrated from Trade pointer later on.
func MatchOrdersInternally(buy *Order, sell *Order) *Trade {
	buyState := buy.GetState()
	sellState := sell.GetState()

	// Sanity checks
	if buyState.State == OrderStateCreated {
		panic("buy order incorect state: created")
	} else if buyState.State == OrderStateFilled {
		panic("buy order incorect state: filled")
	} else if buyState.State == OrderStateCancelled {
		panic("buy order incorect state: cancelled")
	} else if sellState.State == OrderStateCreated {
		panic("sell order incorect state: created")
	} else if sellState.State == OrderStateFilled {
		panic("sell order incorect state: filled")
	} else if sellState.State == OrderStateCancelled {
		panic("sell order incorect state: cancelled")
	}
	if buy.Ticker != sell.Ticker {
		panic("cant match orders with different tickers")
	}

	// Main procedure
	var extTime time.Time
	var buyInit bool
	var matchVolume float64
	var price float64

	// Get match time
	if buyState.TimeExt.After(sellState.TimeExt) {
		extTime = buyState.TimeExt
	} else {
		extTime = sellState.TimeExt
	}

	// Get init order
	if buyState.TimeExt != sellState.TimeExt {
		if buyState.TimeExt.After(sellState.TimeExt) {
			buyInit = true
		} else {
			buyInit = false
		}
	} else {
		if buyState.Time.After(sellState.Time) {
			buyInit = true
		} else {
			buyInit = false
		}
	}

	// Price
	if buyInit {
		price = sell.Price
	} else {
		price = buy.Price
	}

	// Volume matched
	if sellState.Balance > buyState.Balance {
		matchVolume = buyState.Balance
	} else {
		matchVolume = sellState.Balance
	}

	trade := &Trade{
		Ticker:       buy.Ticker,
		Time:         time.Now(),
		TimeExt:      extTime,
		BuyInitiated: buyInit,
		Volume:       matchVolume,
		Price:        price,
	}
	sell.RecordTrade(trade)
	buy.RecordTrade(trade)

	return trade
}
