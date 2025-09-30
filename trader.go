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
			id:    rand.Uint64(),
			otype: LIMIT,
			side:  oside,
			size:  rand.Int31n(10) + 1,
			price: price, // revolve around 10
		})
	}

	return &orderbook
}

// Market Maker adds a limit order near the midprice.
// Executes every ~2 seconds.
func addTradeMM(orderbook *OrderBook, obLock *sync.Mutex) {
	for {
		time.Sleep(time.Second * time.Duration(rand.Intn(2)))
		// if orderbook.Spread() > 0.001 {
		order := createRandomOrder(orderbook.Midprice())
		obLock.Lock()
		orderbook.Add(order)
		obLock.Unlock()
		// }
	}
}

func userOrder() *Order {
	PrintNewOrderHelp()

	var input string
	fmt.Scan(&input)
	ainput := strings.Split(input, ",")

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
		price, err = strconv.ParseFloat(ainput[3], 32)
		if err != nil {
			return nil
		}
	}

	return &Order{
		id:    rand.Uint64(),
		otype: otype,
		side:  oside,
		size:  i32(size),
		price: f32(price),
	}
}

func main() {
	orderbook := OrderBook{}
	orderbook = *bootOB()
	obLock := &sync.Mutex{}

	MM := 5
	for range MM {
		go addTradeMM(&orderbook, obLock)
	}
	stop := make(chan bool)

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
		case "new", "n":
			fmt.Println("We'll create a new order.")
			order := userOrder()
			if order != nil {
				orderbook.Add(order)
				fmt.Println("Order submitted succesfully!")
			}
		case "clear", "c":
			stop <- true
			fmt.Print("\033[H\033[2J")
		case "display", "d":
			go PrintDisplay(&orderbook, stop)
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
