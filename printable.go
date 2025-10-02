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
	fmt.Printf("P: [%d] & [$%.2f]\n", portfolio.asset, portfolio.cash)
}
func (portfolio *Portfolio) PPrint() {
	fmt.Println("Your portfolio balance:")
	portfolio.Print()
	fmt.Println("[assets] & [cash]")
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
