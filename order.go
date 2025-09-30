package main

import (
	"time"
)

type OrderSide bool
type OrderType i8

const BID OrderSide = false
const ASK OrderSide = true
const LIMIT OrderType = 0
const MARKET OrderType = 1

type Order struct {
	portfolio  *Portfolio
	id         u64
	otype      OrderType
	side       OrderSide
	size       i32
	price      f32
	created    time.Time
	filled_pct f32
	order_book *OrderBook
}

/*
OB -> calls Fill on Order
	-> calls Fill ok Portfolio
		-> validates if can execute and modifies portfolio
		<- returns ok / nok
	if ok -> save to TransactionHistory
	<- returns ok / nok to OB
*/

func (passive_order *Order) Fill(active_order *Order) *FillReport {
	can_fill_passive_order := passive_order.portfolio.CanFill(active_order)
	can_fill_active_order := active_order.portfolio.CanFill(passive_order)

	if can_fill_passive_order && can_fill_active_order {
		// TODO enable rollback (just in case)
		passive_order.portfolio.Fill(active_order)
		active_order.portfolio.Fill(passive_order)

		// TODO modify portfolio valies

		// passive_order.order_book.transaction_history.Append(
		// 	passive_order, active_order,
		// )
		return fill(passive_order, active_order)
	} else {
		return &FillReport{}
	}
}

func fill(passive_order, active_order *Order) *FillReport {
	fill_report := FillReport{
		price: passive_order.price,
	}

	if active_order.size >= passive_order.size {
		// Filling active_order partially
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

	return &fill_report
}
