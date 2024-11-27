package service

import (
	"github.com/rayjiu/quantt/data/internal/db/model"
	"github.com/rayjiu/quantt/data/internal/db/repository"
)

var SectorConsService sectorConsService

type sectorConsService struct {
	repo *repository.SectorConsitutuentRepository
}

// NewSectorConsService creates a new klineDayService instance
func NewSectorConsService(repo *repository.SectorConsitutuentRepository) *sectorConsService {
	return &sectorConsService{repo: repo}
}

// BatchUpsert performs batch upsert on sec_constituent records
func (s *sectorConsService) BatchUpsert(cons []*model.SectorConstituent) error {
	return s.repo.BatchUpsert(cons)
}
