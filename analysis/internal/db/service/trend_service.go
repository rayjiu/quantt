package service

import (
	"context"

	"github.com/rayjiu/quantt/analysis/internal/db/model"
	"github.com/rayjiu/quantt/analysis/internal/db/repository"
)

var TrenddService *trendService

type trendService struct {
	TrendRepository *repository.TrendRepository
}

// NewTrendService 初始化 TrendService
func NewTrendService(repo *repository.TrendRepository) *trendService {
	return &trendService{
		TrendRepository: repo,
	}
}

// GetTrendsByStock 通过 stockId、marketType 和 marketDate 获取 Trend 列表
func (s *trendService) GetTrendsByStock(ctx context.Context, stockId string, marketType uint32, marketDate uint32) ([]model.TrendItem, error) {
	return s.TrendRepository.GetTrendsByStock(ctx, stockId, marketType, marketDate)
}
