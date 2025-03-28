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
	queue_ask []Order
	queue_bid []Order
}

type FillReport struct {
	price      f32
	size       i32
	filled_pct f32
}

// TODO do we need a struct for Queue?

func (order_book *OrderBook) Add(order *Order) FillReport {
	if order.otype == MARKET {
		return order_book.addMarket(order)
	} else if order.otype == LIMIT {
		return order_book.addLimit(order)
	} else {
		panic("Unknown OrderType")
	}
}

func (order_book *OrderBook) addMarket(order *Order) FillReport {
	init_order_size := order.size
	pending_size := order.size
	var total_spent f32
	var filled i32

	filling_order := &(*order_book.GetQueueFlip(order.side))[0]
	for pending_size > 0 && len(*order_book.GetQueueFlip(order.side)) > 0 && limitPrice(order, filling_order.price) {
		// fmt.Println(len(queue))
		// fmt.Println(filling_order)
		if pending_size >= filling_order.size {
			filled, _ = filling_order.Fill(filling_order.size)
			total_spent += filling_order.price * f32(filled)
			// fmt.Println(filling_order)
			order_book.remove(filling_order.side, 0)
		} else {
			filled, _ = filling_order.Fill(pending_size)
			total_spent += filling_order.price * f32(filled)
			// fmt.Println(filling_order)
		}
		pending_size = pending_size - filled
	}
	filled_size := init_order_size - pending_size

	return FillReport{
		price:      f32(total_spent) / f32(filled_size),
		size:       filled_size,
		filled_pct: f32(filled_size) / f32(order.size),
	}
}

func limitPrice(o *Order, price f32) bool {
	if o.otype == MARKET {
		return true
	} else if o.side == ASK {
		return o.price <= price
	} else if o.side == BID {
		return o.price >= price
	} else {
		panic("Unknown OrderSide")
	}
}

func (order_book *OrderBook) addLimit(order *Order) FillReport {
	queue := *order_book.GetQueue(order.side)
	fill_report := FillReport{}

	if shouldFillLimitOrder(order, order_book) {
		fill_report = order_book.addMarket(order)
		order.size -= fill_report.size
		order.filled_pct = fill_report.filled_pct
	}

	if order.side == ASK {
		// TODO sort by datetime as well
		index := sort.Search(len(queue), func(i int) bool {
			return queue[i].price > order.price
		})
		order_book.queue_ask = slices.Insert(queue, index, *order)
	} else if order.side == BID {
		index := sort.Search(len(queue), func(i int) bool {
			return queue[i].price < order.price
		})
		order_book.queue_bid = slices.Insert(queue, index, *order)
	} else {
		panic("Unknown OrderSide")
	}

	return fill_report
}

func shouldFillLimitOrder(order *Order, order_book *OrderBook) bool {
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

func (order_book *OrderBook) Remove(order *Order) (*Order, error) {
	queue := *order_book.GetQueue(order.side)
	order_idxs := []int{}
	// TODO can make faster with bisect search (queue is sorted by price)
	for i := range queue {
		if queue[i].id == order.id {
			order_idxs = append(order_idxs, i)
		}
	}

	if len(order_idxs) == 1 {
		o := order_book.remove(order.side, order_idxs[0])
		return &o, nil
	} else {
		return nil, &BaseError{"Found orders matching your index:", len(order_idxs)}
	}
}

func (order_book *OrderBook) remove(side OrderSide, i int) Order {
	var order Order
	if side == ASK {
		order = order_book.queue_ask[i]
		order_book.queue_ask = slices.Delete(order_book.queue_ask, i, i+1)
	} else if side == BID {
		order = order_book.queue_bid[i]
		order_book.queue_bid = slices.Delete(order_book.queue_bid, i, i+1)
	} else {
		panic("Unknown OrderSide")
	}

	return order
}

func (order_book *OrderBook) GetQueue(side OrderSide) *[]Order {
	if side == BID {
		return &order_book.queue_bid
	} else if side == ASK {
		return &order_book.queue_ask
	} else {
		panic("Unknown OrderSide")
	}
}

func (order_book *OrderBook) GetQueueFlip(side OrderSide) *[]Order {
	if side == BID {
		return &order_book.queue_ask
	} else if side == ASK {
		return &order_book.queue_bid
	} else {
		panic("Unknown OrderSide")
	}
}

func (order_book *OrderBook) Midprice() f32 {
	// TODO test
	if len(order_book.queue_ask) > 0 && len(order_book.queue_bid) > 0 {
		return (order_book.queue_ask[0].price + order_book.queue_bid[0].price) / 2
	} else {
		return 0
	}
}
