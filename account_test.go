package main

import (
	"testing"
)

func TestKucoinBalance(t *testing.T) {
	if bal, err := testClients.kucoin.GetCoinBalance("USDT"); err != nil {
		t.FailNow()
	} else {
		t.Logf("%f %f\n", bal.Balance, bal.FreezeBalance)
	}
}

func TestKucoinListMergedDealtOrders(t *testing.T) {
	if _, err := testClients.kucoin.ListMergedDealtOrders("ETH-BTC", "BUY", 20, 1, 0, 0); err != nil {
		t.FailNow()
	}
}

func TestGetBalances(t *testing.T) {
	testClients.GetBalances()
}

func TestGetTrades(t *testing.T) {
	testClients.GetTrades()
}

func TestCheckExistOrder(t *testing.T) {
	oID := "123"
	if checkExistOrder(oID) {
		t.FailNow()
	}
	if !checkExistOrder(oID) {
		t.FailNow()
	}
}

func TestDuplicatedGetTrades(t *testing.T) {
	testClients.GetTrades()
	testClients.GetTrades()
}
