package helper

import (
	"fmt"

	"github.com/rayjiu/quantt/analysis/internal/db/model"
)

// 滑动窗口处理
func SlidingWindow(klines []model.TrendItem, windowSize int, step int, analyzeFuncs ...func([]model.TrendItem, float64)) {
	total := len(klines)
	var openPrice float64
	for i := 0; i+windowSize <= total; i += step {
		window := klines[i : i+windowSize]
		if i == 0 {
			openPrice = window[0].Open
		}
		for _, aFunc := range analyzeFuncs {
			aFunc(window, openPrice)
		}
	}
}

// 分析函数：判断放量滞涨
func AnalyzeVolumeStagnation(window []model.TrendItem, dayOpenPrice float64) {
	// 窗口中的当前K线
	current := window[len(window)-1]

	// 计算窗口内均量
	var totalVolume float64
	var upCnt, downCnt int
	for i, kline := range window[:len(window)-1] {
		totalVolume += float64(kline.TradeVolume)
		if i > 0 {
			if kline.Close >= window[i-1].Close {
				upCnt += 1
			} else {
				downCnt += 1
			}
		}
	}
	averageVolume := totalVolume / float64(len(window)-1)

	// 检查放量条件（当前成交量 > 均量的2倍）
	isHighVolume := float64(current.TradeVolume) > averageVolume*2

	// 检查滞涨条件（涨幅低于0.1%）
	priceChange := (current.Close - current.Open) / current.Open
	// priceChange := (current.Close - dayOpenPrice) / dayOpenPrice
	// priceChange := (current.Close - window[0].Close) / window[0].Close
	isStagnant := priceChange <= 0.001
	// isTrendup := upCnt > downCnt

	// 输出分析结果
	// fmt.Printf("Timestamp: %d | High Volume: %t | Stagnant: %t\n", current.MarketTime, isHighVolume, isStagnant)
	if isHighVolume && isStagnant {
		fmt.Printf("Detected Volume Stagnation! at:%d   \n", current.MarketTime)
		// if isTrendup {
		// 	fmt.Print("Sell it. \n")
		// } else {
		// 	fmt.Print("Buy it. \n")
		// }
	}
}

// 判断放量上涨
func AnalyzeVolumeIncrease(window []model.TrendItem, dayOpenPrice float64) {
	// 窗口中的当前K线
	current := window[len(window)-1]

	// 计算窗口内均量（不包含当前K线）
	var totalVolume float64
	for _, kline := range window[:len(window)-1] {
		totalVolume += float64(kline.TradeVolume)
	}
	averageVolume := totalVolume / float64(len(window)-1)

	// 检查放量条件（当前成交量 > 均量的2倍）
	isHighVolume := float64(current.TradeVolume) > averageVolume*2

	// 检查涨幅条件（涨幅 > 0.5%）
	// priceChange := (current.Close - current.Open) / current.Open
	priceChange := (current.Close - window[0].Close) / window[0].Close
	isSignificantPriceIncrease := priceChange > 0.005

	// 输出分析结果
	// fmt.Printf("Timestamp: %d | High Volume: %t | Significant Price Increase: %t\n", current.MarketTime, isHighVolume, isSignificantPriceIncrease)
	if isHighVolume && isSignificantPriceIncrease {
		fmt.Printf("Detected Volume Increase with Significant Price Movement! time:%d, priceChange:%v \n", current.MarketTime, priceChange)
	}
}
