package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	orderbook := OrderBook{}
	lock := &sync.Mutex{}
	orderbook = *BootOrderbook(40, lock)

	auto_traders := AddAutoTraders(&orderbook)

	user_portfolio := Portfolio{asset: 0, cash: 100_000, is_user: true}

	stop := make(chan bool)
	for {
		fmt.Println("You're trading, type 'h' for help")
		fmt.Println(">")
		var input string
		fmt.Scan(&input)
		switch input {
		case "display", "d":
			go PrintDisplay(&orderbook, stop)
		case "auto", "a":
			go PrintAutoTradersPortfolio(auto_traders, &orderbook, stop)
		case "txs", "t":
			getTxHistory().PPrint()
			fmt.Println()
			fmt.Println("Midprice:", orderbook.Midprice())
			fmt.Println("TWAP:    ", GetAvgPrice())
			fmt.Println("VWAP:    ", GetAvgPriceWeighted())

		case "new", "n":
			fmt.Println("We'll create a new order.")
			order := UserOrder(&user_portfolio)
			if order != nil {
				orderbook.lock.Lock()
				report := orderbook.Add(order)
				orderbook.lock.Unlock()
				fmt.Println("Order submitted succesfully!")
				report.PPrint()
			}
		case "portfolio", "p":
			user_portfolio.PPrint(orderbook.Midprice())

		case "clear", "close", "c":
			stop <- true
			fmt.Print("\033[H\033[2J")
		case "reset", "r":
			_ob := BootOrderbook(40, orderbook.lock)
			_ob.lock.Lock()
			orderbook = *_ob
			_ob.lock.Unlock()
			// no need to unlock because the lock gets defined again
			// TODO reset history of tx
		case "quit", "q":
			os.Exit(0)
		default:
			PrintHelp()
		}
		fmt.Println()
	}
}
