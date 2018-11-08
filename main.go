package main

import (
	"os"
	"strings"

	"github.com/jasonlvhit/gocron"

	"github.com/hexoul/go-coinmarketcap/statistics"
	"github.com/hexoul/go-coinmarketcap/types"
)

var (
	targetSymbol string
	targetAddr   string
	accessKey    map[string]string
	secretKey    map[string]string
)

func init() {
	accessKey = map[string]string{}
	secretKey = map[string]string{}

	for _, val := range os.Args {
		arg := strings.Split(val, "=")
		if len(arg) < 2 {
			continue
		} else if arg[0] == "-targetSymbol" {
			targetSymbol = arg[1]
		} else if arg[0] == "-targetAddr" {
			targetAddr = arg[1]
		} else if strings.Contains(arg[0], "accesskey") {
			accessKey[strings.Split(arg[0], ":")[0][1:]] = arg[1]
		} else if strings.Contains(arg[0], "secretkey") {
			secretKey[strings.Split(arg[0], ":")[0][1:]] = arg[1]
		}
	}

	if targetSymbol == "" || targetAddr == "" {
		panic("TARGET INFO REQUIRED")
	}
}

func main() {
	quotes := "USD,BTC,ETH"

	statistics.GatherCryptoQuote(&types.Options{
		Symbol:  targetSymbol,
		Convert: quotes,
	}, gocron.Every(30).Seconds())

	statistics.GatherTokenMetric(targetSymbol, targetAddr, gocron.Every(2).Seconds())

	<-gocron.Start()
}
