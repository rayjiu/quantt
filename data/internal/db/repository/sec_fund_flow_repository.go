package repository

import (
	"errors"

	"github.com/rayjiu/quantt/data/internal/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SecFundFlowRepository defines methods for accessing the sec_fund_flow table
type SecFundFlowRepository struct {
	db *gorm.DB
}

// NewSecFundFlowRepository creates a new SecFundFlowRepository instance
func NewSecFundFlowRepository(db *gorm.DB) *SecFundFlowRepository {
	return &SecFundFlowRepository{db: db}
}

// GetAll retrieves all records from the sec_fund_flow table
func (repo *SecFundFlowRepository) GetAll() ([]model.SecFundFlow, error) {
	var secFundFlows []model.SecFundFlow
	if err := repo.db.Find(&secFundFlows).Error; err != nil {
		return nil, err
	}
	return secFundFlows, nil
}

// GetBySecCode retrieves a sec_fund_flow record by sec_code
func (repo *SecFundFlowRepository) GetBySecCode(secCode string) (*model.SecFundFlow, error) {
	var secFundFlow model.SecFundFlow
	if err := repo.db.Where("sec_code = ?", secCode).First(&secFundFlow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &secFundFlow, nil
}

// Create adds a new sec_fund_flow record
func (repo *SecFundFlowRepository) Create(secFundFlow *model.SecFundFlow) error {
	return repo.db.Create(secFundFlow).Error
}

// Update updates an existing sec_fund_flow record
func (repo *SecFundFlowRepository) Update(secFundFlow *model.SecFundFlow) error {
	return repo.db.Save(secFundFlow).Error
}

// Delete removes a sec_fund_flow record by sec_code
func (repo *SecFundFlowRepository) Delete(secCode string) error {
	return repo.db.Where("sec_code = ?", secCode).Delete(&model.SecFundFlow{}).Error
}

func (repo *SecFundFlowRepository) BatchUpsert(secFundFlows []*model.SecFundFlow) error {
	err := repo.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "sec_code"}, {Name: "market_date"}}, // 定义冲突字段
		DoUpdates: clause.AssignmentColumns([]string{
			"main_buyer_netin",
			"main_buyer_ratio",
			"l1_netin",
			"l1_netin_ratio",
			"l2_netin",
			"l2_netin_ratio",
			"l3_netin",
			"l3_netin_ratio",
			"l4_netint",
			"l4_net_ratio",
			"update_time",
		}), // 冲突时更新的字段
	}).Create(secFundFlows).Error

	return err
}
