package spotclient

import "github.com/Gengxin-Zhang/exchangeclient/model"

type HuobiSpotClientImpl struct{}

func GetHuobiSpotClient() SpotClient {
	return &HuobiSpotClientImpl{}
}

func (h *HuobiSpotClientImpl) GetDepth(req *model.GetSpotDepthRequest) (*model.GetSpotCandlesResponse, error)
func (h *HuobiSpotClientImpl) GetTrades(req *model.GetSpotTradesRequest) (*model.GetSpotTradesResponse, error)
func (h *HuobiSpotClientImpl) GetCandles(req *model.GetSpotCandlesRequest) (*model.GetSpotCandlesResponse, error)
func (h *HuobiSpotClientImpl) GetTicker(req *model.GetSpotTickerRequest) (*model.GetSpotTickerResponse, error)

func (h *HuobiSpotClientImpl) Order() {}
