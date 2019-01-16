package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	abcc "github.com/hexoul/go-abcc"
	abccTypes "github.com/hexoul/go-abcc/types"
	coinsuper "github.com/hexoul/go-coinsuper"
	kucoin "github.com/hexoul/go-kucoin"

	log "github.com/sirupsen/logrus"
)

const (
	existLen = 1000
)

var (
	balanceLogger *log.Logger
	tradeLogger   *log.Logger
	existOrderID  [existLen]string
	existIdx      = 0
)

func checkExistOrder(orderID string) bool {
	for i := 0; i < existLen; i++ {
		if existOrderID[i] == orderID {
			return true
		}
	}
	existOrderID[existIdx] = orderID
	existIdx++
	if existIdx >= existLen {
		existIdx = 0
	}
	return false
}

func init() {
	// Initialize logger
	balanceLogger = log.New()
	tradeLogger = log.New()

	// Set formatter
	jsonFormatter := &log.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}
	balanceLogger.Formatter = jsonFormatter
	tradeLogger.Formatter = jsonFormatter

	// Set writer
	if f, err := os.OpenFile("./balance.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
		balanceLogger.Out = io.MultiWriter(f, os.Stdout)
	} else {
		fmt.Print("Failed to open log file: you can miss important log")
		balanceLogger.Out = os.Stdout
	}
	if f, err := os.OpenFile("./trade.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
		tradeLogger.Out = io.MultiWriter(f, os.Stdout)
	} else {
		fmt.Print("Failed to open log file: you can miss important log")
		tradeLogger.Out = os.Stdout
	}

	// Set level
	balanceLogger.SetLevel(log.InfoLevel)
	tradeLogger.SetLevel(log.InfoLevel)
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
	if coinsuper.GetInstance() == nil {
		return
	}
	if info, err := coinsuper.GetInstance().UserAssetInfo(nil); err == nil {
		meta = info.Assets["META"].Total
		eth = info.Assets["ETH"].Total
		btc = info.Assets["BTC"].Total
		logBalance("coinsuper", meta, eth, btc)
	}
	return
}

func getAbccBalnace() (meta, eth, btc float64) {
	if abcc.GetInstance() == nil {
		return
	}
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
func (c *Clients) GetBalances() {
	getKucoinBalnace(c.kucoin)
	getCoinsuperBalnace()
	getAbccBalnace()
}

func logTrade(pair, exchange, orderID, side, createdAt string, price, amount, fee, volume float64) {
	tradeLogger.WithFields(log.Fields{
		"pair":      pair,
		"exchange":  exchange,
		"orderID":   orderID,
		"side":      side,
		"price":     price,
		"amount":    amount,
		"fee":       fee,
		"volume":    volume,
		"createdAt": createdAt,
	}).Info("GatherTrades")
}

func getKucoinTrades(k *kucoin.Kucoin, pair string) {
	if ret, err := k.ListMergedDealtOrders(pair, "BUY", 20, 1, 0, 0); err == nil {
		for _, v := range ret.Datas {
			if !checkExistOrder(v.OrderOid) {
				logTrade(pair, "kucoin", v.OrderOid, "BUY", toDateStr(v.CreatedAt/1000), v.DealPrice, v.Amount, v.Fee, v.DealValue)
			}
		}
	}
	if ret, err := k.ListMergedDealtOrders(pair, "SELL", 20, 1, 0, 0); err == nil {
		for _, v := range ret.Datas {
			if !checkExistOrder(v.OrderOid) {
				logTrade(pair, "kucoin", v.OrderOid, "SELL", toDateStr(v.CreatedAt/1000), v.DealPrice, v.Amount, v.Fee, v.DealValue)
			}
		}
	}
}

func getAbccTrades(pair string) {
	if abcc.GetInstance() == nil {
		return
	}
	if ret, err := abcc.GetInstance().Trades(&abccTypes.Options{
		MarketCode: pair,
		PerPage:    "100",
	}); err == nil {
		for _, v := range ret.Trades {
			oID := strconv.FormatInt(v.ID, 10)
			if !checkExistOrder(oID) {
				price, err1 := strconv.ParseFloat(v.Price, 32)
				funds, err2 := strconv.ParseFloat(v.Funds, 32)
				fee, err3 := strconv.ParseFloat(v.Fee, 32)
				volume, err4 := strconv.ParseFloat(v.Volume, 32)
				if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
					logTrade(pair, "abcc", oID, v.Side, v.CreatedAt, price, volume, fee, funds)
				}
			}
		}
	}
}

// GetTrades records trades
func (c *Clients) GetTrades() {
	getKucoinTrades(c.kucoin, strings.ToUpper(targetSymbol)+"-ETH")
	getAbccTrades(strings.ToLower(targetSymbol) + "eth")
}
