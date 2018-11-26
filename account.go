package main

import (
	"fmt"
	"io"
	"os"
	"time"

	kucoin "github.com/eeonevision/kucoin-go"

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

	balanceLogger.WithFields(log.Fields{
		"exchange": "kucoin",
		"meta":     meta,
		"eth":      eth,
		"btc":      btc,
	}).Info("GatherBalance")
	return
}

// GetBalances records balances
func GetBalances(c *Clients) {
	getKucoinBalnace(c.kucoin)
}
