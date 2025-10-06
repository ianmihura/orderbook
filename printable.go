package main

import (
	"fmt"
	"strings"
)

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
	max := orderbook.queue_ask.Bottom().price
	min := orderbook.queue_bid.Bottom().price

	point_count := 30
	point_size := (max - min) / f32(point_count)

	pprintQueueAsk(orderbook.queue_ask, point_size, point_count/2)
	fmt.Printf("\nMidprice:$%.2f, Spread:%.2f%%\n\n", orderbook.Midprice(), orderbook.Spread()*100)
	pprintQueueBid(orderbook.queue_bid, point_size, point_count/2)
	fmt.Println("            10   20   30   40   50   60   70   80   90   100")
	//                  ^2468^2468^2468^2468^2468^2468^2468^2468^2468^2468^
}
func pprintQueueAsk(queue Queue, point_size f32, h_point_count int) {
	min := queue.Top().price
	_print := make([]i32, h_point_count)
	o := len(queue.v) - 1
out:
	for p := range h_point_count {
		if p != 0 {
			_print[p] = _print[p-1]
		}
		for queue.v[o].price < min+f32(p+1)*point_size {
			_print[p] += queue.v[o].size
			o--
			if o < 0 {
				break out
			}
		}
	}
	for i := len(_print) - 1; i >= 0; i-- {
		depth := strings.Repeat("█", int(_print[i]/2))
		fmt.Printf("$%5.2f %s\n", min+f32(i)*point_size, depth)
	}
}
func pprintQueueBid(queue Queue, point_size f32, h_point_count int) {
	max := queue.Top().price
	_print := make([]i32, h_point_count)
	o := len(queue.v) - 1
out:
	for p := range h_point_count {
		if p != 0 {
			_print[p] = _print[p-1]
		}
		for queue.v[o].price > max-f32(p+1)*point_size {
			_print[p] += queue.v[o].size
			o--
			if o < 0 {
				break out
			}
		}
		depth := strings.Repeat("█", int(_print[p]/2))
		fmt.Printf("$%5.2f %s\n", max-f32(p)*point_size, depth)
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

func (portfolio *Portfolio) Print(midprice f32) {
	fmt.Printf("[$%.2f] & [x%d] = $%.2f\n", portfolio.cash, portfolio.asset, f32(portfolio.asset)*midprice+portfolio.cash)
}
func (portfolio *Portfolio) PPrint(midprice f32) {
	fmt.Println("Your portfolio is worth:")
	fmt.Printf("$%.2f\n", f32(portfolio.asset)*midprice+portfolio.cash)
	fmt.Println("Portfolio balance:")
	fmt.Printf("[$%.2f] & [x%d]\n", portfolio.cash, portfolio.asset)
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
