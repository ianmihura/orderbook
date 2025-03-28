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

	for i := len(order_book.queue_ask) - 1; i >= 0; i-- {
		order_book.queue_ask[i].Print()
	}
	fmt.Println()
	for i := range len(order_book.queue_bid) {
		order_book.queue_bid[i].Print()
	}
}

func (order_book *OrderBook) PPrint() {
	// TODO max value dynamically
	var quantity i32
	var depth string
	var ask_print []Pair

	fmt.Println()
	fmt.Println("ASK")
	for i := range len(order_book.queue_ask) {
		quantity += order_book.queue_ask[i].size
		for range quantity/1000 + 1 {
			depth += "█"
		}
		ask_print = append(ask_print, Pair{order_book.queue_ask[i].price, depth})
	}
	for i := len(ask_print) - 1; i >= 0; i-- {
		fmt.Printf("$%f %s\n", ask_print[i].a, ask_print[i].b)
	}

	fmt.Println()
	depth = ""
	quantity = 0
	for i := range len(order_book.queue_bid) {
		quantity += order_book.queue_bid[i].size
		for range quantity/1000 + 1 {
			depth += "█"
		}
		fmt.Printf("$%f %s\n", order_book.queue_bid[i].price, depth)
	}
	fmt.Println("BID")
	fmt.Println()
}

func (order Order) Print() {
	var side string
	if order.side == BID {
		side = "BID"
	} else {
		side = "ASK"
	}
	fmt.Printf("%s : $%f (q:%d)\n", side, order.price, order.size)
}

func (order Order) PPrint() {
	order.Print()
}
