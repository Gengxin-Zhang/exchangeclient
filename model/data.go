package model

import (
	"time"

	"github.com/Gengxin-Zhang/exchangeclient/pkg/constant"
	"github.com/shopspring/decimal"
)

type Depth struct {
	Symbol string
	Ts     time.Time
	Bids   [][]decimal.Decimal
	Asks   [][]decimal.Decimal
}

type Trade struct {
	Symbol string
	Side   constant.Side
	Volume decimal.Decimal
	Price  decimal.Decimal
	Tid    string
	Oid    string
	Ts     time.Time
}

type Ticker struct {
	Symbol    string
	Ts        time.Time
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Amount    decimal.Decimal
	Volume    decimal.Decimal
	Count     decimal.Decimal
	Bid       decimal.Decimal
	BidSize   decimal.Decimal
	Ask       decimal.Decimal
	AskSize   decimal.Decimal
	LastPrice decimal.Decimal
	LastSize  decimal.Decimal
}

/*
class Candle(NamedTuple):
    period: CandlePeriod
    ts: Timestamp
    open_price: Decimal
    high_price: Decimal
    low_price: Decimal
    close_price: Decimal
    base_volume: Decimal
    quote_volume: Decimal = Decimal(math.nan)
    no_trades: Optional[int] = None
*/

type Candle struct {
	Symbol      string
	Period      constant.Period
	Ts          time.Time
	Open        decimal.Decimal
	Close       decimal.Decimal
	High        decimal.Decimal
	Low         decimal.Decimal
	BaseVolume  decimal.Decimal
	QuoteVolume decimal.Decimal
	Trads       int32
}

/*
	symbol: str
    status: OrderEvent
    side: Side
    price: Decimal = Decimal(math.nan)
    volume: Decimal = Decimal(math.nan)
    filled: Decimal = Decimal(math.nan)
    cid: Optional[str] = None
    oid: Optional[str] = None
    create_ts: Optional[Timestamp] = None
    finish_ts: Optional[Timestamp] = None
    cancel_ts: Optional[Timestamp] = None
    update_ts: Optional[Timestamp] = None
*/

type OrderStatus int32

const (
	OrderOpen OrderStatus = iota
)

type Order struct {
	Symbol string
	Status OrderStatus
}
