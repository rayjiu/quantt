package service

import (
	"fmt"
	"sync"

	"github.com/rayjiu/quantt/data/internal/db/model"
	"github.com/rayjiu/quantt/data/internal/db/repository"
)

var SecotorService sectorService

// sectorService 定义服务层，处理业务逻辑
type sectorService struct {
	sectorRepo *repository.SectorRepository
	cache      map[string]*model.Sector
	mu         sync.RWMutex
}

// NewSectorService 创建一个新的 sectorService 实例
func NewSectorService(sectorRepo *repository.SectorRepository) *sectorService {
	return &sectorService{
		sectorRepo: sectorRepo,
		cache:      make(map[string]*model.Sector),
	}
}

// 缓存刷新方法，假设可以周期性地刷新缓存
func (s *sectorService) RefreshCache() error {
	sectors, err := s.sectorRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to refresh cache: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 清空旧的缓存
	s.cache = make(map[string]*model.Sector)

	// 填充新的缓存
	for _, sector := range sectors {
		cacheKey := fmt.Sprintf("%s-%d", sector.SecCode, sector.SecType)
		s.cache[cacheKey] = &sector
	}
	return nil
}

// batchUpsert 批量插入数据，并且根据缓存进行去重操作
func (s *sectorService) BatchUpsert(sectors []model.Sector) error {
	// 使用局部缓存防止并发冲突
	uniqueSectors := make([]model.Sector, 0)
	seen := make(map[string]bool)

	s.mu.RLock()
	// 遍历传入的 sector 数据，并根据业务逻辑进行去重
	for _, sector := range sectors {
		cacheKey := fmt.Sprintf("%s-%d", sector.SecCode, sector.SecType)

		// 如果缓存中存在，判断是否重复，若重复则进行逻辑去重
		if _, exists := s.cache[cacheKey]; exists {
			// 数据存在，进行去重处理（如仅更新名称，跳过插入）
			// 假设我们这里的去重逻辑是：如果名称不同，则更新
			existingSector := s.cache[cacheKey]
			if existingSector.SecName != sector.SecName {
				existingSector.SecName = sector.SecName // 执行更新逻辑
				uniqueSectors = append(uniqueSectors, *existingSector)
			}
		} else {
			// 如果缓存中没有，则视为新数据
			if _, ok := seen[cacheKey]; !ok {
				uniqueSectors = append(uniqueSectors, sector)
				seen[cacheKey] = true
			}
		}
	}
	s.mu.RUnlock()

	// 执行批量插入
	if len(uniqueSectors) > 0 {
		err := s.sectorRepo.BatchUpsert(uniqueSectors)
		if err != nil {
			return fmt.Errorf("failed to batch upsert: %v", err)
		}
	}

	return nil
}
