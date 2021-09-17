package public

import "github.com/Gengxin-Zhang/exchangeclient/pkg/constant"

type PublicClient interface {
	Start()
	Subscribe(channels []string)
	SubCandle(symbol string, period constant.Period)
	SubDepth(symbol string)
	SubTrade(symbol string)
	SubTicker(symbol string)
}

func GetPublicWSClient(exchange constant.ExchangeID, opts ...Opt) PublicClient {
	o := &Options{}
	o.Use(opts...)
	switch exchange {
	case constant.HUOBI:
		return GetHuobiPublicClient(o)
	default:
		return nil
	}
}
