package calculator

import (
	"errors"
	"math"
	"sort"
	"time"

	"github.com/rayjiu/quantt/analysis/internal/db/model"
)

// RSIResult 存储RSI计算结果
type RSIResult struct {
	RSI           float64   // RSI指标值
	Period        int       // 计算周期
	CalculateTime time.Time // 计算时间
}

// RSICalculator 提供RSI指标计算的结构体
type RSICalculator struct {
	periods []int // 支持多个周期的RSI计算
}

// NewRSICalculator 创建RSI计算器
func NewRSICalculator(periods ...int) *RSICalculator {
	if len(periods) == 0 {
		periods = []int{14} // 默认14周期
	}
	return &RSICalculator{
		periods: periods,
	}
}

// calculateRSI 核心RSI计算逻辑
func (rc *RSICalculator) calculateRSI(changes []float64, period int) (float64, error) {
	if len(changes) < period {
		return 0, errors.New("价格序列长度不足以计算RSI")
	}

	// // 计算价格变化
	// changes := make([]float64, len(prices)-1)
	// for i := 1; i < len(prices); i++ {
	// 	changes[i-1] = prices[i] - prices[i-1]
	// }

	// 分离上涨和下跌
	ups := make([]float64, 0)
	downs := make([]float64, 0)

	for _, change := range changes {
		if change > 0 {
			ups = append(ups, change)
		} else if change < 0 {
			downs = append(downs, math.Abs(change))
		}
	}

	// 计算平均上涨和下跌
	avgUp := calculateAverage(ups, period)
	avgDown := calculateAverage(downs, period)

	// 避免除零
	if avgDown == 0 {
		return 100, nil
	}

	// RSI计算
	rs := avgUp / avgDown
	rsi := 100 - (100 / (1 + rs))

	return rsi, nil
}

// calculateAverage 计算移动平均
func calculateAverage(values []float64, period int) float64 {
	if len(values) == 0 {
		return 0
	}

	// 如果数据少于周期，使用全部数据
	if len(values) < period {
		period = len(values)
	}

	// 取最近的数据
	recentValues := values[len(values)-period:]
	sum := 0.0
	for _, val := range recentValues {
		sum += val
	}
	return sum / float64(period)
}

// CalculateRSI 计算给定分时数据的RSI
func (rc *RSICalculator) CalculateRSI(trends []model.TrendItem, priceType string) ([]RSIResult, error) {
	if len(trends) == 0 {
		return nil, errors.New("分时数据为空")
	}

	// 根据交易时间排序
	sort.Slice(trends, func(i, j int) bool {
		return trends[i].TradeTime < trends[j].TradeTime
	})

	// 提取价格序列
	var prices []float64
	switch priceType {
	case "close":
		prices = extractPrices(trends, func(t model.TrendItem) float64 { return t.Close - t.Open })
	case "last":
		prices = extractPrices(trends, func(t model.TrendItem) float64 { return t.LastPrice - t.Open })
	case "avg":
		prices = extractPrices(trends, func(t model.TrendItem) float64 { return t.AvgPrice - t.Open })
	default:
		return nil, errors.New("无效的价格类型")
	}

	// 计算多个周期的RSI
	results := make([]RSIResult, 0, len(rc.periods))
	for _, period := range rc.periods {
		rsi, err := rc.calculateRSI(prices, period)
		if err != nil {
			// return nil, err
			continue
		}

		results = append(results, RSIResult{
			RSI:           rsi,
			Period:        period,
			CalculateTime: time.Now(),
		})
	}

	return results, nil
}

// extractPrices 从分时数据提取特定价格
func extractPrices(trends []model.TrendItem, priceExtractor func(model.TrendItem) float64) []float64 {
	prices := make([]float64, len(trends))
	for i, trend := range trends {
		prices[i] = priceExtractor(trend)
	}
	return prices
}
