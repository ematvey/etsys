package main

import (
	. "etsys"
	"fmt"
)

func main() {
	ts := NewTickSource()
	bg := MakeVolumeBarGen(3.0)
	tickchan := make(chan *Tick)
	barchan := bg.Attach(tickchan)
	bs := MakeEndlessStream()
	bsc := make(chan []*Bar)
	bs.AddReciever(bsc)
	bs.Attach(barchan)
	go ts.Generate(tickchan)
	for {
		x := <-bsc
		y := x[len(x)-1]
		fmt.Printf("%v\t%v\n", y.Ended, y.C)
		// <-bsc
	}
}
