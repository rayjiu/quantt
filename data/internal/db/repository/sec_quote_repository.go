package repository

import (
	"errors"

	"github.com/rayjiu/quantt/data/internal/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SecQuoteRepository defines methods for accessing the sec_quote table
type SecQuoteRepository struct {
	db *gorm.DB
}

// NewSecQuoteRepository creates a new SecQuoteRepository instance
func NewSecQuoteRepository(db *gorm.DB) *SecQuoteRepository {
	return &SecQuoteRepository{db: db}
}

// GetAll retrieves all records from the sec_quote table
func (repo *SecQuoteRepository) GetAll() ([]model.SecQuote, error) {
	var secQuotes []model.SecQuote
	if err := repo.db.Find(&secQuotes).Error; err != nil {
		return nil, err
	}
	return secQuotes, nil
}

// GetByCode retrieves a sec_quote record by code
func (repo *SecQuoteRepository) GetByCode(code string) (*model.SecQuote, error) {
	var secQuote model.SecQuote
	if err := repo.db.Where("code = ?", code).First(&secQuote).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &secQuote, nil
}

// Create adds a new sec_quote record
func (repo *SecQuoteRepository) Create(secQuote *model.SecQuote) error {
	return repo.db.Create(secQuote).Error
}

func (repo *SecQuoteRepository) BatchUpsert(secQuotes []*model.SecQuote) error {
	// 使用 GORM 的批量插入特性和 ON CONFLICT 子句
	err := repo.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "sec_code"}, {Name: "market_type"}},                                    // 指定冲突字段
		DoUpdates: clause.AssignmentColumns([]string{"last_price", "chg_ratio", "exchg_ratio", "total_mkt_cap"}), // 更新这些字段
	}).Create(secQuotes).Error

	return err
}

// Update updates an existing sec_quote record
func (repo *SecQuoteRepository) Update(secQuote *model.SecQuote) error {
	return repo.db.Save(secQuote).Error
}

// Delete removes a sec_quote record by code
func (repo *SecQuoteRepository) Delete(code string) error {
	return repo.db.Where("code = ?", code).Delete(&model.SecQuote{}).Error
}
