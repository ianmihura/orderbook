package main

import (
	"fmt"
	"time"
)

func PrintDisplay(ob *OrderBook, stop <-chan bool) {
	for {
		select {
		case <-stop:
			fmt.Print("\033[H\033[2J")
			return
		default:
			time.Sleep(time.Millisecond * 100)
			fmt.Print("\033[H\033[2J")
			ob.PPrint()
		}
	}
}

func PrintHelp() {
	fmt.Println("")
	fmt.Println("List of commands:")
	fmt.Println("  o, orderbook: prints orderbook")
	fmt.Println("  m, mid: prints midprice")
	fmt.Println("  n, new: create new order")
	fmt.Println("  c, clear: clear the console")
	fmt.Println("  d, display: watch the ob go by")
	fmt.Println("  r, reset: gives new initial values to orderbook")
	fmt.Println("  q, quit: quits the program")
}

func PrintNewOrderHelp() {
	fmt.Println("Order format: A|B,LIMIT|MARKET,size,price")
	fmt.Println("Where: A|B is the side Ask or Bid")
	fmt.Println("Where: L|M is the type Limit or Market")
	fmt.Println("  Note that Market orders will ignore price")
	fmt.Println("Where: Size in an int")
	fmt.Println("Where: Price is a float")
	fmt.Println("* Values must be separated by a single comma (, char)")
	fmt.Println()
	fmt.Println("Some examples...")
	fmt.Println("A,L,3,4.2")
	fmt.Println("B,M,12,0")
	fmt.Println()
	fmt.Println(">")
}

func PrintMakersPortfolio(MM *[]*Portfolio, stop <-chan bool) {
	for {
		select {
		case <-stop:
			fmt.Print("\033[H\033[2J")
			return
		default:
			time.Sleep(time.Millisecond * 100)
			fmt.Print("\033[H\033[2J")
			for _, p := range *MM {
				p.Print()
			}
		}
	}
}
