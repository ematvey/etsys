package etsys

type BarStream struct {
	source    chan *Bar
	recievers []chan []*Bar
	bars      []*Bar
	cutoff    func([]*Bar) bool
	reduce    func([]*Bar) []*Bar
}

func (bs *BarStream) router() {
	defer func() {
		close(bs.source)
	}()
	for bar := range bs.source {
		bs.bars = append(bs.bars, bar)
		if bs.cutoff(bs.bars) {
			bs.bars = bs.reduce(bs.bars)
		}
		for _, rec := range bs.recievers {
			rec <- bs.bars
		}
	}
}

func (bs *BarStream) AddReciever(rec chan []*Bar) {
	bs.recievers = append(bs.recievers, rec)
}

func (bs *BarStream) Attach(source chan *Bar) {
	bs.source = source
	go bs.router()
}

// --- Factories ---

// Makes BarStream that never reduces itself
func MakeEndlessStream() BarStream {
	return BarStream{
		cutoff: func([]*Bar) bool { return false },
	}
}

// Makes BarStream with length up to fixed size
func MakeLengthStream(length int) BarStream {
	return BarStream{
		cutoff: func(bars []*Bar) bool {
			return len(bars) > length
		},
		reduce: func(bars []*Bar) []*Bar {
			return bars[len(bars)-length:]
		},
	}
}
