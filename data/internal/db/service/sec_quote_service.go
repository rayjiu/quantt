package service

import (
	"github.com/rayjiu/quantt/data/internal/db/model"
	"github.com/rayjiu/quantt/data/internal/db/repository"
)

var SecQuoteService secQuoteService

// SecQuoteService defines the business logic for SecQuote
type secQuoteService struct {
	repo *repository.SecQuoteRepository
}

// NewSecQuoteService creates a new SecQuoteService instance
func NewSecQuoteService(repo *repository.SecQuoteRepository) *secQuoteService {
	return &secQuoteService{repo: repo}
}

// GetAllSecQuotes retrieves all sec_quote records
func (s *secQuoteService) GetAllSecQuotes() ([]model.SecQuote, error) {
	return s.repo.GetAll()
}

// GetSecQuoteByCode retrieves a sec_quote record by code
func (s *secQuoteService) GetSecQuoteByCode(code string) (*model.SecQuote, error) {
	return s.repo.GetByCode(code)
}

// CreateSecQuote adds a new sec_quote record
func (s *secQuoteService) CreateSecQuote(secQuote *model.SecQuote) error {
	return s.repo.Create(secQuote)
}

// UpdateSecQuote updates an existing sec_quote record
func (s *secQuoteService) UpdateSecQuote(secQuote *model.SecQuote) error {
	return s.repo.Update(secQuote)
}

// DeleteSecQuote deletes a sec_quote record by code
func (s *secQuoteService) DeleteSecQuote(code string) error {
	return s.repo.Delete(code)
}

// BatchUpsert batch update the sec_quote records if the records exits; or else, insert it.
func (s *secQuoteService) BatchUpsert(secQuotes []*model.SecQuote) error {
	return s.repo.BatchUpsert(secQuotes)
}
