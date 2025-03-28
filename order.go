package main

import (
	"time"
)

type OrderSide i8
type OrderType i8

const BID OrderSide = 0
const ASK OrderSide = 1
const LIMIT OrderType = 0
const MARKET OrderType = 1

type Order struct {
	id         u64
	otype      OrderType
	side       OrderSide
	size       i32
	price      f32
	created    time.Time
	filled_pct f32
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
	// TODO validate with Portfolio
	// TODO if valid, add to TransactionHistory

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
