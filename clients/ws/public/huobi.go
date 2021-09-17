package public

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/Gengxin-Zhang/exchangeclient/model"
	"github.com/Gengxin-Zhang/exchangeclient/pkg/constant"
	"github.com/Gengxin-Zhang/exchangeclient/pkg/wsclient"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type HuobiPublicClient struct {
	*wsclient.WSClient
	opt *Options

	Started bool

	DepthChannel  chan *model.Depth
	TradeChannel  chan *model.Trade
	CandleChannel chan *model.Candle
	TickerChannel chan *model.Ticker
}

const HUOBI_HOST = "wss://api.huobi.pro/ws"

func GetHuobiPublicClient(o *Options) PublicClient {
	huobiClient := &HuobiPublicClient{opt: o}
	return huobiClient
}

func (h *HuobiPublicClient) Start() {
	opt := wsclient.NewOptions(wsclient.WithHost(HUOBI_HOST), wsclient.WithInflate(h.inflact))
	h.WSClient = &wsclient.WSClient{
		Opt:            opt,
		MessageHandler: h.onMessage,
	}
	h.DepthChannel = make(chan *model.Depth, 50)
	h.CandleChannel = make(chan *model.Candle, 50)
	h.TradeChannel = make(chan *model.Trade, 50)
	h.TickerChannel = make(chan *model.Ticker, 50)
	h.Channels = make(map[string]bool)
	h.WSClient.Start()
	h.Started = true
}

func (h *HuobiPublicClient) inflact(input []byte) ([]byte, error) {
	buf := bytes.NewBuffer(input)
	reader, gzipErr := gzip.NewReader(buf)
	if gzipErr != nil {
		return nil, gzipErr
	}
	defer reader.Close()

	result, readErr := ioutil.ReadAll(reader)
	if readErr != nil {
		return nil, readErr
	}

	return result, nil
}

func (h *HuobiPublicClient) Subscribe(channels []string) {
	for _, channel := range channels {
		if h.Started {
			h.Channels[channel] = true
			h.Send([]byte(channel))
		}
	}
}

func (h *HuobiPublicClient) SubCandle(symbol string, period constant.Period) {
	channel := fmt.Sprintf("{\"sub\": \"market.%v.kline.%v\"}", symbol, constant.HuobiPeriodMap[period])
	if h.Started {
		h.Subscribe([]string{channel})
	}
}

func (h *HuobiPublicClient) SubDepth(symbol string) {
	channel := fmt.Sprintf("{\"sub\": \"market.%v.depth.%v\"}", symbol, "step0")
	if h.Started {
		h.Subscribe([]string{channel})
	}
}

func (h *HuobiPublicClient) SubTicker(symbol string) {
	channel := fmt.Sprintf("{\"sub\": \"market.%v.ticker\"}", symbol)
	if h.Started {
		h.Subscribe([]string{channel})
	}
}

func (h *HuobiPublicClient) SubTrade(symbol string) {
	channel := fmt.Sprintf("{\"sub\": \"market.%v.trade.detail\"}", symbol)
	if h.Started {
		h.Subscribe([]string{channel})
	}
}

func (h *HuobiPublicClient) onMessage(message []byte) {
	var data map[string]interface{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		logrus.Errorf("unmarshal message failed, message: %v", string(message))
		return
	}
	if ts, exist := data["ping"]; exist {
		ts = int64(ts.(float64))
		pongData := fmt.Sprintf("{\"pong\":%d }", ts)
		h.Send([]byte(pongData))
	} else if topic, exist := data["subbed"]; exist {
		logrus.Infof("\"subbed\": \"%v\"", topic)
	} else if topic, exist := data["unsubbed"]; exist {
		logrus.Info("\"unsubbed\": \"%s\"", topic)
	} else if topic, exist := data["ch"]; exist {
		elements := strings.Split(topic.(string), ".")
		if len(elements) < 3 {
			return
		}
		symbol := elements[1]
		tick := data["tick"].(map[string]interface{})
		switch elements[2] {
		case "ticker":
			h.tickerHandler(symbol, message, tick)
		case "depth":
			h.depthHandler(symbol, message, tick)
		case "kline":
			h.candleHandler(symbol, message, tick)
		case "trade":
			h.tradeHandler(symbol, message, tick)
		default:
			logrus.Warnf("unknow topic type, topic: %v", topic)
			return
		}
	} else if topic, exist := data["rep"]; exist {
		logrus.Infof("recive req data, topic: %v, data: %v", topic, data)
	} else if code, exist := data["err-code"]; exist {
		msg := data["err-msg"]
		logrus.Error("%d:%s", code, msg)
	} else {
		logrus.Info("WebSocket received unknow data: %s", data)
	}
}

func (h *HuobiPublicClient) tickerHandler(symbol string, raw []byte, data map[string]interface{}) {
	type TickerMessage struct {
		Ts   int64 `json:"ts"`
		Tick struct {
			Open      decimal.Decimal `json:"open"`
			High      decimal.Decimal `json:"high"`
			Low       decimal.Decimal `json:"low"`
			Close     decimal.Decimal `json:"close"`
			Amount    decimal.Decimal `json:"amount"`
			Volume    decimal.Decimal `json:"vol"`
			Count     decimal.Decimal `json:"count"`
			Bid       decimal.Decimal `json:"bid"`
			BidSize   decimal.Decimal `json:"bidSize"`
			Ask       decimal.Decimal `json:"ask"`
			AskSize   decimal.Decimal `json:"askSize"`
			LastPrice decimal.Decimal `json:"lastPrice"`
			LastSize  decimal.Decimal `json:"lastSize"`
		} `json:"tick"`
	}
	var tickerMessage TickerMessage
	err := json.Unmarshal(raw, &tickerMessage)
	if err != nil {
		logrus.WithError(err).Error("json unmarshal ticker message failed")
		return
	}
	ticker := &model.Ticker{
		Symbol:    symbol,
		Ts:        time.Unix(tickerMessage.Ts/1000, tickerMessage.Ts%1000),
		Open:      tickerMessage.Tick.Open,
		Close:     tickerMessage.Tick.Close,
		High:      tickerMessage.Tick.High,
		Low:       tickerMessage.Tick.Low,
		Amount:    tickerMessage.Tick.Amount,
		Volume:    tickerMessage.Tick.Volume,
		Count:     tickerMessage.Tick.Count,
		Bid:       tickerMessage.Tick.Bid,
		BidSize:   tickerMessage.Tick.BidSize,
		Ask:       tickerMessage.Tick.Ask,
		AskSize:   tickerMessage.Tick.AskSize,
		LastPrice: tickerMessage.Tick.LastPrice,
		LastSize:  tickerMessage.Tick.LastSize,
	}
	if h.opt.UseChannel {
		h.TickerChannel <- ticker
	}
}

func (h *HuobiPublicClient) depthHandler(symbol string, raw []byte, data map[string]interface{}) {
	type DepthMessage struct {
		Tick struct {
			Asks [][]decimal.Decimal `json:"asks"`
			Bids [][]decimal.Decimal `json:"bids"`
			Ts   int64               `json:"ts"`
		} `json:"tick"`
	}

	var depthMessage DepthMessage
	err := json.Unmarshal(raw, &depthMessage)
	if err != nil {
		logrus.WithError(err).Errorf("unmarshal depth message failed")
	}
	depth := &model.Depth{
		Symbol: symbol,
		Ts:     time.Unix(depthMessage.Tick.Ts/1000, depthMessage.Tick.Ts%1000),
		Asks:   depthMessage.Tick.Asks,
		Bids:   depthMessage.Tick.Bids,
	}

	if h.opt.UseChannel {
		h.DepthChannel <- depth
	}
}

func (h *HuobiPublicClient) candleHandler(symbol string, raw []byte, data map[string]interface{}) {
	type CandleMessage struct {
		Tick struct {
			Id     int64           `json:"id"`
			Open   decimal.Decimal `json:"open"`
			Close  decimal.Decimal `json:"close"`
			High   decimal.Decimal `json:"high"`
			Low    decimal.Decimal `json:"low"`
			Amount decimal.Decimal `json:"amount"`
			Volume decimal.Decimal `json:"vol"`
			Count  int32           `json:"count"`
		} `json:"tick"`
	}
	var candleMessage CandleMessage
	err := json.Unmarshal(raw, &candleMessage)
	if err != nil {
		logrus.WithError(err).Error("json unmarshal candle message failed")
		return
	}
	candle := &model.Candle{
		Symbol:      symbol,
		Ts:          time.Unix(candleMessage.Tick.Id, 0),
		Trads:       candleMessage.Tick.Count,
		Open:        candleMessage.Tick.Open,
		Close:       candleMessage.Tick.Close,
		High:        candleMessage.Tick.High,
		Low:         candleMessage.Tick.Low,
		BaseVolume:  candleMessage.Tick.Amount,
		QuoteVolume: candleMessage.Tick.Volume,
	}
	if h.opt.UseChannel {
		h.CandleChannel <- candle
	}
}

func (h *HuobiPublicClient) tradeHandler(symbol string, raw []byte, data map[string]interface{}) {
	type TradeMessage struct {
		Tick struct {
			Data []struct {
				Ts      int64   `json:"ts"`
				TradeId int64   `json:"tradeId"`
				Amount  float64 `json:"amount"`
				Price   float64 `json:"price"`
				Side    string  `json:"direction"`
			} `json:"data"`
		} `json:"tick"`
	}
	var tradeMessage TradeMessage
	err := json.Unmarshal(raw, &tradeMessage)
	if err != nil {
		logrus.WithError(err).Errorf("unmarshal trade message failed")
		return
	}
	for _, tradeData := range tradeMessage.Tick.Data {
		trade := &model.Trade{
			Symbol: symbol,
			Ts:     time.Unix(int64(tradeData.Ts/1000), int64(tradeData.Ts%1000)),
			Price:  decimal.NewFromFloat(tradeData.Price),
			Volume: decimal.NewFromFloat(tradeData.Amount),
			Tid:    strconv.FormatInt(tradeData.TradeId, 10),
		}
		if tradeData.Side == "sell" {
			trade.Side = constant.SELL
		} else {
			trade.Side = constant.BUY
		}

		if h.opt.UseChannel {
			h.TradeChannel <- trade
		}
	}
}
