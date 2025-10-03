package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
			fmt.Print("\nType c to close")
		}
	}
}

func PrintHelp() {
	fmt.Println("")
	fmt.Println("List of commands:")
	fmt.Println("  Orderbook:")
	fmt.Println("    d, display: watch the ob go by")
	fmt.Println("    a, auto: display portfolio of all auto traders")
	fmt.Println("    t, txs: see all past transactions")
	fmt.Println("  User:")
	fmt.Println("    n, new: create new order")
	fmt.Println("    p, portfolio: check your portfolio")
	fmt.Println("  App:")
	fmt.Println("    c, clear, close: clear the console & close any ongoing display")
	fmt.Println("    r, reset: gives new initial values to orderbook")
	fmt.Println("    q, quit: quits the program")
}

func PrintNewOrderHelp() {
	fmt.Println("Order format: B|S,TYPE,size,price")
	fmt.Println("Where: B|S is the side Buy or Sell")
	fmt.Println("Where: L|M|D|V|T is the type")
	fmt.Println("  Limit, Market, Midprice, VWAP, TWAP")
	fmt.Println("  Note that Market, Midprice, VWAP, TWAP orders will ignore price")
	fmt.Println("Where: Size in an int")
	fmt.Println("Where: Price is a float")
	fmt.Println("* Values must be separated by a single comma (, char)")
	fmt.Println()
	fmt.Println("Some examples...")
	fmt.Println("S,L,3,4.2")
	fmt.Println("B,M,12")
	fmt.Println()
	fmt.Println(">")
}

func PrintAutoTradersPortfolio(auto_traders *[]*Portfolio, orderbook *OrderBook, stop <-chan bool) {
	for {
		select {
		case <-stop:
			fmt.Print("\033[H\033[2J")
			return
		default:
			time.Sleep(time.Millisecond * 100)
			fmt.Print("\033[H\033[2J")
			for _, p := range *auto_traders {
				p.Print(orderbook.Midprice())
			}
			fmt.Print("\nType c to close")
		}
	}
}

func UserOrder(portfolio *Portfolio) *Order {
	PrintNewOrderHelp()

	var input string
	fmt.Scan(&input)
	ainput := strings.Split(input, ",")
	if len(ainput) < 3 {
		return nil
	}

	var otype OrderType
	var oside OrderSide
	var price f64
	switch ainput[0] {
	case "S", "s":
		oside = ASK
	case "B", "b":
		oside = BID
	default:
		return nil
	}

	switch ainput[1] {
	case "L", "l":
		otype = LIMIT
	case "M", "m":
		otype = MARKET
	case "D", "d":
		otype = MID
	case "V", "v":
		otype = VWAP
	case "T", "t":
		otype = TWAP
	default:
		return nil
	}

	size, err := strconv.Atoi(ainput[2])
	if err != nil {
		return nil
	}

	if otype == LIMIT {
		if len(ainput) < 4 {
			return nil
		}
		price, err = strconv.ParseFloat(ainput[3], 32)
		if err != nil {
			return nil
		}
	}

	return &Order{
		id:        rand.Uint64(),
		portfolio: portfolio,
		otype:     otype,
		side:      oside,
		size:      i32(size),
		price:     f32(price),
	}
}
