package main

type Transaction struct {
	active  Order
	passive []Order
}

type TransactionHistory struct {
	orders []Transaction
}
