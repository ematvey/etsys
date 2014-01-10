package etsys

import (
	"math/rand"
	"time"
)

type TickRandomSource struct {
	currentId     int64
	currentTime   time.Time
	last          float64
	maxVolatility float64
}

func NewTickSource() TickRandomSource {
	return TickRandomSource{
		currentId:     1,
		currentTime:   time.Now(),
		last:          1000.0,
		maxVolatility: 10.0,
	}
}

func (ts *TickRandomSource) Tick() *Tick {
	p := ts.last + (rand.Float64()-0.5)*ts.maxVolatility
	now := time.Now()
	b := &Tick{
		Id:     ts.currentId,
		Time:   now,
		Price:  p,
		Volume: 1,
	}
	ts.currentId++
	ts.currentTime = now
	ts.last = p
	return b
}

func (ts *TickRandomSource) Generate(tickchan chan *Tick) {
	defer func() { close(tickchan) }()
	for {
		tickchan <- ts.Tick()
		time.Sleep(100 * time.Millisecond)
	}
}
