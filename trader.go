package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func createRandomOrder(midprice f32) *Order {
	var oside OrderSide
	if rand.Float32() > 0.5 {
		oside = BID
	} else {
		oside = ASK
	}

	return &Order{
		id:    rand.Uint64(),
		otype: LIMIT,
		side:  oside,
		size:  rand.Int31n(10) + 1,
		price: midprice + RandPrice()*2 - 1,
	}
}

func bootOB() *OrderBook {
	orderbook := OrderBook{}
	p := Portfolio{asset: 999_999, cash: 999_999}

	for range 10 {
		var oside OrderSide
		var price f32
		if rand.Float32() > 0.5 {
			oside = BID
			price = RandPrice()*10 + 5
		} else {
			oside = ASK
			price = RandPrice()*10 + 10
		}

		orderbook.Add(&Order{
			id:        rand.Uint64(),
			portfolio: &p,
			otype:     LIMIT,
			side:      oside,
			size:      rand.Int31n(10) + 1,
			price:     price, // revolve around 10
		})
	}

	return &orderbook
}

// Market Maker adds a limit order near the midprice.
// Executes every ~2 seconds.
func addTradeMM(orderbook *OrderBook, portfolio *Portfolio, obLock *sync.Mutex) {
	for {
		time.Sleep(time.Second * time.Duration(rand.Intn(2)))
		// if orderbook.Spread() > 0.001 {
		order := createRandomOrder(orderbook.Midprice())
		order.portfolio = portfolio
		obLock.Lock()
		orderbook.Add(order)
		obLock.Unlock()
		// }
	}
}

func userOrder(portfolio *Portfolio) *Order {
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
	case "A":
		oside = ASK
	case "B":
		oside = BID
	default:
		return nil
	}

	switch ainput[1] {
	case "L":
		otype = LIMIT
	case "M":
		otype = MARKET
	default:
		return nil
	}

	size, err := strconv.Atoi(ainput[2])
	if err != nil {
		return nil
	}

	if otype != MARKET {
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

func main() {
	orderbook := OrderBook{}
	orderbook = *bootOB()
	obLock := &sync.Mutex{}

	MM := []*Portfolio{}
	for range 5 {
		portfolio := Portfolio{asset: 999_999, cash: 999_999}
		go addTradeMM(&orderbook, &portfolio, obLock)
		MM = append(MM, &portfolio)
	}
	stop := make(chan bool)

	user_portfolio := Portfolio{asset: 0, cash: 100_000}

	for {
		fmt.Println("You're trading, type 'help' for help")
		fmt.Println(">")
		var input string
		fmt.Scan(&input)
		switch input {
		case "orderbook", "ob", "o":
			orderbook.PPrint()
		case "mid", "m":
			fmt.Println(orderbook.Midprice())
		case "display", "d":
			go PrintDisplay(&orderbook, stop)

		case "new", "n":
			fmt.Println("We'll create a new order.")
			order := userOrder(&user_portfolio)
			if order != nil {
				orderbook.Add(order)
				fmt.Println("Order submitted succesfully!")
			}
		case "portfolio", "p":
			user_portfolio.PPrint()
		case "makers", "mm":
			go PrintMakersPortfolio(&MM, stop)

		case "clear", "c":
			stop <- true
			fmt.Print("\033[H\033[2J")
		case "reset", "r":
			orderbook = *bootOB()
		case "quit", "q":
			os.Exit(0)
		default:
			PrintHelp()
		}
		fmt.Println()
	}
}
