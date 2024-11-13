package crawler

import "github.com/rayjiu/quantt/data/internal/db/service"

func Start() {
	service.SecotorService.RefreshCache()
	sectorBaseInfoCrawler.startCrawSectorInfo()
	ffCrawler.startHistoryFundFlowData()
	kCrawler.startCrawlKlineData("601318", 1)
}
