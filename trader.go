package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// var i int
	// fmt.Scan(&i)
	// fmt.Println(i)

	ob := OrderBook{}
	for i := range 1000 {
		var oside OrderSide
		if rand.Float32() > 0.5 {
			oside = BID
		} else {
			oside = ASK
		}

		new_order := Order{
			id:    rand.Uint64(),
			otype: LIMIT,
			side:  oside,
			size:  rand.Int31n(10) + 1,
			price: f32(Truncate(rand.Float64(), 2)),
		}
		// fmt.Print("\033[H\033[2J")
		if i%10 == 0 {
			new_order.Print()
			ob.Add(&new_order)
			ob.PPrint()
			fmt.Println("-----------------------")
		} else {
			ob.Add(&new_order)
		}
		time.Sleep(time.Millisecond * 100)
	}
	ob.PPrint()
}

// TODO boot orderbook
// TODO simulate active trading with multiple threads
