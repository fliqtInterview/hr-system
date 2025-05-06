package services

import (
	"context"
	"log"
	"time"

	"hr-system/internal/repositories"
)

type PrewarmService struct {
	employeeRepo *repositories.EmployeeRepository
	leaveRepo    *repositories.LeaveRepository
	cacheService *CacheService
}

func NewPrewarmService(
	employeeRepo *repositories.EmployeeRepository,
	leaveRepo *repositories.LeaveRepository,
	cacheService *CacheService,
) *PrewarmService {
	return &PrewarmService{
		employeeRepo: employeeRepo,
		leaveRepo:    leaveRepo,
		cacheService: cacheService,
	}
}

// StartPrewarming 開始預熱緩存
func (s *PrewarmService) StartPrewarming(ctx context.Context) {
	// 立即執行一次預熱
	s.prewarmCache(ctx)

	// 每30分鐘執行一次預熱
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				s.prewarmCache(ctx)
			}
		}
	}()
}

// prewarmCache 執行緩存預熱
func (s *PrewarmService) prewarmCache(ctx context.Context) {
	log.Println("Starting cache prewarming...")

	// 獲取所有員工數據
	employees, err := s.employeeRepo.GetAll()
	if err != nil {
		log.Printf("Failed to get employees for prewarming: %v", err)
		return
	}

	// 獲取所有請假記錄
	leaves, err := s.leaveRepo.GetAll()
	if err != nil {
		log.Printf("Failed to get leaves for prewarming: %v", err)
		return
	}

	// 執行預熱
	if err := s.cacheService.PrewarmCache(ctx, employees, leaves); err != nil {
		log.Printf("Failed to prewarm cache: %v", err)
		return
	}
}
