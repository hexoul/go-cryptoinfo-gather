# go-cryptoinfo-gather
[![License](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hexoul/go-cryptoinfo-gather/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hexoul/go-cryptoinfo-gather)](https://goreportcard.com/report/github.com/hexoul/go-cryptoinfo-gather)
[![GoDoc](https://godoc.org/github.com/hexoul/go-cryptoinfo-gather?status.svg)](https://godoc.org/github.com/hexoul/go-cryptoinfo-gather)

> Crypto info gather

## Build
```bash
dep ensure
go build
```

## Usage
All options are not mandatory even if you do not use related APIs.

```bash
go run main.go \
-gitID=[GIT_ID] \
-gitPW=[GIT_PW] \
-logpath=[LOG_PATH] \
-targetSymbol=[TOKEN_SYMBOL] \
-targetAddr=[TOKEN_CONTRACT_ADDR] \
-targetQuotes=USD,BTC,ETH \
-targetSlugs=binance,okex \
-cmcApikey=[CMC_API_KEY] \
-coinsuper:accesskey=[COINSUPER_ACCESS_KEY] \
-coinsuper:secretkey=[COINSUPER_SECRET_KEY] \
-kucoin:accesskey=[KUCOIN_ACCESS_KEY] \
-kucoin:secretkey=[KUCOIN_SECRET_KEY]
```
...

## License
MIT