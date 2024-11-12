package crawler

import "github.com/rayjiu/quantt/data/internal/db/model"

type funflowCrawler struct {
	urlChain          chan fundFlowUrlInfo
	stopDataWriteChan chan int
	fundFlowChan      chan []*model.SecFundFlow
}
