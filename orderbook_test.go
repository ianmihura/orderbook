package main

import (
	"cmp"
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

func TestGetQueues(t *testing.T) {
	oq_ask := []Order{{}}
	oq_bid := []Order{{}, {}}
	ob := OrderBook{queue_ask: oq_ask, queue_bid: oq_bid}

	q_bid := *ob.GetQueue(BID)
	q_ask := *ob.GetQueue(ASK)

	Assert(t, len(q_bid) == 2, "Got wrong queue ", q_bid)
	Assert(t, len(q_ask) == 1, "Got wrong queue ", q_ask)

	q_ask_flip := *ob.GetQueueFlip(BID)
	q_bid_flip := *ob.GetQueueFlip(ASK)

	Assert(t, len(q_bid_flip) == 2, "Got wrong queue ", q_bid_flip)
	Assert(t, len(q_ask_flip) == 1, "Got wrong queue ", q_ask_flip)
}

func TestRemove(t *testing.T) {
	oq_bid := []Order{{id: 1}, {id: 2}, {id: 3}}
	ob := OrderBook{queue_bid: oq_bid}

	deleted_order, err := ob.Remove(&oq_bid[1])

	Assert(t, err == nil, err)
	Assert(t, deleted_order.id == 2, "Identified wrong order to delete")
	Assert(t, len(ob.queue_bid) == 2, "Length of queue should be 2")
	Assert(t,
		ob.queue_bid[0].id == 1 && ob.queue_bid[1].id == 3,
		"Deleted wrong Order",
	)
}

func TestAddLimit(t *testing.T) {
	ob := OrderBook{}
	for range 100 {
		var oSide OrderSide
		if rand.Float32() > 0.5 {
			oSide = BID
		} else {
			oSide = ASK
		}

		o := Order{
			id:      rand.Uint64(),
			otype:   LIMIT,
			side:    oSide,
			price:   rand.Float32() + float32(oSide), // ensures we dont autofill limit orders
			created: time.Now().Add(time.Duration(rand.Uint64())),
		}
		ob.Add(&o)
	}

	Assert(t, (len(ob.queue_ask)+len(ob.queue_bid)) == 100, "Should insert 100 elements")
	Assert(t, !slices.IsSortedFunc(ob.queue_ask, func(a, b Order) int {
		return cmp.Compare(a.id, b.id)
	}), "ASK queue Should not be sorted by ID")
	Assert(t, !slices.IsSortedFunc(ob.queue_bid, func(a, b Order) int {
		return cmp.Compare(a.id, b.id)
	}), "BID queue Should not be sorted by ID")
	Assert(t, slices.IsSortedFunc(ob.queue_ask, func(a, b Order) int {
		return cmp.Compare(a.price, b.price)
	}), "ASK queue Should be sorted by price")
	Assert(t, IsSortedFuncDesc(ob.queue_bid, func(a, b Order) int {
		return cmp.Compare(a.price, b.price)
	}), "BID queue Should be sorted by price")

	// TODO not yet programmed
	// Assert(t, slices.IsSortedFunc(ob.queue_ask, func(a, b Order) int {
	// 	return a.created.Compare(b.created)
	// }), "ASK queue Should be sorted by created datetime")
	// Assert(t, slices.IsSortedFunc(ob.queue_bid, func(a, b Order) int {
	// 	return a.created.Compare(b.created)
	// }), "BID queue Should be sorted by created datetime")
}

func TestAddAutofillLimit(t *testing.T) {
	ob := OrderBook{}
	for i := range 50 {
		ob.Add(&Order{
			id:    rand.Uint64(),
			otype: LIMIT,
			side:  ASK,
			size:  5,
			price: f32(101 - i),
		})
	}

	// Partial fills 1 order of size 5, at prices (52)
	fill := ob.Add(&Order{
		id:    rand.Uint64(),
		otype: LIMIT,
		side:  BID,
		size:  6,
		price: 52,
	})
	Assert(t, fill.size == 5, "Should fill all the order", fill.size)
	Assert(t, fill.price == 52, "Should fill the order at correct price", fill.price)
	Assert(t, fill.filledPct == 5.0/6.0, "Should fill all filledPct", fill.filledPct)
	Assert(t, ob.queue_bid[0].size == 1, "We should see the rest of the order in the queue")
	Assert(t, ob.queue_bid[0].filledPct == 5.0/6.0, "We should see the rest of the order in the queue", ob.queue_ask[0].filledPct)
}

func TestAddMarket(t *testing.T) {
	ob := OrderBook{}
	for i := range 50 {
		ob.Add(&Order{
			id:    rand.Uint64(),
			otype: LIMIT,
			side:  BID,
			size:  5,
			price: f32(i + 1),
		})
		ob.Add(&Order{
			id:    rand.Uint64(),
			otype: LIMIT,
			side:  ASK,
			size:  5,
			price: f32(101 - i),
		})
	}

	// Fills 4 orders of size 5 each, at prices (50, 49, 48, 47)
	fill := ob.Add(&Order{
		id:    rand.Uint64(),
		otype: MARKET,
		side:  ASK,
		size:  20,
	})
	Assert(t, fill.size == 20, "Should fill all the order")
	Assert(t, fill.price == 48.5, "Should fill the order at correct price", fill.price)
	Assert(t, fill.filledPct == 1, "Should fill all filledPct", fill.filledPct)
	Assert(t, len(ob.queue_bid) == 46, "Should take 4 orders off the BID queue")
	Assert(t, len(ob.queue_ask) == 50, "Should leave the ASK queue intact")

	// Fills 6 orders of size 5 each, at prices (52, 53, 54, 55, 56, 57)
	fill = ob.Add(&Order{
		id:    rand.Uint64(),
		otype: MARKET,
		side:  BID,
		size:  30,
	})
	Assert(t, fill.size == 30, "Should fill all the order")
	Assert(t, fill.price == 54.5, "Should fill the order at correct price", fill.price)
	Assert(t, fill.filledPct == 1, "Should fill all filledPct", fill.filledPct)
	Assert(t, len(ob.queue_bid) == 46, "Should leave the BID queue intact")
	Assert(t, len(ob.queue_ask) == 44, "Should take 6 orders off the BID queue")

	// Partial fills 2 orders of size 5+3, at prices (58, 59)
	fill = ob.Add(&Order{
		id:    rand.Uint64(),
		otype: MARKET,
		side:  BID,
		size:  8,
	})
	Assert(t, fill.size == 8, "Should fill all the order", fill.size)
	Assert(t, fill.price == 58.375, "Should fill the order at correct price", fill.price)
	Assert(t, fill.filledPct == 1, "Should fill all filledPct", fill.filledPct)
	Assert(t, ob.queue_ask[0].size == 2, "Should leave first order size partially filled", ob.queue_ask[0].size)
	Assert(t, ob.queue_ask[0].filledPct == 0.6, "Should leave first order 'filled' partially filled", ob.queue_ask[0].filledPct)
}

// With 1m Orders, takes about 100s to execute
func xTestStress(t *testing.T) {
	ob := OrderBook{}
	for range 1_000_000 {
		var oSide OrderSide
		if rand.Float32() > 0.5 {
			oSide = BID
		} else {
			oSide = ASK
		}

		ob.Add(&Order{
			id:    rand.Uint64(),
			otype: LIMIT,
			side:  oSide,
			size:  rand.Int32() / 10000,
			price: rand.Float32(),
			// created: time.Now().Add(time.Duration(rand.Uint64())),
		})
	}
	Assert(t, true, "")
}

func TestDisplay(t *testing.T) {
	ob := OrderBook{}
	for range 100 {
		var oSide OrderSide
		if rand.Float32() > 0.5 {
			oSide = BID
		} else {
			oSide = ASK
		}

		// TODO try some random func thats not so smooth
		ob.Add(&Order{
			id:    rand.Uint64(),
			otype: LIMIT,
			side:  oSide,
			size:  rand.Int32() / 10000000,
			price: rand.Float32(),
			// created: time.Now().Add(time.Duration(rand.Uint64())),
		})
	}
	ob.PPrint()
	Assert(t, true, "")
}
