package calculator

import (
	"context"

	"github.com/rayjiu/quantt/data/internal/db/service"
	log "github.com/sirupsen/logrus"
)

func StartCalHistoryVolumeRate() {
	service.SecotorService.RefreshCache()
	var sectorMap = service.SecotorService.GetCachedSector()

	var stockCodeList []struct {
		StockCode  string
		MarketType int16
	}
	for _, v := range sectorMap {

		stockCodeList = append(stockCodeList,
			struct {
				StockCode  string
				MarketType int16
			}{
				StockCode:  v.SecCode,
				MarketType: 90,
			},
		)
	}

	if err := service.KlineExtraService.CalculateAndSaveVolumeRatios(context.Background(), stockCodeList); err != nil {
		log.Errorf("出现错误:%v", err)
	}
	log.Infof("量比计算完成。")
}
