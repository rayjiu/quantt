package calculator

import (
	"errors"
	"sort"

	"github.com/rayjiu/quantt/analysis/internal/db/model"
)

// RollingRSICalculator 提供滚动窗口的RSI计算
type RollingRSICalculator struct {
	periods []int // 支持多个周期的RSI计算
}

// RollingRSIResult 存储每分钟的RSI结果
type RollingRSIResult struct {
	TrendItem    model.TrendItem // 对应的分时数据
	RSIResults   []RSIResult     // 每个周期的RSI值
	CalculateErr error           // 计算错误信息
}

// NewRollingRSICalculator 创建滚动RSI计算器
func NewRollingRSICalculator(periods ...int) *RollingRSICalculator {
	if len(periods) == 0 {
		periods = []int{14} // 默认14周期
	}
	return &RollingRSICalculator{
		periods: periods,
	}
}

// CalculateRollingRSI 对整个分时数据进行滚动RSI计算
func (rc *RollingRSICalculator) CalculateRollingRSI(trends []model.TrendItem, priceType string) ([]RollingRSIResult, error) {
	if len(trends) == 0 {
		return nil, errors.New("分时数据为空")
	}

	// 按交易时间排序（确保数据有序）
	sort.Slice(trends, func(i, j int) bool {
		return trends[i].TradeTime < trends[j].TradeTime
	})

	// 存储每个时间点的RSI结果
	rollingResults := make([]RollingRSIResult, 0, len(trends))

	// 创建RSI计算器
	rsiCalculator := NewRSICalculator(rc.periods...)

	// 滚动窗口计算RSI
	for i := range trends {
		// 取从开始到当前位置的切片
		currentWindowTrends := trends[:i+1]

		// 计算RSI
		rsiResults, err := rsiCalculator.CalculateRSI(currentWindowTrends, priceType)

		// 记录每个时间点的结果
		rollingResults = append(rollingResults, RollingRSIResult{
			TrendItem:    trends[i],
			RSIResults:   rsiResults,
			CalculateErr: err,
		})
	}

	return rollingResults, nil
}
