package etsys

// Aggregates Ticks into Bars
// Cutoff defines wether given tick should go into this bar or the next
type BarGen struct {
	tickSource chan *Tick
	barPipe    chan *Bar
	currentBar *Bar
	cutoff     func(*Bar, *Tick) bool
}

func (bg *BarGen) pollLoop() {
	defer func() {
		close(bg.tickSource)
		close(bg.barPipe)
	}()
	for tick := range bg.tickSource {
		bg.processTick(tick)
	}
}

func (bg *BarGen) stepCondition(tick *Tick) bool {
	return bg.cutoff(bg.currentBar, tick)
}

func (bg *BarGen) updateBar(tick *Tick) {
	if tick.Price > bg.currentBar.H {
		bg.currentBar.H = tick.Price
	} else if tick.Price < bg.currentBar.L {
		bg.currentBar.L = tick.Price
	}
	bg.currentBar.Ended = tick.Time
	bg.currentBar.C = tick.Price
	bg.currentBar.V += tick.Volume
}

func (bg *BarGen) processTick(tick *Tick) {
	if bg.currentBar == nil {
		bg.currentBar = &Bar{
			Started: tick.Time,
			O:       tick.Price,
			L:       tick.Price,
		}
	} else if bg.stepCondition(tick) {
		bg.barPipe <- bg.currentBar
		bg.currentBar = &Bar{
			Started: tick.Time,
			O:       tick.Price,
			L:       tick.Price,
		}
	}
	bg.updateBar(tick)
}

func (bg *BarGen) Attach(tickchan chan *Tick) chan *Bar {
	bg.tickSource = tickchan
	bg.barPipe = make(chan *Bar)
	go bg.pollLoop()
	return bg.barPipe
}

// --- factories

func MakeVolumeBarGen(volume float64) BarGen {
	return BarGen{
		cutoff: func(b *Bar, t *Tick) bool {
			if b.V+t.Volume > volume {
				return true
			}
			return false
		},
	}
}
