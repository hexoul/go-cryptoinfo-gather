package main

import (
	abcc "github.com/hexoul/go-abcc"
	coinsuper "github.com/hexoul/go-coinsuper"
	kucoin "github.com/hexoul/go-kucoin"
	upbit "github.com/hexoul/go-upbit"
	bittrex "github.com/toorop/go-bittrex"
)

var (
	testAccessKey = map[string]string{
		"kucoin":    "YOUR_ACCESS_KEY",
		"coinsuper": "YOUR_ACCESS_KEY",
		"abcc":      "YOUR_ACCESS_KEY",
		"bittrex":   "YOUR_ACCESS_KEY",
		"upbit":     "YOUR_ACCESS_KEY",
	}
	testSecretKey = map[string]string{
		"kucoin":    "YOUR_SECRET_KEY",
		"coinsuper": "YOUR_SECRET_KEY",
		"abcc":      "YOUR_SECRET_KEY",
		"bittrex":   "YOUR_SECRET_KEY",
		"upbit":     "YOUR_SECRET_KEY",
	}
	testPassPhrase = map[string]string{
		"kucoin": "YOUR_PASS_PHRASE",
	}

	testClients Clients
)

func init() {
	targetSymbol = "META"
	testClients.kucoin = kucoin.GetInstanceWithKey(testAccessKey["kucoin"], testSecretKey["kucoin"], testPassPhrase["kucoin"])
	testClients.abcc = abcc.GetInstanceWithKey(testAccessKey["abcc"], testSecretKey["abcc"])
	testClients.coinsuper = coinsuper.GetInstanceWithKey(testAccessKey["coinsuper"], testSecretKey["coinsuper"])
	testClients.bittrex = bittrex.New(testAccessKey["bittrex"], testSecretKey["bittrex"])
	testClients.upbit = upbit.GetInstanceWithKey(testAccessKey["upbit"], testSecretKey["upbit"])
}
