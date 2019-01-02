package main

import (
	"fmt"
	"os"
	"strings"

	kucoin "github.com/eeonevision/kucoin-go"
	abcc "github.com/hexoul/go-abcc"
	coinsuper "github.com/hexoul/go-coinsuper"

	"github.com/jasonlvhit/gocron"

	"github.com/hexoul/go-coinmarketcap/statistics"
	"github.com/hexoul/go-coinmarketcap/types"
)

// Clients struct
type Clients struct {
	abcc      *abcc.Client
	kucoin    *kucoin.Kucoin
	coinsuper *coinsuper.Client
}

var (
	targetSymbol string
	targetAddr   string
	targetQuotes = "USD"
	targetSlugs  = "binance"
	accessKey    = map[string]string{}
	secretKey    = map[string]string{}
	clients      Clients
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

	if accessKey["kucoin"] != "" {
		clients.kucoin = kucoin.New(accessKey["kucoin"], secretKey["kucoin"])
	}
}

func main() {
	fmt.Println("Initializing...")

	// Initialize CryptoQuote
	cryptoQuoteOptions := &types.Options{
		Symbol:  targetSymbol,
		Convert: targetQuotes,
	}
	statistics.TaskGatherCryptoQuote(cryptoQuoteOptions)

	// Initialize ExchangeMarketPairs
	var exchangeMarketPairsOptions []*types.Options
	for i, slug := range strings.Split(targetSlugs, ",") {
		exchangeMarketPairsOptions = append(exchangeMarketPairsOptions, &types.Options{
			Slug:    slug,
			Convert: targetQuotes,
			Limit:   200,
		})
		statistics.TaskGatherExchangeMarketPairs(exchangeMarketPairsOptions[i], slug)
	}

	// Initialize TokenMetric
	statistics.TaskGatherTokenMetric(targetSymbol, targetAddr)

	// Initialize OHLCV
	var ohlcvOptions []*types.Options
	for _, quote := range strings.Split(targetQuotes, ",") {
		ohlcvOptions = append(ohlcvOptions, &types.Options{
			Symbol:  targetSymbol,
			Convert: quote,
		})
	}

	// Initialize Balance
	clients.GetBalances()

	// Initialize Trade
	clients.GetTrades()

	fmt.Printf("Done\nScheduling...\n")

	// Schedule CryptoQuote
	statistics.GatherCryptoQuote(cryptoQuoteOptions, gocron.Every(10).Minutes())
	statistics.GatherCryptoQuote(cryptoQuoteOptions, gocron.Every(1).Day().At("23:59"))

	// Schedule ExchangeMarketPairs
	for _, option := range exchangeMarketPairsOptions {
		statistics.GatherExchangeMarketPairs(option, targetSymbol, gocron.Every(10).Minutes())
		statistics.GatherExchangeMarketPairs(option, targetSymbol, gocron.Every(1).Day().At("23:59"))
	}

	// Schedule TokenMetric
	statistics.GatherTokenMetric(targetSymbol, targetAddr, gocron.Every(30).Minutes())
	statistics.GatherTokenMetric(targetSymbol, targetAddr, gocron.Every(1).Day().At("23:59"))

	// Schedule OHLCV
	for _, option := range ohlcvOptions {
		statistics.GatherOhlcv(option, gocron.Every(1).Day().At("10:00"))
	}

	// Schedule Balance
	gocron.Every(10).Minutes().Do(clients.GetBalances)

	// Schedule Trade
	gocron.Every(2).Minutes().Do(clients.GetTrades)

	// Schedule Git commit and push
	gocron.Every(2).Hours().Do(gitPushChanges)

	fmt.Printf("Done\nStart!!\n")
	<-gocron.Start()
}
