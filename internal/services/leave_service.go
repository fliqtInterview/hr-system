package services

import (
	"context"
	"errors"
	"log"
	"time"

	"hr-system/internal/models"
	"hr-system/internal/repositories"
)

type LeaveService struct {
	leaveRepo    *repositories.LeaveRepository
	employeeRepo *repositories.EmployeeRepository
	cacheService *CacheService
}

func NewLeaveService(leaveRepo *repositories.LeaveRepository, employeeRepo *repositories.EmployeeRepository, cacheService *CacheService) *LeaveService {
	return &LeaveService{
		leaveRepo:    leaveRepo,
		employeeRepo: employeeRepo,
		cacheService: cacheService,
	}
}

// CreateLeave 創建請假記錄
func (s *LeaveService) CreateLeave(leave *models.Leave) error {
	// 檢查員工是否存在
	_, err := s.employeeRepo.GetByID(leave.EmployeeID)
	if err != nil {
		return errors.New("employee not found")
	}

	// 檢查日期是否有效
	if leave.StartDate.After(leave.EndDate) {
		return errors.New("start date must be before end date")
	}

	// 檢查是否有重疊的請假記錄
	// TODO: 實現日期重疊檢查

	if err := s.leaveRepo.Create(leave); err != nil {
		return err
	}

	// 添加到緩存
	ctx := context.Background()
	if err := s.cacheService.SetLeave(ctx, leave); err != nil {
		// 緩存失敗不影響主流程，只記錄日誌
		log.Printf("Failed to cache leave: %v", err)
	}

	return nil
}

// GetLeave 獲取請假記錄
func (s *LeaveService) GetLeave(id uint) (*models.Leave, error) {
	ctx := context.Background()

	// 先從緩存獲取
	leave, err := s.cacheService.GetLeave(ctx, id)
	if err == nil {
		return leave, nil
	}

	// 緩存未命中，從數據庫獲取
	leave, err = s.leaveRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 添加到緩存
	if err := s.cacheService.SetLeave(ctx, leave); err != nil {
		log.Printf("Failed to cache leave: %v", err)
	}

	return leave, nil
}

// UpdateLeaveStatus 更新請假狀態
func (s *LeaveService) UpdateLeaveStatus(id uint, status string, remark string) error {
	leave, err := s.leaveRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 檢查狀態是否有效
	if status != "approved" && status != "rejected" {
		return errors.New("invalid status")
	}

	// 更新狀態
	leave.Status = status
	leave.ApproveRemark = remark
	now := time.Now()
	leave.ApproveTime = &now

	if err := s.leaveRepo.Update(leave); err != nil {
		return err
	}

	// 更新緩存
	ctx := context.Background()
	if err := s.cacheService.SetLeave(ctx, leave); err != nil {
		log.Printf("Failed to update leave cache: %v", err)
	}

	return nil
}

// DeleteLeave 刪除請假記錄
func (s *LeaveService) DeleteLeave(id uint) error {
	if err := s.leaveRepo.Delete(id); err != nil {
		return err
	}

	// 刪除緩存
	ctx := context.Background()
	if err := s.cacheService.DeleteLeave(ctx, id); err != nil {
		log.Printf("Failed to delete leave cache: %v", err)
	}

	return nil
}

// ListLeaves 獲取所有請假記錄
func (s *LeaveService) ListLeaves() ([]models.Leave, error) {
	return s.leaveRepo.GetAll()
}
