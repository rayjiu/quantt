package repository

import (
	"errors"

	"github.com/rayjiu/quantt/data/internal/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KlineDayRepository struct {
	db *gorm.DB
}

// NewKlineDayRepository creates a new KlineDayRepository instance
func NewKlineDayRepository(db *gorm.DB) *KlineDayRepository {
	return &KlineDayRepository{db: db}
}

// Create inserts a new KlineDay record into the database
func (r *KlineDayRepository) Create(kline *model.KlineDay) error {
	return r.db.Create(kline).Error
}

func (r *KlineDayRepository) GetByStockCode(secCode string, marketType int16) ([]model.KlineDay, error) {
	var klines []model.KlineDay
	if err := r.db.Where("stock_code = ? AND market_type = ?", secCode, marketType).Order("market_date asc").Find(&klines).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return klines, nil
}

// GetByPrimaryKey retrieves a KlineDay record based on the unique constraints (sec_code, market_type, market_date)
func (r *KlineDayRepository) GetByPrimaryKey(secCode string, marketType int16, marketDate int32) (*model.KlineDay, error) {
	var kline model.KlineDay
	if err := r.db.Where("stock_code = ? AND market_type = ? AND market_date = ?", secCode, marketType, marketDate).First(&kline).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &kline, nil
}

// Update updates an existing KlineDay record
func (r *KlineDayRepository) Update(kline *model.KlineDay) error {
	return r.db.Save(kline).Error
}

// Delete removes a KlineDay record based on the unique constraints
func (r *KlineDayRepository) Delete(secCode string, marketType int16, marketDate int32) error {
	return r.db.Where("stock_code = ? AND market_type = ? AND market_date = ?", secCode, marketType, marketDate).Delete(&model.KlineDay{}).Error
}

// ListAll retrieves all KlineDay records from the database
func (r *KlineDayRepository) ListAll() ([]model.KlineDay, error) {
	var klines []model.KlineDay
	if err := r.db.Find(&klines).Error; err != nil {
		return nil, err
	}
	return klines, nil
}

func (r *KlineDayRepository) BatchUpsert(klines []*model.KlineDay) error {
	if len(klines) == 0 {
		return nil
	}

	err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "stock_code"},
			{Name: "market_type"},
			{Name: "market_date"},
		}, // 定义冲突字段
		DoUpdates: clause.AssignmentColumns([]string{
			"open",
			"close",
			"low",
			"high",
			"volume",
			"amount",
			"amplitude",
			"change",
			"changepct",
			"turnover",
			"update_time",
		}), // 冲突时更新的字段
	}).Create(&klines).Error

	return err
}
