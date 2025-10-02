package main

import "fmt"

type Printable interface {
	Print()
	PPrint()
}

func (orderbook *OrderBook) Print() {
	fmt.Println()
	fmt.Println("OrderBook")
	fmt.Println()

	for i := range orderbook.queue_ask.Len() {
		orderbook.queue_ask.v[i].Print()
	}
	fmt.Println()
	for i := orderbook.queue_bid.Len() - 1; i >= 0; i-- {
		orderbook.queue_bid.v[i].Print()
	}
}
func (orderbook *OrderBook) PPrint() {
	// TODO max value dynamically
	var quantity i32
	var depth string
	var ask_print []Tuple

	for i := orderbook.queue_ask.Len() - 1; i >= 0; i-- {
		quantity += orderbook.queue_ask.v[i].size
		for range quantity/1000 + 1 {
			depth += "█"
		}
		ask_print = append(ask_print, Tuple{orderbook.queue_ask.v[i].price, depth, orderbook.queue_ask.v[i].portfolio.is_user})
	}
	for i := len(ask_print) - 1; i >= 0; i-- {
		fmt.Printf("$%.2f %s", ask_print[i].a, ask_print[i].b)
		if ask_print[i].c == true {
			fmt.Print(" (U)")
		}
		fmt.Print("\n")
	}

	fmt.Printf("\nMidprice:$%.2f, Spread:%.2f%%\n\n", orderbook.Midprice(), orderbook.Spread()*100)
	depth = ""
	quantity = 0
	for i := orderbook.queue_bid.Len() - 1; i >= 0; i-- {
		quantity += orderbook.queue_bid.v[i].size
		for range quantity/1000 + 1 {
			depth += "█"
		}
		fmt.Printf("$%.2f %s", orderbook.queue_bid.v[i].price, depth)
		if orderbook.queue_bid.v[i].portfolio.is_user {
			fmt.Print(" (U)")
		}
		fmt.Print("\n")
	}
}

func (order *Order) Print() {
	var side string
	if order.side == BID {
		side = "BID"
	} else {
		side = "ASK"
	}
	fmt.Printf("%s (%s) : $%.2f x%d\n", side, string(order.otype), order.price, order.size)
}
func (order *Order) PPrint() {
	order.Print()
}

func (fill *FillReport) Print() {
	fmt.Printf("Filled %.0f%% at $%.2f x%d", fill.filled_pct*100, fill.price, fill.size)
	if fill.is_active {
		fmt.Print(" (A)\n")
	} else {
		fmt.Print("\n")
	}
}
func (fill *FillReport) PPrint() {
	fill.Print()
}

func (portfolio *Portfolio) Print() {
	fmt.Printf("[$%.2f] & [x%d]\n", portfolio.cash, portfolio.asset)
}
func (portfolio *Portfolio) PPrint(midprice f32) {
	fmt.Println("Your portfolio is worth:")
	fmt.Printf("$%.2f\n", f32(portfolio.asset)*midprice+portfolio.cash)
	fmt.Println("Portfolio balance:")
	portfolio.Print()
	fmt.Println("[cash] & [assets]")
	fmt.Println()
}

func (tx *Transaction) Print() {
	fmt.Println("Past transaction:")
	fmt.Println("Active order:")
	tx.active.Print()
	fmt.Println("Last report:")
	tx.reports[len(tx.reports)-1].Print()
	fmt.Println()
}
func (tx *Transaction) PPrint() {
	fmt.Println("Past transaction:")
	fmt.Println("  Active order:")
	tx.active.Print()
	fmt.Println("  Passive orders:")
	for _, order := range tx.passive {
		order.Print()
	}
	fmt.Println("  Reports:")
	for _, report := range tx.reports {
		report.Print()
	}
	fmt.Println()
}

func (txh *TransactionHistory) Print() {
	for _, tx := range txh.txs {
		tx.Print()
	}
}
func (txh *TransactionHistory) PPrint() {
	for _, tx := range txh.txs {
		tx.PPrint()
	}
	fmt.Println("Total transactions:", len(txh.txs))
}
