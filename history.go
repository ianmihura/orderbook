package main

type Transaction struct {
	active  Order
	passive []Order
}

type TransactionHistory struct {
	orders []Transaction
}

func (transaction_history *TransactionHistory) Append(passive_order, active_order *Order) {
	// TODO search (last one) if there's another active order transaction
	// if ok, append
	// else, create a new one
}
