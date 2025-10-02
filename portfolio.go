package main

type Portfolio struct {
	cash    f32
	asset   i32
	is_user bool
}

// Checks if current balance (either cash or asset)
// can accomodate the incoming order.
// Only accomodates full orders, does not fill orders partially.
func (portfolio *Portfolio) CanFill(order *Order) bool {
	if order.side == BID {
		// I'm selling
		return portfolio.asset >= order.size
	} else {
		// I'm buying
		order_cost := order.price * f32(order.size)
		return portfolio.cash >= order_cost
	}
}

// Incoming order is the other side of the trade
func (portfolio *Portfolio) Fill(order *Order) {
	if order.side == BID {
		// I'm selling
		portfolio.cash += order.price * f32(order.size)
		portfolio.asset -= order.size
	} else {
		// I'm buying
		portfolio.cash -= order.price * f32(order.size)
		portfolio.asset += order.size
	}
}
