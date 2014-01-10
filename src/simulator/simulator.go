package main

import (
	. "etsys"
	. "etsys/sim"
	"fmt"
	"os"
	"os/signal"
)

func main() {

	fmt.Println("SIMULATION STARTED")

	sigpipe := make(chan os.Signal)
	signal.Notify(sigpipe, os.Interrupt)

	ex := MakeSomeSimulatedExchange()
	ex.Run()

	for {
		select {
		case t := <-ex.Tradelog:
			fmt.Printf("%+v\n", t)
		case o := <-ex.Orderlog:
			if o.State == OrderStateCreated {
				fmt.Printf("order: %+v\n", o)
			}
		case sig := <-sigpipe:
			if sig == os.Interrupt {
				fmt.Println("\nSIGINT\nSTATE DUMP\n")
				for _, t := range ex.GetTickers() {
					fmt.Println(ex.GetMarket(t).OrderBook)
				}
				os.Exit(1)
			}
		}
	}
}
