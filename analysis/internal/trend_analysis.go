package internal

import (
	"context"

	"github.com/rayjiu/quantt/analysis/internal/db/model"
	"github.com/rayjiu/quantt/analysis/internal/db/service"
	"github.com/rayjiu/quantt/analysis/internal/helper"
	log "github.com/sirupsen/logrus"
)

var trAnalysis trendAnalysis = trendAnalysis{}

type trendAnalysis struct {
	datas []model.TrendItem
}

func (analysis *trendAnalysis) loadTrendData(stockCode string, marketType uint32, marketDate uint32) {
	var trendDatas, err = service.TrenddService.GetTrendsByStock(context.Background(), stockCode, marketType, marketDate)
	if err != nil {
		log.Errorf("获取分时数据出错:%+v", err)
		return
	}
	analysis.datas = trendDatas
}

func (analysis *trendAnalysis) doAnalysis(stockCode string, marketType uint32, marketDate uint32) {
	analysis.loadTrendData(stockCode, marketType, marketDate)

	// 滑动窗口参数
	windowSize := 10 // 窗口大小：5分钟
	step := 1        // 滑动步长：1分钟
	helper.SlidingWindow(analysis.datas, windowSize, step, helper.AnalyzeVolumeStagnation, helper.AnalyzeVolumeIncrease)
}
