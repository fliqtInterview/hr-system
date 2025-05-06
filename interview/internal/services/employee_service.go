package services

import (
	"errors"
	"hr-system/internal/models"
	"hr-system/internal/repositories"
)

type EmployeeService struct {
	repo *repositories.EmployeeRepository
}

func NewEmployeeService(repo *repositories.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// CreateEmployee 創建新員工
func (s *EmployeeService) CreateEmployee(employee *models.Employee) error {
	// 檢查員工編號是否已存在
	existing, _ := s.repo.GetByEmployeeID(employee.EmployeeID)
	if existing != nil {
		return errors.New("employee ID already exists")
	}
	return s.repo.Create(employee)
}

// GetEmployee 獲取員工信息
func (s *EmployeeService) GetEmployee(id uint) (*models.Employee, error) {
	return s.repo.GetByID(id)
}

// GetEmployeeByEmployeeID 根據員工編號獲取員工
func (s *EmployeeService) GetEmployeeByEmployeeID(employeeID string) (*models.Employee, error) {
	return s.repo.GetByEmployeeID(employeeID)
}

// UpdateEmployee 更新員工信息
func (s *EmployeeService) UpdateEmployee(employee *models.Employee) error {
	// 檢查員工是否存在
	existing, err := s.repo.GetByID(employee.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("employee not found")
	}
	return s.repo.Update(employee)
}

// DeleteEmployee 刪除員工
func (s *EmployeeService) DeleteEmployee(id uint) error {
	// 檢查員工是否存在
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("employee not found")
	}
	return s.repo.Delete(id)
}

// ListEmployees 獲取員工列表
func (s *EmployeeService) ListEmployees(page, pageSize int) ([]models.Employee, int64, error) {
	return s.repo.List(page, pageSize)
}
