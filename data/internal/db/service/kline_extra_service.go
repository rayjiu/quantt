package service

import (
	"context"
	"sync"
	"time"

	"github.com/panjf2000/ants"
	"github.com/rayjiu/quantt/data/internal/db/model"
	"github.com/rayjiu/quantt/data/internal/db/repository"
	"github.com/sirupsen/logrus"
)

const batchSize = 1000 // 每批处理的记录数

var KlineExtraService klineExtraService

type klineExtraService struct {
	repo     *repository.KlineExtraRepository
	klineSvc *klineDayService
	lock     sync.Mutex
}

func NewKlineExtraService(repo *repository.KlineExtraRepository, klineSvc *klineDayService) *klineExtraService {
	return &klineExtraService{
		repo:     repo,
		klineSvc: klineSvc,
		lock:     sync.Mutex{},
	}
}

// CalculateAndSaveVolumeRatios 计算并保存量比数据
func (s *klineExtraService) CalculateAndSaveVolumeRatios(ctx context.Context, stockList []struct {
	StockCode  string
	MarketType int16
}) error {
	var batch []model.KlineExtra
	var pools, _ = ants.NewPool(10)
	var wg sync.WaitGroup
	for _, stock1 := range stockList {
		stock := stock1
		wg.Add(1)
		pools.Submit(func() {
			// 获取该股票的所有K线数据
			klines, err := s.klineSvc.GetKlineByCode(stock.StockCode, stock.MarketType)
			logrus.Infof("%v ->获取Kxian的长度为:%v", stock.StockCode, len(klines))
			if err != nil {
				wg.Done()
				panic(err)
			}

			// 计算每一天的量比
			for i := 0; i < len(klines); i++ {
				volumeRatio := calculateDailyVolumeRatio(klines, i)
				// logrus.Infof("%v -> %v 量比计算完成", klines[i].StockCode, klines[i].MarketDate)
				s.lock.Lock()
				batch = append(batch, model.KlineExtra{
					StockCode:   stock.StockCode,
					MarketType:  stock.MarketType,
					MarketDate:  klines[i].MarketDate,
					VolumeRatio: volumeRatio,
					UpdateTime:  time.Now(),
				})

				// 当批次达到指定大小时，执行批量保存
				if len(batch) >= batchSize {
					if err := s.repo.BatchUpsert(ctx, batch); err != nil {
						panic(err)
					}
					batch = batch[:0] // 清空批次
				}
				s.lock.Unlock()
			}
			wg.Done()
		})
	}
	wg.Wait()
	// 处理剩余的数据
	if len(batch) > 0 {
		if err := s.repo.BatchUpsert(ctx, batch); err != nil {
			return err
		}
	}

	return nil
}

// calculateDailyVolumeRatio 计算单个交易日的量比
func calculateDailyVolumeRatio(klines []model.KlineDay, currentIndex int) float64 {
	currentVolume := float64(klines[currentIndex].Volume) / 241

	// 获取前5个交易日的数据范围
	startIdx := max(0, currentIndex-5)
	var totalPreviousVolume float64
	var daysCount int

	// 计算前5个交易日的平均成交量
	for i := startIdx; i < currentIndex; i++ {
		totalPreviousVolume += float64(klines[i].Volume)
		daysCount++
	}

	// 如果没有历史数据，使用当天数据
	if daysCount == 0 {
		return 1.0
	}

	// 计算前N天平均每日成交量
	averagePreviousVolume := totalPreviousVolume / (float64(daysCount) * 241)

	// 计算量比：当日成交量 / 前N日平均成交量
	if averagePreviousVolume > 0 {
		return currentVolume / averagePreviousVolume
	}

	return 0
}

// max returns the larger of x or y
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
