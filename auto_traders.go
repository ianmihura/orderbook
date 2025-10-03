package main

import (
	"math/rand"
	"sync"
	"time"
)

// type TraderSpec struct {
// 	condition func() bool
// 	size      func() i32
// 	price     func() i32
// 	side      func() OrderSide
// 	otype     func() OrderType
// 	wait      int
// }

// Returns a clean orderbook with n random orders,
// revolving around 10
func BootOrderbook(n int, lock *sync.Mutex) *OrderBook {
	orderbook := OrderBook{}
	orderbook.lock = lock
	p := Portfolio{asset: 999_999, cash: 999_999}

	for range n {
		var oside OrderSide
		var price f32

		/*
			To get a rough shape of an ordrebook:
			we assume we want something like inverse norm dist,
			revolving around least likely mean_ob=mid_price,
			with a defined range of (2*std_ob,-2*std_ob).
			Therefore we get 2 dist, with mean = mean_ob±(2*std) =>
			10±2 => 8 & 12 (with mean_ob=10, std_ob=1).
			In each case we get half the dist (the half that's nearer
			to the mean_ob)

			The two dist:
			- N(8,1) only the positive half (X>mean)
			- N(12,1) only the negative half (X<mean)
		*/
		if rand.Float32() > 0.5 {
			oside = BID
			price = PHalfNormFloat32T(8, 1.0, 2)
		} else {
			oside = ASK
			price = NHalfNormFloat32T(12, 1.0, 2)
		}

		orderbook.Add(&Order{
			id:        rand.Uint64(),
			portfolio: &p,
			otype:     LIMIT,
			side:      oside,
			size:      rand.Int31n(10) + 1,
			price:     price,
		})
	}

	return &orderbook
}

// Market Maker adds a limit order near the midprice.
// Executes every second, if the spread is >0.1%
func addMM(orderbook *OrderBook, portfolio *Portfolio) {
	for {
		// time.Sleep(time.Second)
		if orderbook.Spread() > 0.001 {
			order := &Order{
				id:        rand.Uint64(),
				otype:     LIMIT,
				side:      RandChoice(BID, ASK),
				size:      rand.Int31n(10) + 1,
				price:     NormFloat32T(orderbook.Midprice(), 0.1, 2),
				portfolio: portfolio,
			}
			orderbook.lock.Lock()
			orderbook.Add(order)
			orderbook.lock.Unlock()
		}
	}
}

func addRetail(orderbook *OrderBook, portfolio *Portfolio) {
	for {
		time.Sleep(time.Second / 3)
		mid := orderbook.Midprice()
		order := &Order{
			id:        rand.Uint64(),
			otype:     RandChoice(LIMIT, MARKET),
			side:      RandChoice(BID, ASK),
			size:      rand.Int31n(3) + 1,           // range(1,3)
			price:     NormFloat32T(mid, mid/10, 2), // std ~ 1
			portfolio: portfolio,
		}
		orderbook.lock.Lock()
		orderbook.Add(order)
		orderbook.lock.Unlock()
	}
}

func addInstitutional(orderbook *OrderBook, portfolio *Portfolio) {
	for {
		time.Sleep(time.Second * time.Duration(10))
		mid := orderbook.Midprice()
		order := &Order{
			id:        rand.Uint64(),
			otype:     LIMIT,
			side:      RandChoice(BID, ASK),
			size:      rand.Int31n(100) + 1,         // range(1,100)
			price:     NormFloat32T(mid, mid/20, 2), // std ~ 0.5
			portfolio: portfolio,
		}
		orderbook.lock.Lock()
		orderbook.Add(order)
		orderbook.lock.Unlock()
	}
}

func addHolder(orderbook *OrderBook, portfolio *Portfolio) {
	for {
		time.Sleep(time.Second * time.Duration(20))
		order := &Order{
			id:        rand.Uint64(),
			otype:     MARKET,
			side:      BID,
			size:      rand.Int31n(10) + 20, // range(20,30)
			portfolio: portfolio,
		}
		orderbook.lock.Lock()
		orderbook.Add(order)
		orderbook.lock.Unlock()
	}
}

func AddAutoTraders(orderbook *OrderBook) *[]*Portfolio {
	autos := []*Portfolio{}
	// Adds 2 market makers
	for range 2 {
		portfolio := Portfolio{asset: 100_000, cash: 1_000_000}
		go addMM(orderbook, &portfolio)
		autos = append(autos, &portfolio)
	}

	// Adds 5 retail intraday, noise traders, many random small orders
	for range 5 {
		portfolio := Portfolio{asset: 1_000, cash: 10_000}
		go addRetail(orderbook, &portfolio)
		autos = append(autos, &portfolio)
	}

	// Adds 2 institutions, ocational random big orders
	for range 2 {
		portfolio := Portfolio{asset: 100_000, cash: 1_000_000}
		go addInstitutional(orderbook, &portfolio)
		autos = append(autos, &portfolio)
	}

	// Adds 1 very ocational long term buy-and-holder
	for range 1 {
		portfolio := Portfolio{asset: 0, cash: 200_000}
		go addHolder(orderbook, &portfolio)
		autos = append(autos, &portfolio)
	}

	// TODO add traders bull/bear

	return &autos
}
