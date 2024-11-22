package calculator

import (
	"context"
	"fmt"

	"github.com/rayjiu/quantt/analysis/internal/db/model"
	"github.com/rayjiu/quantt/analysis/internal/db/service"
	"github.com/sirupsen/logrus"
)

type SupportResistanceCalc struct {
	datas []model.TrendItem
}

func NewSupportResistanceCalc(datas []model.TrendItem) *SupportResistanceCalc {
	return &SupportResistanceCalc{
		datas: datas,
	}
}

// calculateSupportAndResistance 根据数据计算支撑位和压力位
func (sc *SupportResistanceCalc) calculateSupportAndResistance(data []model.TrendItem) (float64, float64) {

	if len(data) == 0 {
		return 0, 0
	}

	// 初始化支撑位和压力位
	low := data[0].LastPrice
	high := data[0].LastPrice

	// 遍历数据找出最低点和最高点
	for _, item := range data {
		if item.LastPrice < low {
			low = item.LastPrice
		}
		if item.LastPrice > high {
			high = item.LastPrice
		}
	}

	return low, high
}

// ProcessTrendItems 实时轮询分时数据，每 5 分钟更新一次支撑位和压力位
func (sc *SupportResistanceCalc) ProcessTrendItems() {
	var trendItems []model.TrendItem = sc.datas
	var window []model.TrendItem // 保存最近 5 分钟的分时数据
	var support, resistance float64
	var config StrategyConfig = StrategyConfig{}
	var position *Position = &Position{}
	for _, item := range trendItems {
		// 将当前分时数据加入窗口
		window = append(window, item)

		// 如果 MarketTime 是 5 的倍数，则计算支撑位和压力位
		if item.MarketTime%5 == 0 {
			// 计算支撑位和压力位
			support, resistance = sc.calculateSupportAndResistance(window)
			// fmt.Printf("时间: %d, 支撑位: %.3f, 压力位: %.3f\n", item.MarketTime, support, resistance)
		}

		if support != 0 && resistance != 0 {
			config.SupportLevel = support
			config.ResistanceLevel = resistance
			config.VolumeMultiplier = 1.5
			config.TakeProfitPercentage = 0.43
			config.StopLossPercentage = 1
			sc.executeWaveStrategy(window, config, position)
		}
	}
}

// 策略配置
type StrategyConfig struct {
	SupportLevel         float64 // 支撑位
	ResistanceLevel      float64 // 压力位
	VolumeMultiplier     float64 // 成交量放大倍数阈值
	StopLossPercentage   float64 // 止损百分比
	TakeProfitPercentage float64 // 止盈百分比
}

// 当前持仓状态
type Position struct {
	IsHolding bool    // 是否持有仓位
	BuyPrice  float64 // 买入价格
	Volume    uint64  // 持有份额
}

// 波段操作策略实现
func (sc *SupportResistanceCalc) executeWaveStrategy(data []model.TrendItem, config StrategyConfig, position *Position) {
	if len(data) < 2 {
		fmt.Println("数据不足，无法执行策略")
		return
	}

	// 获取最近两条分时数据
	current := data[len(data)-1]
	previous := data[len(data)-2]
	fmt.Println(previous)
	var snap, err = service.SnapshotService.GetLatestSnapshot(context.Background(), "159792", 2, 20241121)
	if err != nil {
		logrus.Errorf("err:%v", err)
	}
	fmt.Printf("snapshot %+v \n", snap)
	var existedTradeVolume float64
	for _, d := range data {
		existedTradeVolume += float64(d.TradeVolume)
	}
	var currentReal = model.TrendItem{
		MarketTime:  snap.TradeTime,
		TradeVolume: snap.TradeAmount - uint64(existedTradeVolume),
		LastPrice:   snap.LastPrice,
	}
	fmt.Println("volume->", currentReal.TradeVolume)
	// var previouseReal = data[len(data)-1]
	// 成交量判断：当前成交量是否放大
	volumeIncrease := float64(current.TradeVolume) / float64(previous.TradeVolume)
	fmt.Printf("时间: %d, 最新价格: %.3f, 成交量变化倍数: %.2f\n",
		current.MarketTime, current.LastPrice, volumeIncrease)

	// 策略核心逻辑
	if !position.IsHolding {
		// 低吸条件：价格接近支撑位，且成交量放大
		if current.LastPrice <= config.SupportLevel && volumeIncrease > config.VolumeMultiplier {
			// 买入操作
			position.IsHolding = true
			position.BuyPrice = current.LastPrice
			position.Volume = 1000 // 假设每次操作1000份
			fmt.Printf("[*****买入信号] 时间: %d, 买入价格: %.3f, 持仓份额: %d\n",
				current.MarketTime, current.LastPrice, position.Volume)
		}
	} else {
		// 高抛条件：价格接近压力位，或达到止盈目标
		profitPercentage := (current.LastPrice - position.BuyPrice) / position.BuyPrice * 100
		if current.LastPrice >= config.ResistanceLevel || profitPercentage >= config.TakeProfitPercentage {
			// 卖出操作
			fmt.Printf("[*****卖出信号] 时间: %d, 卖出价格: %.3f, 持仓收益: %.2f%%\n",
				current.MarketTime, current.LastPrice, profitPercentage)
			position.IsHolding = false
			position.Volume = 0
		}

		// 止损条件
		lossPercentage := (position.BuyPrice - current.LastPrice) / position.BuyPrice * 100
		if lossPercentage >= config.StopLossPercentage {
			fmt.Printf("[*****止损信号] 时间: %d, 卖出价格: %.3f, 持仓亏损: %.2f%%\n",
				current.MarketTime, current.LastPrice, lossPercentage)
			position.IsHolding = false
			position.Volume = 0
		}
	}
}
