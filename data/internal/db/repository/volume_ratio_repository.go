// repository/volume_ratio_repository.go
package repository

import (
	"context"

	"github.com/rayjiu/quantt/data/internal/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KlineExtraRepository struct {
	db *gorm.DB
}

func NewKlineExtraRepository(db *gorm.DB) *KlineExtraRepository {
	return &KlineExtraRepository{db: db}
}

// BatchUpsert 批量更新或插入量比数据
func (r *KlineExtraRepository) BatchUpsert(ctx context.Context, ratios []model.KlineExtra) error {
	if len(ratios) == 0 {
		return nil
	}

	err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "stock_code"},
			{Name: "market_type"},
			{Name: "market_date"},
		}, // 定义冲突字段
		DoUpdates: clause.AssignmentColumns([]string{
			"volume_ratio",
			"update_time",
		}), // 冲突时更新的字段
	}).Create(&ratios).Error

	return err
}

// GetByStockCode 获取指定股票的所有量比数据
func (r *KlineExtraRepository) GetByStockCode(ctx context.Context, stockCode string, marketType int16) ([]model.KlineExtra, error) {
	var ratios []model.KlineExtra
	err := r.db.Where("stock_code = ? AND market_type = ?", stockCode, marketType).
		Order("market_date ASC").
		Find(&ratios).Error
	return ratios, err
}
