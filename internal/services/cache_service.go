package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"hr-system/config"
	"hr-system/internal/models"
)

type CacheService struct{}

func NewCacheService() *CacheService {
	return &CacheService{}
}

// 員工緩存操作
func (s *CacheService) GetEmployee(ctx context.Context, id uint) (*models.Employee, error) {
	key := fmt.Sprintf("%s%d", config.EmployeeKeyPrefix, id)
	data, err := config.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var employee models.Employee
	if err := json.Unmarshal(data, &employee); err != nil {
		return nil, err
	}

	return &employee, nil
}

func (s *CacheService) SetEmployee(ctx context.Context, employee *models.Employee) error {
	key := fmt.Sprintf("%s%d", config.EmployeeKeyPrefix, employee.ID)
	data, err := json.Marshal(employee)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(ctx, key, data, config.EmployeeCacheExpiration).Err()
}

func (s *CacheService) DeleteEmployee(ctx context.Context, id uint) error {
	key := fmt.Sprintf("%s%d", config.EmployeeKeyPrefix, id)
	return config.RedisClient.Del(ctx, key).Err()
}

// 請假記錄緩存操作
func (s *CacheService) GetLeave(ctx context.Context, id uint) (*models.Leave, error) {
	key := fmt.Sprintf("%s%d", config.LeaveKeyPrefix, id)
	data, err := config.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var leave models.Leave
	if err := json.Unmarshal(data, &leave); err != nil {
		return nil, err
	}

	return &leave, nil
}

func (s *CacheService) SetLeave(ctx context.Context, leave *models.Leave) error {
	key := fmt.Sprintf("%s%d", config.LeaveKeyPrefix, leave.ID)
	data, err := json.Marshal(leave)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(ctx, key, data, config.LeaveCacheExpiration).Err()
}

func (s *CacheService) DeleteLeave(ctx context.Context, id uint) error {
	key := fmt.Sprintf("%s%d", config.LeaveKeyPrefix, id)
	return config.RedisClient.Del(ctx, key).Err()
}

// PrewarmCache 預熱緩存
func (s *CacheService) PrewarmCache(ctx context.Context, employees []models.Employee, leaves []models.Leave) error {
	// 使用管道批量寫入緩存
	pipeline := config.RedisClient.Pipeline()

	// 預熱員工數據
	for _, employee := range employees {
		key := fmt.Sprintf("%s%d", config.EmployeeKeyPrefix, employee.ID)
		data, err := json.Marshal(employee)
		if err != nil {
			log.Printf("Failed to marshal employee data: %v", err)
			continue
		}
		pipeline.Set(ctx, key, data, config.EmployeeCacheExpiration)
	}

	// 預熱請假記錄數據
	for _, leave := range leaves {
		key := fmt.Sprintf("%s%d", config.LeaveKeyPrefix, leave.ID)
		data, err := json.Marshal(leave)
		if err != nil {
			log.Printf("Failed to marshal leave data: %v", err)
			continue
		}
		pipeline.Set(ctx, key, data, config.LeaveCacheExpiration)
	}

	// 執行管道命令
	_, err := pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute pipeline: %v", err)
	}

	log.Printf("Successfully prewarmed cache with %d employees and %d leaves", len(employees), len(leaves))
	return nil
}
