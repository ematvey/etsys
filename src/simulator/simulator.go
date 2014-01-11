package main

import (
	. "etsys"
	. "etsys/sim"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	log.Print("simulation started")

	sigpipe := make(chan os.Signal)
	signal.Notify(sigpipe, os.Interrupt)

	sys := MakeSys()
	sys.Run()

	time.Sleep(time.Millisecond * 50)

	for c := range sys.Connectors {
		log.Printf("[conn] %+v: %+v", c, sys.Connectors[c])
	}

	for {
		select {
		case t := <-sys.Tradelog:
			log.Printf("[t] %+v", t)
		case o := <-sys.Orderlog:
			if o.State == OrderStateCreated {
				// log.Printf("%+v", o)
			}
		case sig := <-sigpipe:
			if sig == os.Interrupt {
				fmt.Println("\nSIGINT\nSTATE DUMP\n")
				log.Printf("%+v", sys.Tickers())
				for _, t := range sys.Tickers() {
					fmt.Println(sys.Connectors[t].StateDump())
				}
				os.Exit(1)
			}
		}
	}
}
