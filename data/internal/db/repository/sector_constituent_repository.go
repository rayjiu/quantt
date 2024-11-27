package repository

import (
	"github.com/rayjiu/quantt/data/internal/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SectorConsitutuentRepository struct {
	db *gorm.DB
}

// NewSectorConsitutuentRepository creates a new SectorConsitutuentRepository instance
func NewSectorConsitutuentRepository(db *gorm.DB) *SectorConsitutuentRepository {
	return &SectorConsitutuentRepository{db: db}
}

func (s *SectorConsitutuentRepository) BatchUpsert(cons []*model.SectorConstituent) error {
	if len(cons) == 0 {
		return nil
	}

	err := s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "sec_code"},
			{Name: "stock_code"},
			{Name: "market_type"},
		}, // 定义冲突字段
		DoUpdates: clause.AssignmentColumns([]string{
			"sec_code",
			"stock_code",
			"market_type",
			"update_time",
		}),
	}).Create(&cons).Error

	return err
}
