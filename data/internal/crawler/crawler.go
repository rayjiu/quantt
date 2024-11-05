package crawler

import "github.com/rayjiu/quantt/data/internal/db/service"

func Start() {
	service.SecotorService.RefreshCache()
	eastmoney.startCrawSectorInfo()
}
