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
	upbit "github.com/hexoul/go-upbit"
	upbitTypes "github.com/hexoul/go-upbit/types"
	bittrex "github.com/toorop/go-bittrex"

	log "github.com/sirupsen/logrus"
)

const (
	existLen = 10000
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

func getKucoinBalnace() (meta, eth, btc float64) {
	if kucoin.GetInstance() == nil {
		return
	}
	if accounts, err := kucoin.GetInstance().ListAccounts(nil); err == nil {
		for _, v := range accounts {
			bal, pErr := strconv.ParseFloat(v.Balance, 32)
			if pErr != nil {
				continue
			}
			switch v.Currency {
			case "META":
				meta += bal
				break
			case "ETH":
				eth += bal
				break
			case "BTC":
				btc += bal
				break
			}
		}
		logBalance("kucoin", meta, eth, btc)
	}
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

func getBittrexBalnace(b *bittrex.Bittrex) (meta, eth, btc float64) {
	if balances, err := b.GetBalances(); err == nil {
		for _, v := range balances {
			if v.Currency == "META" {
				meta, _ = v.Balance.Float64()
			} else if v.Currency == "ETH" {
				eth, _ = v.Balance.Float64()
			} else if v.Currency == "BTC" {
				btc, _ = v.Balance.Float64()
			}
		}
		logBalance("bittrex", meta, eth, btc)
	}
	return
}

func getUpbitBalnace() (meta, eth, btc float64) {
	if upbit.GetInstance() == nil {
		return
	}
	if ret, err := upbit.GetInstance().Accounts(); err == nil {
		for _, v := range ret {
			if v.Currency == "META" {
				meta = sumStrFloat(v.Balance, v.Locked)
			} else if v.Currency == "ETH" {
				eth = sumStrFloat(v.Balance, v.Locked)
			} else if v.Currency == "BTC" {
				btc = sumStrFloat(v.Balance, v.Locked)
			}
		}
		logBalance("upbit", meta, eth, btc)
	}
	return
}

// GetBalances records balances
func (c *Clients) GetBalances() {
	getKucoinBalnace()
	getCoinsuperBalnace()
	getAbccBalnace()
	getBittrexBalnace(c.bittrex)
	getUpbitBalnace()
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

func getKucoinTrades(pair string) {
	if kucoin.GetInstance() == nil {
		return
	}
	if ret, err := kucoin.GetInstance().ListFills(nil); err == nil {
		for _, v := range ret {
			if v.Symbol == pair && !checkExistOrder(v.TradeID) {
				price, pErr1 := strconv.ParseFloat(v.Price, 32)
				size, pErr2 := strconv.ParseFloat(v.Size, 32)
				fee, pErr3 := strconv.ParseFloat(v.Fee, 32)
				funds, pErr4 := strconv.ParseFloat(v.Funds, 32)
				if pErr1 == nil && pErr2 == nil && pErr3 == nil && pErr4 == nil {
					logTrade(pair, "kucoin", v.TradeID, v.Side, toDateStr(v.CreatedAt/1000), price, size, fee, funds)
				}
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

func getBittrexTrades(b *bittrex.Bittrex, pair string) {
	if b == nil {
		return
	}
	// if orders, err := b.GetOrderHistory("BTC-META"); err == nil {
	// 	for _, v := range orders {

	// 	}
	// }
}

func getUpbitTrades(pair string) {
	if upbit.GetInstance() == nil {
		return
	}
	if ret, err := upbit.GetInstance().GetOrders(&upbitTypes.Options{
		Market: pair,
		State:  upbitTypes.StateOptions.Done,
	}); err == nil {
		for _, v := range ret {
			if !checkExistOrder(v.UUID) {
				price, err1 := strconv.ParseFloat(v.Price, 32)
				fee, err2 := strconv.ParseFloat(v.PaidFee, 32)
				volume, err3 := strconv.ParseFloat(v.ExecutedVolume, 32)
				if err1 == nil && err2 == nil && err3 == nil {
					logTrade(pair, "upbit", v.UUID, v.Side, v.CreatedAt, price, volume, fee, price*volume)
				}
			}
		}
	}
}

// GetTrades records trades
func (c *Clients) GetTrades() {
	getKucoinTrades(strings.ToUpper(targetSymbol) + "-ETH")
	getAbccTrades(strings.ToLower(targetSymbol) + "eth")
	btcMeta := "BTC-" + strings.ToUpper(targetSymbol)
	getBittrexTrades(c.bittrex, btcMeta)
	getUpbitTrades(btcMeta)
}
