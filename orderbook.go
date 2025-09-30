package main

import (
	"sort"
)

type OrderBook struct {
	queue_ask           Queue
	queue_bid           Queue
	transaction_history *TransactionHistory
}

type FillReport struct {
	price      f32
	size       i32
	filled_pct f32
}

// Adds Order to the OrderBook, with autofilling.
// Returns a FillReport. Will panic if unknown OrderType
func (order_book *OrderBook) Add(order *Order) FillReport {
	// TODO: Will execute Fill function in the Order _
	switch order.otype {
	case MARKET:
		return addToMarket(order_book, order)
	case LIMIT:
		return addLimit(order_book, order)
	default:
		// TODO more order types coming
		panic("Unknown OrderType")
	}
}

// Fills an Order with existing market orders.
// If it's a LIMIT order, will fill until the limit price.
// If the Order size is bigger than the queue, or its limit price is reached,
// we'll exit without completing the intended order size.
//
// Side effects: Edits the active_order size and filled_pct,
// removes filled orders from the order_book, edits partial fill orders.
//
// Returns a FillReport. Filled information is available here.
func addToMarket(order_book *OrderBook, active_order *Order) FillReport {
	init_order_size := active_order.size // the order size will change as it gets filled
	pending_size := active_order.size
	var total_spent f32
	var fill_report_tmp FillReport
	queue_flip := order_book.GetQueueFlip(active_order.side)

	for pending_size > 0 && !queue_flip.IsEmpty() &&
		shouldFillOrder(active_order, queue_flip.Top().price) {

		fill_report_tmp = *queue_flip.Top().Fill(active_order)

		if fill_report_tmp.filled_pct == 1 {
			queue_flip.Pop()
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

// Will check if we should fill the limit Order, and
// will fill the Order up to the specified limit price.
// Then, will add the remaining order to the queue.
//
// Returns a FillReport.
func addLimit(order_book *OrderBook, order *Order) FillReport {
	queue := *order_book.GetQueue(order.side)
	fill_report := FillReport{}

	if shouldFillLimitOrder(order_book, order) {
		fill_report = addToMarket(order_book, order)
		if fill_report.filled_pct == 1.0 {
			return fill_report
		}
	}

	if order.side == BID {
		index := sort.Search(queue.Len(), func(i int) bool {
			return queue.v[i].price > order.price
		})
		order_book.queue_bid.Insert(index, *order)
	} else {
		index := sort.Search(queue.Len(), func(i int) bool {
			return queue.v[i].price < order.price
		})
		order_book.queue_ask.Insert(index, *order)
	}

	order.order_book = order_book
	return fill_report
}

// Will check if a limit Order needs to be filled.
func shouldFillLimitOrder(order_book *OrderBook, order *Order) bool {
	queue := *order_book.GetQueueFlip(order.side)
	if queue.IsEmpty() {
		return false
	} else if order.side == ASK && order.price <= queue.Top().price {
		return true
	} else if order.side == BID && order.price >= queue.Top().price {
		return true
	} else {
		return false
	}
}

// Will check if an order still needs to be filled, by checking otype or limit price
func shouldFillOrder(order *Order, price f32) bool {
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

// Finds and removes an Order from an OrderBook.
// Does not fill orders, only removes element from stack.
//
// Returns removed Order, or error if finds != 1 Orders with same ID.
func (order_book *OrderBook) Remove(order *Order) (*Order, error) {
	queue := *order_book.GetQueue(order.side)
	order_idxs := queue.FindAll(*order)

	if len(order_idxs) == 1 {
		o := order_book.rawRemove(order.side, order_idxs[0])
		return &o, nil
	} else {
		return nil, &BaseError{"Found orders matching your index:", len(order_idxs)}
	}
}

// Removes the Order at the request side and index (i).
// Does not fill orders, only removes element from array.
//
// Returns removed Order.
func (order_book *OrderBook) rawRemove(side OrderSide, i int) Order {
	var order Order
	if side == BID {
		order = order_book.queue_bid.Remove(i)
	} else {
		order = order_book.queue_ask.Remove(i)
	}

	return order
}

// Gets corresponding queue.
func (order_book *OrderBook) GetQueue(side OrderSide) *Queue {
	if side == BID {
		return &order_book.queue_bid
	} else {
		return &order_book.queue_ask
	}
}

// Gets opposite corresponding queue. Useful to fill Orders,
// as orders get filled by opposite side Orders.
func (order_book *OrderBook) GetQueueFlip(side OrderSide) *Queue {
	if side == BID {
		return &order_book.queue_ask
	} else {
		return &order_book.queue_bid
	}
}

// Returns the midpoint (avg) between the two Top prices, known as the MidPrice
func (order_book *OrderBook) Midprice() f32 {
	if order_book.queue_ask.IsEmpty() || order_book.queue_bid.IsEmpty() {
		return 0.0
	} else {
		return (order_book.queue_ask.Top().price + order_book.queue_bid.Top().price) / 2
	}
}
