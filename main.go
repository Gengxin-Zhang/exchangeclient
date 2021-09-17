package main

import (
	"github.com/Gengxin-Zhang/exchangeclient/clients/ws/public"
	"github.com/Gengxin-Zhang/exchangeclient/pkg/constant"
	"github.com/sirupsen/logrus"
)

func main() {
	// 双链接
	huobi := public.GetPublicWSClient(constant.HUOBI, public.WithChannel(true)).(*public.HuobiPublicClient)
	huobi2 := public.GetPublicWSClient(constant.HUOBI, public.WithChannel(true)).(*public.HuobiPublicClient)
	huobi.Start()
	huobi2.Start()
	symbols := []string{"btcusdt", "ethusdt", "trxusdt", "eosusdt", "filusd", "filusdt", "xrpusdt", "bchusdt", "ltcusdt", "htusdt", "ethbtc", "trxbtc", "htbtc", "eosbtc", "hteth", "usdthusd", "ltcbtc", "bchbtc", "xrpbtc", "topbtc", "eoseth", "trxeth", "etcbtc", "topht", "hptht", "eosht", "bchht", "xtzusdt", "okbusdt", "adausdt", "linkusdt", "xmrusdt", "xlmusdt", "leousdt", "etcusdt", "dashusdt", "neousdt", "crousdt", "atomusdt", "iotausdt", "zecusdt", "xemusdt", "ontusdt", "vetusdt", "batusdt", "dogeusdt", "algousdt", "qtumusdt", "dcrusdt"}
	for _, symbol := range symbols {
		huobi.SubCandle(symbol, constant.MIN1)
		huobi.SubDepth(symbol)
		huobi.SubTicker(symbol)
		huobi.SubTrade(symbol)
		huobi2.SubCandle(symbol, constant.MIN1)
		huobi2.SubDepth(symbol)
		huobi2.SubTicker(symbol)
		huobi2.SubTrade(symbol)
	}
	// TODO：双链接去重
	for {
		select {
		case message := <-huobi.DepthChannel:
			logrus.Infof("depth: %v", message)
		case message := <-huobi.CandleChannel:
			logrus.Infof("candle: %v", message)
		case message := <-huobi.TradeChannel:
			logrus.Infof("trade: %v", message)
		case message := <-huobi.TickerChannel:
			logrus.Infof("ticker: %v", message)
		case message := <-huobi2.DepthChannel:
			logrus.Infof("depth: %v", message)
		case message := <-huobi2.CandleChannel:
			logrus.Infof("candle: %v", message)
		case message := <-huobi2.TradeChannel:
			logrus.Infof("trade: %v", message)
		case message := <-huobi2.TickerChannel:
			logrus.Infof("ticker: %v", message)
		}
	}
}
