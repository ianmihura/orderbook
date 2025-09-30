package main

import (
	"slices"
	"sort"
)

/*
...
2: 53
1: 52
0: 51 ASK (queue_ask)

0: 49 BID (queue_bid)
1: 48
2: 47
...

TODO consider switching queue indexing for better perf
TODO (test that the hypotesis above)
*/

type OrderBook struct {
	queue_ask           []Order
	queue_bid           []Order
	transaction_history *TransactionHistory
}

type FillReport struct {
	price      f32
	size       i32
	filled_pct f32
}

// TODO do we need a struct for Queue?

// Adds Order to the OrderBook, with autofilling.
// TODO: Will execute Fill function in the Order _
//
// Returns a FillReport. Will panic if unknown OrderType
func (order_book *OrderBook) Add(order *Order) FillReport {
	switch order.otype {
	case MARKET:
		return addMarket(order_book, order)
	case LIMIT:
		return addLimit(order_book, order)
	default:
		panic("Unknown OrderType")
	}
}

// Fills an Order with existing market orders.
// If it's a LIMIT order, will fill until the limit price.
// If the Order size is bigger than the queue or its limit price is reached,
// we'll exit without completing the intended order size.
// Filled information is available in the returning FillReport.
//
// Side effects: Edits the active_order size and filled_pct,
// removes filled orders from the order_book, edits partial fill orders.
//
// Returns a FillReport.
func addMarket(order_book *OrderBook, active_order *Order) FillReport {
	init_order_size := active_order.size // the order size will change as it gets filled
	pending_size := active_order.size
	var total_spent f32
	var fill_report_tmp FillReport

	passive_order := &(*order_book.GetQueueFlip(active_order.side))[0]
	for pending_size > 0 &&
		len(*order_book.GetQueueFlip(active_order.side)) > 0 &&
		limitPrice(active_order, passive_order.price) {

		fill_report_tmp = *passive_order.Fill(active_order)

		if fill_report_tmp.filled_pct == 1 {
			order_book.RawRemove(passive_order.side, 0)
		}

		total_spent += fill_report_tmp.price * f32(fill_report_tmp.size)
		pending_size -= fill_report_tmp.size
	}
	filled_size := init_order_size - pending_size

	return FillReport{
		price:      f32(total_spent) / f32(filled_size),
		size:       filled_size,
		filled_pct: f32(filled_size) / f32(init_order_size),
	}
}

// TODO document
func limitPrice(order *Order, price f32) bool {
	if order.otype == MARKET {
		return true
	} else if order.side == ASK {
		return order.price <= price
	} else if order.side == BID {
		return order.price >= price
	} else {
		panic("Unknown OrderSide")
	}
}

// Will check if we should fill the limit Order, and
// will fill the Order up to the specified limit price.
// Then will add the remaining order to the queue.
//
// Returns a FillReport. Will panic if unknown OrderSide.
func addLimit(order_book *OrderBook, order *Order) FillReport {
	queue := *order_book.GetQueue(order.side)
	fill_report := FillReport{}

	if shouldFillLimitOrder(order_book, order) {
		fill_report = addMarket(order_book, order)
	}

	switch order.side {
	case ASK:
		// TODO sort by datetime as well
		index := sort.Search(len(queue), func(i int) bool {
			return queue[i].price > order.price
		})
		order_book.queue_ask = slices.Insert(queue, index, *order)
	case BID:
		index := sort.Search(len(queue), func(i int) bool {
			return queue[i].price < order.price
		})
		order_book.queue_bid = slices.Insert(queue, index, *order)
	default:
		panic("Unknown OrderSide")
	}

	order.order_book = order_book
	return fill_report
}

// Will check if a limit Order needs to be filled.
func shouldFillLimitOrder(order_book *OrderBook, order *Order) bool {
	queue := *order_book.GetQueueFlip(order.side)
	if len(queue) == 0 {
		return false
	} else if order.side == ASK && order.price <= queue[0].price {
		return true
	} else if order.side == BID && order.price >= queue[0].price {
		return true
	} else {
		return false
	}
}

// Finds and removes an Order from an OrderBook.
// Does not fill orders, only removes element from array.
//
// Returns removed Order, or error if finds != 1 Orders with same ID.
func (order_book *OrderBook) Remove(order *Order) (*Order, error) {
	queue := *order_book.GetQueue(order.side)
	order_idxs := []int{}
	for i := range queue {
		// TODO faster bisect search (is sorted by price)
		if queue[i].id == order.id {
			order_idxs = append(order_idxs, i)
		}
	}

	if len(order_idxs) == 1 {
		o := order_book.RawRemove(order.side, order_idxs[0])
		return &o, nil
	} else {
		return nil, &BaseError{"Found orders matching your index:", len(order_idxs)}
	}
}

// Removes the Order at the request side and index (i).
// Does not fill orders, only removes element from array.
//
// Returns removed Order. Panics if unknown OrderSide.
func (order_book *OrderBook) RawRemove(side OrderSide, i int) Order {
	var order Order
	switch side {
	case ASK:
		order = order_book.queue_ask[i]
		order_book.queue_ask = slices.Delete(order_book.queue_ask, i, i+1)
	case BID:
		order = order_book.queue_bid[i]
		order_book.queue_bid = slices.Delete(order_book.queue_bid, i, i+1)
	default:
		panic("Unknown OrderSide")
	}

	return order
}

// Gets corresponding queue.
func (order_book *OrderBook) GetQueue(side OrderSide) *[]Order {
	switch side {
	case BID:
		return &order_book.queue_bid
	case ASK:
		return &order_book.queue_ask
	default:
		panic("Unknown OrderSide")
	}
}

// Gets opposite corresponding queue. Useful to fill Orders,
// as orders get filled by opposite side Orders.
func (order_book *OrderBook) GetQueueFlip(side OrderSide) *[]Order {
	switch side {
	case BID:
		return &order_book.queue_ask
	case ASK:
		return &order_book.queue_bid
	default:
		panic("Unknown OrderSide")
	}
}

// TODO document
func (order_book *OrderBook) Midprice() f32 {
	// TODO test
	if len(order_book.queue_ask) > 0 && len(order_book.queue_bid) > 0 {
		return (order_book.queue_ask[0].price + order_book.queue_bid[0].price) / 2
	} else {
		return 0
	}
}
