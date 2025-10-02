package main

import (
	"cmp"
	"math/rand"
	"slices"
	"testing"
)

func createQueue(length int) Queue {
	q := Queue{[]Order{}}
	for range length {
		q.Push(Order{})
	}
	return q
}

func TestGetQueues(t *testing.T) {
	oq_ask := createQueue(1)
	oq_bid := createQueue(2)
	ob := OrderBook{queue_ask: oq_ask, queue_bid: oq_bid}

	q_bid := *ob.GetQueue(BID)
	q_ask := *ob.GetQueue(ASK)

	Assert(t, q_bid.Len() == 2, "Got wrong queue ", q_bid)
	Assert(t, q_ask.Len() == 1, "Got wrong queue ", q_ask)

	q_ask_flip := *ob.GetQueueFlip(BID)
	q_bid_flip := *ob.GetQueueFlip(ASK)

	Assert(t, q_bid_flip.Len() == 2, "Got wrong queue ", q_bid_flip)
	Assert(t, q_ask_flip.Len() == 1, "Got wrong queue ", q_ask_flip)
}

func TestRemove(t *testing.T) {
	oq_bid := Queue{}
	oq_bid.CopyFromSlice([]Order{{id: 1}, {id: 2}, {id: 3}})
	ob := OrderBook{queue_bid: oq_bid}

	deleted_order, err := ob.Remove(&oq_bid.v[1])

	Assert(t, err == nil, err)
	Assert(t, deleted_order.id == 2, "Identified wrong order to delete")
	Assert(t, ob.queue_bid.Len() == 2, "Length of queue should be 2")
	Assert(t,
		ob.queue_bid.v[0].id == 1 && ob.queue_bid.v[1].id == 3,
		"Deleted wrong Order",
	)
}

func TestAddLimit(t *testing.T) {
	ob := OrderBook{}
	for range 100 {
		var oside OrderSide
		var price f32
		if rand.Float32() > 0.5 {
			oside = BID
			price = rand.Float32()
		} else {
			oside = ASK
			price = rand.Float32() + 1.0 // offset ensures we dont autofill limit orders
		}

		o := Order{
			id:    rand.Uint64(),
			otype: LIMIT,
			side:  oside,
			price: price,
		}
		ob.Add(&o)
	}

	Assert(t, (ob.queue_ask.Len()+ob.queue_bid.Len()) == 100, "Should insert 100 elements")
	Assert(t, !slices.IsSortedFunc(ob.queue_ask.v, func(a, b Order) int {
		return cmp.Compare(a.id, b.id)
	}), "ASK queue Should not be sorted by ID")
	Assert(t, !slices.IsSortedFunc(ob.queue_bid.v, func(a, b Order) int {
		return cmp.Compare(a.id, b.id)
	}), "BID queue Should not be sorted by ID")
	Assert(t, slices.IsSortedFunc(ob.queue_ask.v, func(a, b Order) int {
		return cmp.Compare(b.price, a.price)
	}), "ASK queue Should be sorted by price")
	Assert(t, IsSortedFuncDesc(ob.queue_bid.v, func(a, b Order) int {
		return cmp.Compare(b.price, a.price)
	}), "BID queue Should be sorted by price")

	// TODO not implemented
	// Assert(t, slices.IsSortedFunc(ob.queue_ask, func(a, b Order) int {
	// 	return a.created.Compare(b.created)
	// }), "ASK queue Should be sorted by created datetime")
	// Assert(t, slices.IsSortedFunc(ob.queue_bid, func(a, b Order) int {
	// 	return a.created.Compare(b.created)
	// }), "BID queue Should be sorted by created datetime")
}

func TestAddAutofillLimit(t *testing.T) {
	ob := OrderBook{}
	p := Portfolio{asset: 999_999, cash: 999_999}
	for i := range 50 {
		ob.Add(&Order{
			id:        rand.Uint64(),
			otype:     LIMIT,
			side:      ASK,
			size:      5,
			price:     f32(101 - i),
			portfolio: &p,
		})
	}

	// Partial fills itself with 1 order of size 5, at price (52)
	fill := ob.Add(&Order{
		id:        rand.Uint64(),
		otype:     LIMIT,
		side:      BID,
		size:      6,
		price:     52,
		portfolio: &p,
	})
	Assert(t, fill.size == 5, "Should fill all the order", fill.size)
	Assert(t, fill.price == 52, "Should fill the order at correct price", fill.price)
	Assert(t, fill.filled_pct == 5.0/6.0, "Should fill all filled_pct", fill.filled_pct)
	Assert(t, ob.queue_ask.Top().size == 5, "We should see the untouched order in the queue, with size", ob.queue_ask.Top().size)
	Assert(t, ob.queue_ask.Top().price != 52, "We should see the partial filled order in the queue, with price", ob.queue_ask.Top().price)
	Assert(t, ob.queue_bid.Top().size == 1, "We should see the partial filled order in the queue, with size", ob.queue_bid.Top().size)
	Assert(t, ob.queue_bid.Top().filled_pct == 5.0/6.0, "We should see the partial fill order in the queue, with pct", ob.queue_ask.Top().filled_pct)
}

func TestAddMarket(t *testing.T) {
	ob := OrderBook{}
	p := Portfolio{asset: 999_999, cash: 999_999}
	for i := range 50 {
		ob.Add(&Order{
			id:        rand.Uint64(),
			otype:     LIMIT,
			side:      BID,
			size:      5,
			price:     f32(i + 1),
			portfolio: &p,
		})
		ob.Add(&Order{
			id:        rand.Uint64(),
			otype:     LIMIT,
			side:      ASK,
			size:      5,
			price:     f32(101 - i),
			portfolio: &p,
		})
	}

	// Fills 4 orders of size 5 each, at prices (50, 49, 48, 47)
	fill := ob.Add(&Order{
		id:        rand.Uint64(),
		otype:     MARKET,
		side:      ASK,
		size:      20,
		portfolio: &p,
	})
	Assert(t, fill.size == 20, "Should fill all the order")
	Assert(t, fill.price == 48.5, "Should fill the order at correct price", fill.price)
	Assert(t, fill.filled_pct == 1, "Should fill all filled_pct", fill.filled_pct)
	Assert(t, ob.queue_bid.Len() == 46, "Should take 4 orders off the BID queue")
	Assert(t, ob.queue_ask.Len() == 50, "Should leave the ASK queue intact")

	// Fills 6 orders of size 5 each, at prices (52, 53, 54, 55, 56, 57)
	fill = ob.Add(&Order{
		id:        rand.Uint64(),
		otype:     MARKET,
		side:      BID,
		size:      30,
		portfolio: &p,
	})
	Assert(t, fill.size == 30, "Should fill all the order")
	Assert(t, fill.price == 54.5, "Should fill the order at correct price", fill.price)
	Assert(t, fill.filled_pct == 1, "Should fill all filled_pct", fill.filled_pct)
	Assert(t, ob.queue_bid.Len() == 46, "Should leave the BID queue intact")
	Assert(t, ob.queue_ask.Len() == 44, "Should take 6 orders off the BID queue")

	// Partial fills 2 orders of size 5+3, at prices (58, 59)
	fill = ob.Add(&Order{
		id:        rand.Uint64(),
		otype:     MARKET,
		side:      BID,
		size:      8,
		portfolio: &p,
	})
	Assert(t, fill.size == 8, "Should fill all the order", fill.size)
	Assert(t, fill.price == 58.375, "Should fill the order at correct price", fill.price)
	Assert(t, fill.filled_pct == 1, "Should fill all filled_pct", fill.filled_pct)
	Assert(t, ob.queue_ask.Top().size == 2, "Should leave first order size partially filled", ob.queue_ask.Top().size)
	Assert(t, ob.queue_ask.Top().filled_pct == 0.6, "Should leave first order 'filled' partially filled", ob.queue_ask.Top().filled_pct)
}

// With 1m Orders, takes about 160s to execute
func xTestStress(t *testing.T) {
	ob := OrderBook{}
	p := Portfolio{asset: 999_999, cash: 999_999}
	for range 1000 {
		var oside OrderSide
		if rand.Float32() > 0.5 {
			oside = BID
		} else {
			oside = ASK
		}

		ob.Add(&Order{
			id:        rand.Uint64(),
			otype:     LIMIT,
			side:      oside,
			size:      i32(rand.Int() / 10000),
			price:     rand.Float32(),
			portfolio: &p,
		})
	}
	Assert(t, true, "")
}

func xTestDisplay(t *testing.T) {
	ob := OrderBook{}
	p := Portfolio{asset: 999_999, cash: 999_999}
	for range 100 {
		var oside OrderSide
		if rand.Float32() > 0.5 {
			oside = BID
		} else {
			oside = ASK
		}

		ob.Add(&Order{
			id:        rand.Uint64(),
			otype:     LIMIT,
			side:      oside,
			size:      rand.Int31n(10),
			price:     f32(Truncate(rand.Float64(), 2)),
			portfolio: &p,
		})
	}
	ob.PPrint()
	Assert(t, true, "")
}
