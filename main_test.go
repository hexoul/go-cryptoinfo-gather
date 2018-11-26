package main

import (
	kucoin "github.com/eeonevision/kucoin-go"
	abcc "github.com/hexoul/go-abcc"
	coinsuper "github.com/hexoul/go-coinsuper"
)

var (
	testAccessKey = map[string]string{
		"kucoin":    "YOUR_ACCESS_KEY",
		"coinsuper": "YOUR_ACCESS_KEY",
		"abcc":      "YOUR_ACCESS_KEY",
	}
	testSecretKey = map[string]string{
		"kucoin":    "YOUR_SECRET_KEY",
		"coinsuper": "YOUR_SECRET_KEY",
		"abcc":      "YOUR_SECRET_KEY",
	}

	testClients Clients
)

func init() {
	testClients.kucoin = kucoin.New(testAccessKey["kucoin"], testSecretKey["kucoin"])
	testClients.abcc = abcc.GetInstanceWithKey(testAccessKey["abcc"], testSecretKey["abcc"])
	testClients.coinsuper = coinsuper.GetInstanceWithKey(testAccessKey["coinsuper"], testSecretKey["coinsuper"])
}
