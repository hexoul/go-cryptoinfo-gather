package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jasonlvhit/gocron"

	"github.com/hexoul/go-coinmarketcap"
	"github.com/hexoul/go-coinmarketcap/statistics"
	"github.com/hexoul/go-coinmarketcap/types"
)

var (
	targetSymbol    string
	targetAddr      string
	kucoinAccesskey string
	kucoinSecretkey string
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
		} else if arg[0] == "-kucoinAccesskey" {
			kucoinAccesskey = arg[1]
		} else if arg[0] == "-kucoinSecret" {
			kucoinSecretkey = arg[1]
		}
	}

	if targetSymbol == "" || targetAddr == "" {
		panic("TARGET INFO REQUIRED")
	}
}

func main() {
	if ret, err := coinmarketcap.GetInstance().CryptoInfo(&types.Options{
		Symbol: targetSymbol,
	}); err == nil {
		fmt.Println(ret.CryptoInfo[targetSymbol].Name)
	}

	statistics.GatherTokenMetric(targetSymbol, targetAddr, gocron.Every(2).Seconds())
	<-gocron.Start()
}
