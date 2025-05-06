package services

import (
	"context"
	"errors"
	"log"

	"hr-system/internal/models"
	"hr-system/internal/repositories"
)

type EmployeeService struct {
	employeeRepo *repositories.EmployeeRepository
	cacheService *CacheService
}

func NewEmployeeService(employeeRepo *repositories.EmployeeRepository, cacheService *CacheService) *EmployeeService {
	return &EmployeeService{
		employeeRepo: employeeRepo,
		cacheService: cacheService,
	}
}

// CreateEmployee 創建員工
func (s *EmployeeService) CreateEmployee(employee *models.Employee) error {
	// 檢查郵箱是否已存在
	existingEmployee, err := s.employeeRepo.GetByEmail(employee.Email)
	if err == nil && existingEmployee != nil {
		return errors.New("email already exists")
	}

	if err := s.employeeRepo.Create(employee); err != nil {
		return err
	}

	// 添加到緩存
	ctx := context.Background()
	if err := s.cacheService.SetEmployee(ctx, employee); err != nil {
		// 緩存失敗不影響主流程，只記錄日誌
		log.Printf("Failed to cache employee: %v", err)
	}

	return nil
}

// GetEmployee 獲取員工信息
func (s *EmployeeService) GetEmployee(id uint) (*models.Employee, error) {
	ctx := context.Background()

	// 先從緩存獲取
	employee, err := s.cacheService.GetEmployee(ctx, id)
	if err == nil {
		return employee, nil
	}

	// 緩存未命中，從數據庫獲取
	employee, err = s.employeeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 添加到緩存
	if err := s.cacheService.SetEmployee(ctx, employee); err != nil {
		log.Printf("Failed to cache employee: %v", err)
	}

	return employee, nil
}

// UpdateEmployee 更新員工信息
func (s *EmployeeService) UpdateEmployee(employee *models.Employee) error {
	// 獲取原來的員工信息
	oldEmployee, err := s.employeeRepo.GetByID(employee.ID)
	if err != nil {
		return err
	}

	// 如果郵箱發生變化，檢查新郵箱是否已存在
	if oldEmployee.Email != employee.Email {
		existingEmployee, err := s.employeeRepo.GetByEmail(employee.Email)
		if err == nil && existingEmployee != nil {
			return errors.New("email already exists")
		}
	}

	if err := s.employeeRepo.Update(employee); err != nil {
		return err
	}

	// 更新緩存
	ctx := context.Background()
	if err := s.cacheService.SetEmployee(ctx, employee); err != nil {
		log.Printf("Failed to update employee cache: %v", err)
	}

	return nil
}

// DeleteEmployee 刪除員工
func (s *EmployeeService) DeleteEmployee(id uint) error {
	if err := s.employeeRepo.Delete(id); err != nil {
		return err
	}

	// 刪除緩存
	ctx := context.Background()
	if err := s.cacheService.DeleteEmployee(ctx, id); err != nil {
		log.Printf("Failed to delete employee cache: %v", err)
	}

	return nil
}

// ListEmployees 獲取所有員工
func (s *EmployeeService) ListEmployees() ([]models.Employee, error) {
	return s.employeeRepo.GetAll()
}
