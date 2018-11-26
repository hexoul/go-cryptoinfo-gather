package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	kucoin "github.com/eeonevision/kucoin-go"
	abcc "github.com/hexoul/go-abcc"
	coinsuper "github.com/hexoul/go-coinsuper"

	log "github.com/sirupsen/logrus"
)

var (
	balanceLogger *log.Logger
)

func init() {
	// Initialize logger
	balanceLogger = log.New()

	// Default configuration
	balanceLogger.Formatter = &log.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}
	if f, err := os.OpenFile("./balance.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
		balanceLogger.Out = io.MultiWriter(f, os.Stdout)
	} else {
		fmt.Print("Failed to open log file: you can miss important log")
		balanceLogger.Out = os.Stdout
	}
	balanceLogger.SetLevel(log.InfoLevel)
}

func logBalance(exchange string, meta, eth, btc interface{}) {
	balanceLogger.WithFields(log.Fields{
		"exchange": exchange,
		"meta":     meta,
		"eth":      eth,
		"btc":      btc,
	}).Info("GatherBalance")
}

func getKucoinBalnace(k *kucoin.Kucoin) (meta, eth, btc float64) {
	if k == nil {
		return
	}
	if bal, err := k.GetCoinBalance("META"); err == nil {
		meta = bal.Balance + bal.FreezeBalance
	}
	if bal, err := k.GetCoinBalance("ETH"); err == nil {
		eth = bal.Balance + bal.FreezeBalance
	}
	if bal, err := k.GetCoinBalance("BTC"); err == nil {
		btc = bal.Balance + bal.FreezeBalance
	}
	logBalance("kucoin", meta, eth, btc)
	return
}

func getCoinsuperBalnace() (meta, eth, btc string) {
	if info, err := coinsuper.GetInstance().UserAssetInfo(nil); err == nil {
		meta = info.Assets["META"].Total
		eth = info.Assets["ETH"].Total
		btc = info.Assets["BTC"].Total
		logBalance("coinsuper", meta, eth, btc)
	}
	return
}

func sumStrFloat(s1, s2 string) (sum float64) {
	f1, err1 := strconv.ParseFloat(s1, 64)
	f2, err2 := strconv.ParseFloat(s2, 64)
	if err1 == nil && err2 == nil {
		sum = f1 + f2
	}
	return
}

func getAbccBalnace() (meta, eth, btc float64) {
	if me, err := abcc.GetInstance().Me(nil); err == nil {
		for _, v := range me.Accounts {
			if v.Currency == "meta" {
				meta = sumStrFloat(v.Balance, v.Locked)
			} else if v.Currency == "eth" {
				eth = sumStrFloat(v.Balance, v.Locked)
			} else if v.Currency == "btc" {
				btc = sumStrFloat(v.Balance, v.Locked)
			}
		}
		logBalance("abcc", meta, eth, btc)
	}
	return
}

// GetBalances records balances
func GetBalances(c *Clients) {
	getKucoinBalnace(c.kucoin)
	getCoinsuperBalnace()
	getAbccBalnace()
}
