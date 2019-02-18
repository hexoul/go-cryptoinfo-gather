package main

import (
	abcc "github.com/hexoul/go-abcc"
	coinsuper "github.com/hexoul/go-coinsuper"
	kucoin "github.com/hexoul/go-kucoin"
)

var (
	testAccessKey = map[string]string{
		"kucoin":    "5c6ac76a1cde7830729e25eb",
		"coinsuper": "YOUR_ACCESS_KEY",
		"abcc":      "YOUR_ACCESS_KEY",
	}
	testSecretKey = map[string]string{
		"kucoin":    "5c0910e-b3ad-41fe-af34-0164d9143322",
		"coinsuper": "YOUR_SECRET_KEY",
		"abcc":      "YOUR_SECRET_KEY",
	}
	testPassPhrase = map[string]string{
		"kucoin": "MetadiumToTheMoon",
	}

	testClients Clients
)

func init() {
	targetSymbol = "META"
	testClients.kucoin = kucoin.GetInstanceWithKey(testAccessKey["kucoin"], testSecretKey["kucoin"], testPassPhrase["kucoin"])
	testClients.abcc = abcc.GetInstanceWithKey(testAccessKey["abcc"], testSecretKey["abcc"])
	testClients.coinsuper = coinsuper.GetInstanceWithKey(testAccessKey["coinsuper"], testSecretKey["coinsuper"])
}
