package repositories

import (
	"hr-system/config"
	"hr-system/internal/models"
)

type EmployeeRepository struct{}

func NewEmployeeRepository() *EmployeeRepository {
	return &EmployeeRepository{}
}

// Create 創建員工
func (r *EmployeeRepository) Create(employee *models.Employee) error {
	return config.DB.Create(employee).Error
}

// GetByID 根據ID獲取員工
func (r *EmployeeRepository) GetByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	err := config.DB.First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetByEmail 根據郵箱獲取員工
func (r *EmployeeRepository) GetByEmail(email string) (*models.Employee, error) {
	var employee models.Employee
	err := config.DB.Where("email = ?", email).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// Update 更新員工信息
func (r *EmployeeRepository) Update(employee *models.Employee) error {
	return config.DB.Save(employee).Error
}

// Delete 刪除員工
func (r *EmployeeRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Employee{}, id).Error
}

// GetAll 獲取所有員工
func (r *EmployeeRepository) GetAll() ([]models.Employee, error) {
	var employees []models.Employee
	err := config.DB.Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}
