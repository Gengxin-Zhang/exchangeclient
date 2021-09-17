package spotclient

import (
	"github.com/Gengxin-Zhang/exchangeclient/model"
)

type SpotClient interface {
	// data
	GetDepth(req *model.GetSpotDepthRequest) (*model.GetSpotCandlesResponse, error)
	GetTrades(req *model.GetSpotTradesRequest) (*model.GetSpotTradesResponse, error)
	GetTicker(req *model.GetSpotTickerRequest) (*model.GetSpotTickerResponse, error)
	GetCandles(req *model.GetSpotCandlesRequest) (*model.GetSpotCandlesResponse, error)
}

type ClientImpls struct {
	Huobi SpotClient
}

var clientImps *ClientImpls

func GetSpotClient() *ClientImpls {
	if clientImps == nil {
		clientImps = &ClientImpls{
			Huobi: GetHuobiSpotClient(),
		}
	}
	return clientImps
}
