package internal

import (
	"context"
	"fmt"

	"github.com/rayjiu/quantt/analysis/internal/calculator"
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

func (analysis *trendAnalysis) doRSIAnalysis(stockCode string, marketType uint32, marketDate uint32) {
	analysis.loadTrendData(stockCode, marketType, marketDate)

	fmt.Printf("1-> %+v \n", analysis.datas[0])
	fmt.Printf("2-> %+v \n", analysis.datas[1])
	fmt.Printf("3-> %+v \n", analysis.datas[2])
	fmt.Printf("4-> %+v \n\n", analysis.datas[3])

	rollingCalculator := calculator.NewRollingRSICalculator(5, 7, 9)

	rollingResults, err := rollingCalculator.CalculateRollingRSI(analysis.datas, "close")
	if err != nil {
		log.Errorf("计算出现错误:%v", err)
	}
	// 遍历结果
	for _, result := range rollingResults {
		// 获取每个时间点的RSI值
		for _, rsiResult := range result.RSIResults {
			fmt.Printf("时间: %v, 周期: %d, RSI: %.2f\n",
				result.TrendItem.MarketTime,
				rsiResult.Period,
				rsiResult.RSI)
		}
	}
}

func (analysis *trendAnalysis) doVolumePriceAnalysis(stockCode string, marketType uint32, marketDate uint32) {
	analysis.loadTrendData(stockCode, marketType, marketDate)

	for i := 10; i < len(analysis.datas); i++ {
		var targetDatas = analysis.datas[:i]
		var volumes []uint64
		var weights []float64
		for i, data := range targetDatas {
			volumes = append(volumes, uint64(data.TradeVolume))
			weights = append(weights, float64(i+1))
		}
		var calculator *calculator.VolumePriceCalculator = calculator.NewPriceCalculator(volumes, weights)
		var volumeAverage = calculator.CalculateSimpleMovingAverage()
		volumeIncrease := float64(analysis.datas[i].TradeVolume) / volumeAverage

		var latest = analysis.datas[i]
		var previous = analysis.datas[i-1]
		var multiplier float64 = 1.5
		if volumeIncrease > multiplier {
			// 价格变化趋势
			if latest.LastPrice < previous.LastPrice {
				fmt.Printf("%d [%f] -> Buy Signal: Stock %s, Last Price: %.4f\n", latest.MarketTime, volumeIncrease, latest.StockId, latest.LastPrice)
			} else if latest.LastPrice > previous.LastPrice {
				fmt.Printf("%d [%f] -> Sell Signal: Stock %s, Last Price: %.4f\n", latest.MarketTime, volumeIncrease, latest.StockId, latest.LastPrice)
			} else {
				fmt.Println("No significant price movement despite volume increase.")
			}
		} else {
			fmt.Println("No significant volume increase.")
		}
	}
}

func (analysis *trendAnalysis) doSupportAndResistance(stockCode string, marketType uint32, marketDate uint32) {
	analysis.loadTrendData(stockCode, marketType, marketDate)
	var calc = calculator.NewSupportResistanceCalc(analysis.datas)
	log.Infof("data加载完成, datas:%v", len(analysis.datas))
	calc.ProcessTrendItems()
}
