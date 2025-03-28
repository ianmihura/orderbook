package main

import (
	"fmt"
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
	price     f32
	size      i32
	filledPct f32
}

// TODO do we need a struct for Queue?

func (orderBook *OrderBook) Add(order *Order) FillReport {
	if order.otype == MARKET {
		return orderBook.addMarket(order)
	} else if order.otype == LIMIT {
		return orderBook.addLimit(order)
	} else {
		panic("Unknown OrderType")
	}
}

func (orderBook *OrderBook) addMarket(order *Order) FillReport {
	initialOrderSize := order.size
	pendingSize := order.size
	var totalSpent f32
	var filled i32

	fillingOrder := &(*orderBook.GetQueueFlip(order.side))[0]
	for pendingSize > 0 && len(*orderBook.GetQueueFlip(order.side)) > 0 && limitPrice(order, fillingOrder.price) {
		// fmt.Println(len(queue))
		// fmt.Println(fillingOrder)
		if pendingSize >= fillingOrder.size {
			filled, _ = fillingOrder.Fill(fillingOrder.size)
			totalSpent += fillingOrder.price * f32(filled)
			// fmt.Println(fillingOrder)
			orderBook.remove(fillingOrder.side, 0)
		} else {
			filled, _ = fillingOrder.Fill(pendingSize)
			totalSpent += fillingOrder.price * f32(filled)
			// fmt.Println(fillingOrder)
		}
		pendingSize = pendingSize - filled
	}
	filledSize := initialOrderSize - pendingSize

	return FillReport{
		price:     f32(totalSpent) / f32(filledSize),
		size:      filledSize,
		filledPct: f32(filledSize) / f32(order.size),
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

func (orderBook *OrderBook) addLimit(order *Order) FillReport {
	queue := *orderBook.GetQueue(order.side)
	fillReport := FillReport{}

	if shouldFillLimitOrder(order, orderBook) {
		fillReport = orderBook.addMarket(order)
		order.size -= fillReport.size
		order.filledPct = fillReport.filledPct
	}

	if order.side == ASK {
		// TODO sort by datetime as well
		index := sort.Search(len(queue), func(i int) bool {
			return queue[i].price > order.price
		})
		orderBook.queue_ask = slices.Insert(queue, index, *order)
	} else if order.side == BID {
		index := sort.Search(len(queue), func(i int) bool {
			return queue[i].price < order.price
		})
		orderBook.queue_bid = slices.Insert(queue, index, *order)
	} else {
		panic("Unknown OrderSide")
	}

	return fillReport
}

func shouldFillLimitOrder(order *Order, orderBook *OrderBook) bool {
	queue := *orderBook.GetQueueFlip(order.side)
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

func (orderBook *OrderBook) Remove(order *Order) (*Order, error) {
	queue := *orderBook.GetQueue(order.side)
	ordersIndex := []int{}
	// TODO can make faster with bisect search (queue is sorted by price)
	for i := range queue {
		if queue[i].id == order.id {
			ordersIndex = append(ordersIndex, i)
		}
	}

	if len(ordersIndex) == 1 {
		o := orderBook.remove(order.side, ordersIndex[0])
		return &o, nil
	} else {
		return nil, &BaseError{"Found orders matching your index:", len(ordersIndex)}
	}
}

func (orderBook *OrderBook) remove(side OrderSide, i int) Order {
	var order Order
	if side == ASK {
		order = orderBook.queue_ask[i]
		orderBook.queue_ask = slices.Delete(orderBook.queue_ask, i, i+1)
	} else if side == BID {
		order = orderBook.queue_bid[i]
		orderBook.queue_bid = slices.Delete(orderBook.queue_bid, i, i+1)
	} else {
		panic("Unknown OrderSide")
	}

	return order
}

func (orderBook *OrderBook) GetQueue(side OrderSide) *[]Order {
	if side == BID {
		return &orderBook.queue_bid
	} else if side == ASK {
		return &orderBook.queue_ask
	} else {
		panic("Unknown OrderSide")
	}
}

func (orderBook *OrderBook) GetQueueFlip(side OrderSide) *[]Order {
	if side == BID {
		return &orderBook.queue_ask
	} else if side == ASK {
		return &orderBook.queue_bid
	} else {
		panic("Unknown OrderSide")
	}
}

func (orderBook *OrderBook) Midprice() f32 {
	// TODO test
	if len(orderBook.queue_ask) > 0 && len(orderBook.queue_bid) > 0 {
		return (orderBook.queue_ask[0].price + orderBook.queue_bid[0].price) / 2
	} else {
		return 0
	}
}

func (orderBook *OrderBook) Print() {
	fmt.Println()
	fmt.Println("OrderBook")
	fmt.Println()

	for i := len(orderBook.queue_ask) - 1; i >= 0; i-- {
		orderBook.queue_ask[i].Print()
	}
	fmt.Println()
	for i := range len(orderBook.queue_bid) {
		orderBook.queue_bid[i].Print()
	}
}

func (orderBook *OrderBook) PPrint() {
	// TODO max value dynamically
	var quantity i32
	var depth string
	var askPrint []Pair

	fmt.Println()
	fmt.Println("ASK")
	for i := range len(orderBook.queue_ask) {
		quantity += orderBook.queue_ask[i].size
		for range quantity/1000 + 1 {
			depth += "█"
		}
		askPrint = append(askPrint, Pair{orderBook.queue_ask[i].price, depth})
	}
	for i := len(askPrint) - 1; i >= 0; i-- {
		fmt.Printf("$%f %s\n", askPrint[i].a, askPrint[i].b)
	}

	fmt.Println()
	depth = ""
	quantity = 0
	for i := range len(orderBook.queue_bid) {
		quantity += orderBook.queue_bid[i].size
		for range quantity/1000 + 1 {
			depth += "█"
		}
		fmt.Printf("$%f %s\n", orderBook.queue_bid[i].price, depth)
	}
	fmt.Println("BID")
	fmt.Println()
}
