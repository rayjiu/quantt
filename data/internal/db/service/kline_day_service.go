package service

import (
	"errors"

	"github.com/rayjiu/quantt/data/internal/db/model"
	"github.com/rayjiu/quantt/data/internal/db/repository"
)

var KlineService klineDayService

type klineDayService struct {
	repo *repository.KlineDayRepository
}

// NewklineDayService creates a new klineDayService instance
func NewklineDayService(repo *repository.KlineDayRepository) *klineDayService {
	return &klineDayService{repo: repo}
}

// CreateKline adds a new KlineDay record
func (s *klineDayService) CreateKline(kline *model.KlineDay) error {
	existing, err := s.repo.GetByPrimaryKey(kline.StockCode, kline.MarketType, kline.MarketDate)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("record already exists")
	}
	return s.repo.Create(kline)
}

// GetKline retrieves a KlineDay record by its primary key fields
func (s *klineDayService) GetKline(secCode string, marketType int16, marketDate int32) (*models.KlineDay, error) {
	return s.repo.GetByPrimaryKey(secCode, marketType, marketDate)
}

// UpdateKline updates an existing KlineDay record
func (s *klineDayService) UpdateKline(kline *model.KlineDay) error {
	return s.repo.Update(kline)
}

// DeleteKline deletes a KlineDay record by its primary key fields
func (s *klineDayService) DeleteKline(secCode string, marketType int16, marketDate int32) error {
	return s.repo.Delete(secCode, marketType, marketDate)
}

// ListAllKlines retrieves all KlineDay records
func (s *klineDayService) ListAllKlines() ([]model.KlineDay, error) {
	return s.repo.ListAll()
}

// BatchUpsert performs batch upsert on sec_fund_flow records
func (s *klineDayService) BatchUpsert(klines []*model.KlineDay) error {
	return s.repo.BatchUpsert(klines)
}
