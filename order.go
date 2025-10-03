package main

type OrderSide bool
type OrderType rune

const BID OrderSide = false
const ASK OrderSide = true
const LIMIT OrderType = 'L'
const MARKET OrderType = 'M'
const MID OrderType = 'D'
const VWAP OrderType = 'V'
const TWAP OrderType = 'T'

type Order struct {
	portfolio *Portfolio
	id        u64
	otype     OrderType
	side      OrderSide
	size      i32
	price     f32
	// created    time.Time
	filled_pct f32
	orderbook  *OrderBook
}

// created: time.Now().Add(time.Duration(rand.Uint64())),

/*
OB -> calls Fill on Order
	Order -> calls Fill on Portfolio
		Portfolio -> validates if can execute and modifies portfolio
		<- returns ok / nok
	<- if nok, returns fillreport
	if ok -> save to TransactionHistory, edit Orders
	<- returns fillreport
*/

// Fills both orders up to the min(order_active.size, order_passive.size).
// Will check in the portfolio if there is enough assets or cash.
//
// Note: portfolio either fills orders fully, or does not fill,
// portfolios do not fill orders partially: if you have insuficient balance
// in a portfolio, it will not fill the order partially
//
// Returns a FillReport, corresponding to the passive_order (as a convention).
func (passive_order *Order) Fill(active_order *Order) *FillReport {
	if active_order.otype == MARKET {
		active_order.price = passive_order.price
	}

	can_fill_passive_order := passive_order.portfolio.CanFill(active_order)
	can_fill_active_order := active_order.portfolio.CanFill(passive_order)

	if !can_fill_passive_order || !can_fill_active_order {
		return &FillReport{}
	}

	// Reflecting changes in the portfolio
	// TODO error: sending the full order to be filled by the portfolio, assuming I'll fill it fully
	passive_order.portfolio.Fill(active_order)
	active_order.portfolio.Fill(passive_order)
	// Save the transactions to history before they get changed
	go AddToTxHistory(*passive_order, *active_order)

	fill_report := FillReport{
		price: passive_order.price,
	}
	if active_order.size >= passive_order.size {
		// Filling active_order partially or totally
		active_order.filled_pct = f32(passive_order.size) / f32(active_order.size)
		active_order.size -= passive_order.size

		fill_report.size = passive_order.size
		fill_report.filled_pct = 1
	} else {
		// Filling passive_order partially
		passive_order.filled_pct = f32(active_order.size) / f32(passive_order.size)
		passive_order.size -= active_order.size

		fill_report.size = active_order.size
		fill_report.filled_pct = passive_order.filled_pct
	}

	go AddFillReport(active_order.id, &fill_report)

	return &fill_report
}
