package main

import "fmt"

type Printable interface {
	Print()
	PPrint()
}

func (order_book *OrderBook) Print() {
	fmt.Println()
	fmt.Println("OrderBook")
	fmt.Println()

	for i := range order_book.queue_ask.Len() {
		order_book.queue_ask.v[i].Print()
	}
	fmt.Println()
	for i := order_book.queue_bid.Len() - 1; i >= 0; i-- {
		order_book.queue_bid.v[i].Print()
	}
}

func (order_book *OrderBook) PPrint() {
	// TODO max value dynamically
	var quantity i32
	var depth string
	var ask_print []Pair

	for i := order_book.queue_ask.Len() - 1; i >= 0; i-- {
		quantity += order_book.queue_ask.v[i].size
		for range quantity/1000 + 1 {
			depth += "█"
		}
		ask_print = append(ask_print, Pair{order_book.queue_ask.v[i].price, depth})
	}
	for i := len(ask_print) - 1; i >= 0; i-- {
		fmt.Printf("$%.2f %s\n", ask_print[i].a, ask_print[i].b)
	}

	fmt.Printf("\nMidprice:$%.2f, Spread:%.2f%%\n\n", order_book.Midprice(), order_book.Spread()*100)
	depth = ""
	quantity = 0
	for i := order_book.queue_bid.Len() - 1; i >= 0; i-- {
		quantity += order_book.queue_bid.v[i].size
		for range quantity/1000 + 1 {
			depth += "█"
		}
		fmt.Printf("$%.2f %s\n", order_book.queue_bid.v[i].price, depth)
	}
}

func (order Order) Print() {
	var side string
	if order.side == BID {
		side = "BID"
	} else {
		side = "ASK"
	}
	fmt.Printf("%s : $%.2f (q:%d)\n", side, order.price, order.size)
}

func (order Order) PPrint() {
	order.Print()
}

func (fill FillReport) Print() {
	fmt.Printf("Filled %.2f%% at price $%.2f (q:%d)\n", fill.filled_pct, fill.price, fill.size)
}

func (fill FillReport) PPrint() {
	fill.Print()
}
