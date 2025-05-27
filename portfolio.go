package main

type Portfolio struct {
	cash  f32
	asset i32
}

// Checks if current balance (either cash or asset)
// can accomodate the incoming order
func (portfolio *Portfolio) CanFill(order *Order) bool {
	// TODO validate
	return true
}

func (portfolio *Portfolio) Fill(order *Order) {
	// TODO modify the portfolio
}

// TODO rollback last fill
