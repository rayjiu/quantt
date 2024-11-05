package repository

import (
	"github.com/google/uuid"
	"github.com/rayjiu/quantt/data/internal/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SectorRepository 提供对 sectors 表的操作方法
type SectorRepository struct {
	db *gorm.DB
}

// NewSectorRepository 返回一个新的 SectorRepository 实例
func NewSectorRepository(db *gorm.DB) *SectorRepository {
	return &SectorRepository{db: db}
}

// Create 创建一个新的 Sector
func (r *SectorRepository) Create(sector *model.Sector) error {
	result := r.db.Create(sector)
	return result.Error
}

// GetByID 根据 SectorID 获取单个 Sector
func (r *SectorRepository) GetByID(sectorID uuid.UUID) (*model.Sector, error) {
	var sector model.Sector
	result := r.db.First(&sector, "sector_id = ?", sectorID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sector, nil
}

// GetBySecCode 获取指定 sec_code 和 sec_type 的 Sector
func (r *SectorRepository) GetBySecCode(secCode string, secType int16) (*model.Sector, error) {
	var sector model.Sector
	result := r.db.First(&sector, "sec_code = ? AND sec_type = ?", secCode, secType)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sector, nil
}

// Update 更新指定 Sector 的数据
func (r *SectorRepository) Update(sector *model.Sector) error {
	result := r.db.Save(sector)
	return result.Error
}

// Delete 删除指定 Sector
func (r *SectorRepository) Delete(sectorID uuid.UUID) error {
	result := r.db.Delete(&model.Sector{}, "sector_id = ?", sectorID)
	return result.Error
}

// GetAll 获取所有 Sector 数据
func (r *SectorRepository) GetAll() ([]model.Sector, error) {
	var sectors []model.Sector
	result := r.db.Find(&sectors)
	if result.Error != nil {
		return nil, result.Error
	}
	return sectors, nil
}

// BatchUpsert 批量插入或更新数据
func (r *SectorRepository) BatchUpsert(sectors []model.Sector) error {

	// 根据主键进行冲突处理，假设在 sec_code 和 sec_type 冲突时更新数据
	result := r.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "sec_code"}, {Name: "sec_type"}},
			DoUpdates: clause.AssignmentColumns([]string{"sec_name"}),
		},
	).Create(&sectors)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
