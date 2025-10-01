package main

import (
	"sync"
	"time"
)

type Transaction struct {
	active  Order
	passive []Order
	reports []FillReport
}

type TransactionHistory struct {
	orders []Transaction
}

var txLock = &sync.Mutex{}

var _tx_history *TransactionHistory

func GetTxHistory() *TransactionHistory {
	if _tx_history == nil {
		txLock.Lock()
		defer txLock.Unlock()
		if _tx_history == nil {
			_tx_history = &TransactionHistory{}
		}
	}
	return _tx_history
}

// Finds a tx in the history, based on the active id
func FindInTxHistory(id u64) *Transaction {
	// loop from the back: its likely that were working with one of the latest added
	tx_history := GetTxHistory()
	for i := len(tx_history.orders) - 1; i >= 0; i-- {
		if tx_history.orders[i].active.id == id {
			return &tx_history.orders[i]
		}
	}
	return nil
}

// We save a copy of these orders before they get changed
func AddToTxHistory(passive_order, active_order Order) {
	tx_history := GetTxHistory()
	tx := FindInTxHistory(active_order.id)
	if tx == nil {
		tx := Transaction{
			active:  active_order,
			passive: []Order{passive_order},
		}
		tx_history.orders = append(tx_history.orders, tx)
	} else {
		tx.passive = append(tx.passive, passive_order)
	}
}

// We save a copy of these orders before they get changed
func AddFillReport(id u64, fill_report *FillReport) {
	tx := FindInTxHistory(id)
	for tx == nil {
		tx = FindInTxHistory(id)
		// Maybe the tx hasnt been created yet
		time.Sleep(time.Second)
	}
	tx.reports = append(tx.reports, *fill_report)
}
