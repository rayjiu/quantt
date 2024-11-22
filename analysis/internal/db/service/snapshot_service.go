package service

import (
	"context"
	"fmt"

	"github.com/rayjiu/quantt/analysis/internal/db/model"
	"github.com/rayjiu/quantt/analysis/internal/db/repository"
)

var SnapshotService *snapshotService

type snapshotService struct {
	SnapshotRepository *repository.SnapshotRepository
}

// NewTrendService 初始化 TrendService
func NewSnapshotService(repo *repository.SnapshotRepository) *snapshotService {
	return &snapshotService{
		SnapshotRepository: repo,
	}
}

// GetTrendsByStock 通过 stockId、marketType 和 marketDate 获取 Trend 列表
func (s *snapshotService) GetLatestSnapshot(ctx context.Context, stockId string, marketType uint32, marketDate uint32) (*model.StockSnapshot, error) {
	var data, err = s.SnapshotRepository.GetLatestSnapshot(ctx, stockId, marketType, marketDate)
	if err != nil {
		fmt.Println(err)
	}

	return data, nil
}
