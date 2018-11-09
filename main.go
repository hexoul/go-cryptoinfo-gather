package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jasonlvhit/gocron"

	"github.com/hexoul/go-coinmarketcap/statistics"
	"github.com/hexoul/go-coinmarketcap/types"
)

var (
	targetSymbol string
	targetAddr   string
	targetQuotes = "USD"
	targetSlugs  = "binance"
	accessKey    = map[string]string{}
	secretKey    = map[string]string{}
)

func init() {
	for _, val := range os.Args {
		arg := strings.Split(val, "=")
		if len(arg) < 2 {
			continue
		} else if arg[0] == "-targetSymbol" {
			targetSymbol = arg[1]
		} else if arg[0] == "-targetAddr" {
			targetAddr = arg[1]
		} else if arg[0] == "-targetQuotes" {
			targetQuotes = arg[1]
		} else if arg[0] == "-targetSlugs" {
			targetSlugs = arg[1]
		} else if strings.Contains(arg[0], "accesskey") {
			accessKey[strings.Split(arg[0], ":")[0][1:]] = arg[1]
		} else if strings.Contains(arg[0], "secretkey") {
			secretKey[strings.Split(arg[0], ":")[0][1:]] = arg[1]
		}
	}

	// if targetSymbol == "" || targetAddr == "" {
	// 	panic("TARGET INFO REQUIRED")
	// }
}

func main() {
	fmt.Println("Scheduling...")

	// ExchangeMarketPairs
	for _, slug := range strings.Split(targetSlugs, ",") {
		statistics.GatherExchangeMarketPairs(&types.Options{
			Slug:    slug,
			Convert: targetQuotes,
		}, targetSymbol, gocron.Every(2).Minutes())
	}

	// CryptoQuote
	statistics.GatherCryptoQuote(&types.Options{
		Symbol:  targetSymbol,
		Convert: targetQuotes,
	}, gocron.Every(2).Minutes())

	// TokenMetric
	statistics.GatherTokenMetric(targetSymbol, targetAddr, gocron.Every(2).Minutes())

	fmt.Printf("Done\nStarting...\n")
	<-gocron.Start()
}
