package crawler

func Start() {
	// service.SecotorService.RefreshCache()
	// sectorBaseInfoCrawler.startCrawSectorInfo()
	// ffCrawler.startHistoryFundFlowData()
	// kCrawler.startCrawlKlineData("515170", 1)
	// kCrawler.startCrawAllSecKliineData()
	kCrawler.startCrawKlineDataByBeiginData(20241126)
	// secCCrrawler.startCrawAllConsData()
}
