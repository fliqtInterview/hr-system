package repositories

import (
	"hr-system/internal/models"
	"time"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// Create 創建新員工
func (r *EmployeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

// GetByID 根據ID獲取員工
func (r *EmployeeRepository) GetByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetByEmployeeID 根據員工編號獲取員工
func (r *EmployeeRepository) GetByEmployeeID(employeeID string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Where("employee_id = ?", employeeID).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// Update 更新員工信息
func (r *EmployeeRepository) Update(employee *models.Employee) error {
	// 獲取現有記錄
	var existingEmployee models.Employee
	if err := r.db.First(&existingEmployee, employee.ID).Error; err != nil {
		return err
	}

	// 保留原有的時間戳和狀態
	employee.CreatedAt = existingEmployee.CreatedAt
	employee.UpdatedAt = time.Now()
	if employee.Status == "" {
		employee.Status = existingEmployee.Status
	}

	// 更新記錄，排除零值字段
	return r.db.Model(&existingEmployee).Updates(map[string]interface{}{
		"employee_id":       employee.EmployeeID,
		"name":              employee.Name,
		"email":             employee.Email,
		"phone":             employee.Phone,
		"position":          employee.Position,
		"department":        employee.Department,
		"level":             employee.Level,
		"salary":            employee.Salary,
		"hire_date":         employee.HireDate,
		"address":           employee.Address,
		"emergency_contact": employee.EmergencyContact,
		"status":            employee.Status,
		"updated_at":        employee.UpdatedAt,
	}).Error
}

// Delete 刪除員工
func (r *EmployeeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Employee{}, id).Error
}

// List 獲取員工列表
func (r *EmployeeRepository) List(page, pageSize int) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	err := r.db.Model(&models.Employee{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&employees).Error
	if err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}
