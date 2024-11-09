package service

import (
	"github.com/rayjiu/quantt/data/internal/db/model"
	"github.com/rayjiu/quantt/data/internal/db/repository"
)

var SecFundFlowService secFundFlowService

// SecFundFlowService defines the business logic for SecFundFlow
type secFundFlowService struct {
	repo *repository.SecFundFlowRepository
}

// NewSecFundFlowService creates a new SecFundFlowService instance
func NewSecFundFlowService(repo *repository.SecFundFlowRepository) *secFundFlowService {
	return &secFundFlowService{repo: repo}
}

// GetAllSecFundFlows retrieves all sec_fund_flow records
func (s *secFundFlowService) GetAllSecFundFlows() ([]model.SecFundFlow, error) {
	return s.repo.GetAll()
}

// GetSecFundFlowBySecCode retrieves a sec_fund_flow record by sec_code
func (s *secFundFlowService) GetSecFundFlowBySecCode(secCode string) (*model.SecFundFlow, error) {
	return s.repo.GetBySecCode(secCode)
}

// CreateSecFundFlow adds a new sec_fund_flow record
func (s *secFundFlowService) CreateSecFundFlow(secFundFlow *model.SecFundFlow) error {
	return s.repo.Create(secFundFlow)
}

// UpdateSecFundFlow updates an existing sec_fund_flow record
func (s *secFundFlowService) UpdateSecFundFlow(secFundFlow *model.SecFundFlow) error {
	return s.repo.Update(secFundFlow)
}

// DeleteSecFundFlow deletes a sec_fund_flow record by sec_code
func (s *secFundFlowService) DeleteSecFundFlow(secCode string) error {
	return s.repo.Delete(secCode)
}

// BatchUpsert performs batch upsert on sec_fund_flow records
func (s *secFundFlowService) BatchUpsert(secFundFlows []*model.SecFundFlow) error {
	return s.repo.BatchUpsert(secFundFlows)
}
