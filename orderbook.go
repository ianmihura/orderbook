package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type OrderBook struct {
	queue_ask Queue
	queue_bid Queue
	lock      *sync.Mutex
}

type FillReport struct {
	price      f32
	size       i32
	filled_pct f32
	is_active  bool
}

// Adds Order to the OrderBook, with autofilling.
// Returns a FillReport. Will panic if unknown OrderType
func (orderbook *OrderBook) Add(order *Order) *FillReport {
	// TODO: Will execute Fill function in the Order _
	// TODO refactor this to a strategy pattern
	switch order.otype {
	case MARKET:
		return addToMarket(orderbook, order)
	case LIMIT:
		return addLimit(orderbook, order)
	case MID:
		return addMidprice(orderbook, order)
	case TWAP:
		return addTWAP(orderbook, order)
	case VWAP:
		return addVWAP(orderbook, order)
	default:
		panic("Unknown OrderType")
	}
}

// Fills an Order with existing market orders.
// If it's a LIMIT order, will fill until the limit price.
// If the Order size is bigger than the queue, or its limit price is reached,
// we'll exit without completing the intended order size.
//
// Side effects: Edits the active_order size and filled_pct,
// removes filled orders from the orderbook, edits partial fill orders.
//
// Returns a FillReport. Filled information is available here.
func addToMarket(orderbook *OrderBook, active_order *Order) *FillReport {
	init_order_size := active_order.size // the order size will change as it gets filled
	pending_size := active_order.size
	var total_spent f32
	var fill_report_tmp FillReport
	queue_flip := orderbook.GetQueueFlip(active_order.side)

	for pending_size > 0 && !queue_flip.IsEmpty() &&
		shouldFillOrder(active_order, queue_flip.Top().price) {

		// fill_report_tmp corresponds to passive_order
		fill_report_tmp = *queue_flip.Top().Fill(active_order)

		if fill_report_tmp.filled_pct == 1 {
			// Removing empty top order
			queue_flip.Pop()
		} else if fill_report_tmp.filled_pct == 0 {
			// Empty report => could not fill
			break
		}

		total_spent += fill_report_tmp.price * f32(fill_report_tmp.size)
		pending_size -= fill_report_tmp.size
	}
	filled_size := init_order_size - pending_size

	fill_report_tmp = FillReport{
		price:      f32(total_spent) / f32(filled_size),
		size:       filled_size,
		filled_pct: f32(filled_size) / f32(init_order_size),
		is_active:  true,
	}
	AddFillReport(active_order.id, &fill_report_tmp)
	return &fill_report_tmp
}

// Will check if we should fill the limit Order, and
// will fill the Order up to the specified limit price.
// Then, will add the remaining order to the queue.
//
// Returns a FillReport.
func addLimit(orderbook *OrderBook, order *Order) *FillReport {
	queue := *orderbook.GetQueue(order.side)
	fill_report := &FillReport{}

	if shouldFillLimitOrder(orderbook, order) {
		fill_report = addToMarket(orderbook, order)
		if fill_report.filled_pct == 1.0 {
			return fill_report
		}
	}

	if order.side == BID {
		index := sort.Search(queue.Len(), func(i int) bool {
			return queue.v[i].price > order.price
		})
		orderbook.queue_bid.Insert(index, *order)
	} else {
		index := sort.Search(queue.Len(), func(i int) bool {
			return queue.v[i].price < order.price
		})
		orderbook.queue_ask.Insert(index, *order)
	}

	order.orderbook = orderbook
	return fill_report
}

// Will check if a limit Order needs to be filled.
func shouldFillLimitOrder(orderbook *OrderBook, order *Order) bool {
	queue := *orderbook.GetQueueFlip(order.side)
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

func addMidprice(orderbook *OrderBook, order *Order) *FillReport {
	order.price = orderbook.Midprice()
	return addLimit(orderbook, order)
}

// Will split up orders in random slots.
// Size will be typical order size +- 20%.
// Spaced out, every few random seconds [0,3).
// Price will always be the Average executed price.
func addTWAP(orderbook *OrderBook, order *Order) *FillReport {
	// TODO finish this thing
	fill_report := FillReport{}
	for order.size >= 0 {
		avgSize := GetAvgSize()
		_order := order
		_order.id = rand.Uint64()
		_order.price = GetAvgPrice()
		_order.size = avgSize + i32(f32(avgSize)*(rand.Float32()-0.5)/2.5) // error term: 0.5/2.5 = ±20%
		_order.otype = LIMIT
		fmt.Println("Doing one")
		fill_report_tmp := addLimit(orderbook, _order)
		time.Sleep(time.Second * time.Duration(rand.Intn(3)))
		order.size -= fill_report_tmp.size

		// How do we get all the filling info?
		// v1. naif
		// v2. we wait for a fill every time before we submit the next
		// v3. goroutine every submit, when we get the fillreport we add the number
		// fill_report.size += fill_report_tmp.size
	}
	// fill_report = FillReport{
	// 	price:      f32(total_spent) / f32(filled_size),
	// 	size:       filled_size,
	// 	filled_pct: f32(filled_size) / f32(init_order_size),
	// 	is_active:  true,
	// }
	return &fill_report
}

// Will split up orders in random slots.
// Size will be typical order size +- 20%.
// Spaced out, every few random seconds [0,4].
// Price will always be the Average executed price. weighted by volume.
func addVWAP(orderbook *OrderBook, order *Order) *FillReport {
	order.price = GetAvgPriceWeighted()
	// TODO split-up
	return addLimit(orderbook, order)
}

// Finds and removes an Order from an OrderBook.
// Does not fill orders, only removes element from stack.
//
// Returns removed Order, or error if finds != 1 Orders with same ID.
func (orderbook *OrderBook) Remove(order *Order) (*Order, error) {
	queue := *orderbook.GetQueue(order.side)
	order_idxs := queue.FindAll(*order)

	if len(order_idxs) == 1 {
		o := orderbook.rawRemove(order.side, order_idxs[0])
		return &o, nil
	} else {
		return nil, &BaseError{"Found orders matching your index:", len(order_idxs)}
	}
}

// Removes the Order at the request side and index (i).
// Does not fill orders, only removes element from array.
//
// Returns removed Order.
func (orderbook *OrderBook) rawRemove(side OrderSide, i int) Order {
	var order Order
	if side == BID {
		order = orderbook.queue_bid.Remove(i)
	} else {
		order = orderbook.queue_ask.Remove(i)
	}

	return order
}

// Gets corresponding queue.
func (orderbook *OrderBook) GetQueue(side OrderSide) *Queue {
	if side == BID {
		return &orderbook.queue_bid
	} else {
		return &orderbook.queue_ask
	}
}

// Gets opposite corresponding queue. Useful to fill Orders,
// as orders get filled by opposite side Orders.
func (orderbook *OrderBook) GetQueueFlip(side OrderSide) *Queue {
	if side == BID {
		return &orderbook.queue_ask
	} else {
		return &orderbook.queue_bid
	}
}

// Returns the midpoint (avg) between the two Top prices, known as the MidPrice
func (orderbook *OrderBook) Midprice() f32 {
	if orderbook.queue_ask.IsEmpty() || orderbook.queue_bid.IsEmpty() {
		return 0.0
	} else {
		return (orderbook.queue_ask.Top().price + orderbook.queue_bid.Top().price) / 2
	}
}

// Bid-Ask spread (%)
func (orderbook *OrderBook) Spread() f32 {
	if orderbook.queue_ask.IsEmpty() || orderbook.queue_bid.IsEmpty() {
		return 0.0
	}
	ask_price := orderbook.queue_ask.Top().price
	bid_price := orderbook.queue_bid.Top().price
	return (ask_price - bid_price) / orderbook.Midprice()
}
