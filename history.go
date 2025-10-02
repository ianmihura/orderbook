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
	txs map[u64]*Transaction
}

var txLock = &sync.Mutex{}

// Singleton pattern
var _tx_history *TransactionHistory

func getTxHistory() *TransactionHistory {
	if _tx_history == nil {
		txLock.Lock()
		defer txLock.Unlock()
		if _tx_history == nil {
			_tx_history = &TransactionHistory{
				txs: map[u64]*Transaction{},
			}
		}
	}
	return _tx_history
}

func getTxById(id u64) *Transaction {
	txs := getTxHistory().txs
	var tx *Transaction
	txLock.Lock()
	tx = txs[id]
	txLock.Unlock()
	return tx
}

// We save a copy of the orders, before they get changed
func AddToTxHistory(passive_order, active_order Order) {
	tx := getTxById(active_order.id)
	if tx == nil {
		tx := Transaction{
			active:  active_order,
			passive: []Order{passive_order},
		}

		tx_history := getTxHistory()
		txLock.Lock()
		tx_history.txs[active_order.id] = &tx
		txLock.Unlock()
	} else {
		txLock.Lock()
		tx.passive = append(tx.passive, passive_order)
		txLock.Unlock()
	}
}

// We save a copy of the reports, before they get changed
func AddFillReport(id u64, fill_report *FillReport) {
	tx := getTxById(id)
	for tx == nil {
		// Maybe the tx hasnt been created yet
		time.Sleep(time.Second)
		tx = getTxById(id)
	}
	txLock.Lock()
	tx.reports = append(tx.reports, *fill_report)
	txLock.Unlock()
}

// Gets the first found active report from a list of reports.
// Returns nil if it finds none.
func getActiveReport(reports []FillReport) *FillReport {
	for _, r := range reports {
		if r.is_active {
			return &r
		}
	}
	return nil
}

// Compute the average of all executed prices, weighted by their size
// We keep this function expensive because we wont use it very often
func GetAvgPriceWeighted() f32 {
	var weighted_price f32
	var sum_weight i32
	txs := getTxHistory().txs // Get a copy, just in case
	for _, tx := range txs {
		active_report := getActiveReport(tx.reports)
		if active_report == nil {
			continue
		}
		weighted_price += active_report.price * f32(active_report.size)
		sum_weight += active_report.size
	}
	return weighted_price / f32(sum_weight)
}

// Compute the average of all executed prices
// We keep this function expensive because we wont use it very often
func GetAvgPrice() f32 {
	var total_price f32
	txs := getTxHistory().txs // Get a copy, just in case
	for _, tx := range txs {
		active_report := getActiveReport(tx.reports)
		if active_report == nil {
			continue
		}
		total_price += active_report.price
	}
	return total_price / f32(len(txs))
}

// Compute the average of all executed sizes
// We keep this function expensive because we wont use it very often
func GetAvgSize() i32 {
	var total_size i32
	txs := getTxHistory().txs // Get a copy, just in case
	for _, tx := range txs {
		active_report := getActiveReport(tx.reports)
		if active_report == nil {
			continue
		}
		total_size += active_report.size
	}
	return total_size / i32(len(txs))
}
