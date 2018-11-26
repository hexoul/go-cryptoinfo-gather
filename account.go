package main

func getKucoinBalnace() (meta, eth, btc float64) {
	if bal, err := clients.kucoin.GetCoinBalance("META"); err == nil {
		meta = bal.Balance + bal.FreezeBalance
	}
	if bal, err := clients.kucoin.GetCoinBalance("ETH"); err == nil {
		eth = bal.Balance + bal.FreezeBalance
	}
	if bal, err := clients.kucoin.GetCoinBalance("BTC"); err == nil {
		btc = bal.Balance + bal.FreezeBalance
	}
	return
}

// GetBalances records balances
func GetBalances() {
	getKucoinBalnace()
}
