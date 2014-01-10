package etsys

import (
	"time"
)

type Tick struct {
	Id     int64
	Time   time.Time
	Price  float64
	Volume float64
	Ticker string
}
