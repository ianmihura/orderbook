package main

import (
	"testing"
)

func resetTxHistory() {
	txLock.Lock()
	defer txLock.Unlock()
	_tx_history = nil
}

func TestTransactionHistoryHappyPath(t *testing.T) {
	// Reset the history singleton before the test (seems that its a good practive)
	resetTxHistory()

	activeOrder1 := Order{
		id:    1001,
		otype: LIMIT,
		side:  BID,
		size:  10,
		price: 100.0,
	}
	passiveOrder1 := Order{
		id:    2001,
		otype: LIMIT,
		side:  ASK,
		size:  5,
		price: 99.0,
	}

	report1Active := FillReport{
		is_active: true,
		size:      5,
		price:     99.0,
	}
	report1Passive := FillReport{
		is_active: false,
		size:      5,
		price:     99.0,
	}

	activeOrder2 := Order{
		id:    1002,
		otype: LIMIT,
		side:  ASK,
		size:  20,
		price: 101.0,
	}
	passiveOrder2 := Order{
		id:    2002,
		otype: LIMIT,
		side:  BID,
		size:  15,
		price: 101.5,
	}

	report2Active := FillReport{
		is_active: true,
		size:      15,
		price:     101.5,
	}

	// Test AddToTxHistory
	AddToTxHistory(passiveOrder1, activeOrder1)
	AddToTxHistory(passiveOrder2, activeOrder2)

	tx1 := getTxById(activeOrder1.id)
	Assert(t, tx1 != nil, "Transaction 1 should exist")
	Assert(t, tx1.active.id == activeOrder1.id, "Active order ID should match")
	Assert(t, len(tx1.passive) == 1, "Should have 1 passive order initially")
	Assert(t, tx1.passive[0].id == passiveOrder1.id, "Passive order ID should match")

	tx2 := getTxById(activeOrder2.id)
	Assert(t, tx2 != nil, "Transaction 2 should exist")
	Assert(t, tx2.active.id == activeOrder2.id, "Active order ID should match")
	Assert(t, len(tx2.passive) == 1, "Should have 1 passive order initially")

	passiveOrder1_2 := Order{
		id:    2003,
		otype: LIMIT,
		side:  ASK,
		size:  3,
		price: 99.5,
	}
	AddToTxHistory(passiveOrder1_2, activeOrder1)
	Assert(t, len(tx1.passive) == 2, "Should have 2 passive orders after second addition")
	Assert(t, tx1.passive[1].id == passiveOrder1_2.id, "Second passive order ID should match")

	// Test AddFillReport
	AddFillReport(activeOrder1.id, &report1Active)
	AddFillReport(activeOrder1.id, &report1Passive)
	AddFillReport(activeOrder2.id, &report2Active)

	Assert(t, len(tx1.reports) == 2, "Transaction 1 should have 2 reports")
	Assert(t, tx1.reports[0].is_active == true, "First report for tx1 should be active")
	Assert(t, tx1.reports[1].is_active == false, "Second report for tx1 should be passive")

	Assert(t, len(tx2.reports) == 1, "Transaction 2 should have 1 report")
	Assert(t, tx2.reports[0].is_active == true, "Report for tx2 should be active")

	//    Test GetAvgPriceWeighted
	// We have two transactions with active reports:
	// tx1: price 99.0, size 5. Weighted price: 99.0 * 5 = 495.0
	// tx2: price 101.5, size 15. Weighted price: 101.5 * 15 = 1522.5
	// Total weighted price: 495.0 + 1522.5 = 2017.5
	// Total weight (sum_weight): 5 + 15 = 20
	// Avg Weighted Price: 2017.5 / 20 = 100.875

	expectedAvgWeighted := f32(100.875)
	actualAvgWeighted := GetAvgPriceWeighted()

	delta := f32(0.0001) // Fucking float comparison
	Assert(t, Abs(actualAvgWeighted-expectedAvgWeighted) < delta,
		"GetAvgPriceWeighted should be correct", actualAvgWeighted)

	//    Test GetAvgPrice
	// We have two transactions with active reports:
	// tx1: price 99.0
	// txx2: price 101.5
	// Total price: 99.0 + 101.5 = 200.5
	// Avg Price: 200.5 / 2 = 100.25

	expectedAvgPrice := f32(100.25)
	actualAvgPrice := GetAvgPrice()
	Assert(t, Abs(actualAvgPrice-expectedAvgPrice) < delta,
		"GetAvgPrice should be correct", actualAvgPrice)
}
