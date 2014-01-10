package sim

import (
	"math/rand"
	"time"
)

// Random Walker
type RandomWalkProcess struct {
	Value      float64
	Volatility float64
	Interval   time.Duration
	Pipe       chan float64
}

func (rw *RandomWalkProcess) walker() {
	defer func() { close(rw.Pipe) }()
	for {
		rw.Value *= 1 + (rand.Float64()-0.5)*rw.Volatility
		rw.Pipe <- rw.Value
		time.Sleep(rw.Interval)
	}
}

func (rw *RandomWalkProcess) Walk() {
	go rw.walker()
}
