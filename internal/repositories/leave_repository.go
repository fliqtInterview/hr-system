package repositories

import (
	"hr-system/config"
	"hr-system/internal/models"
)

type LeaveRepository struct{}

func NewLeaveRepository() *LeaveRepository {
	return &LeaveRepository{}
}

// Create 創建請假記錄
func (r *LeaveRepository) Create(leave *models.Leave) error {
	return config.DB.Create(leave).Error
}

// GetByID 根據ID獲取請假記錄
func (r *LeaveRepository) GetByID(id uint) (*models.Leave, error) {
	var leave models.Leave
	err := config.DB.Preload("Employee").First(&leave, id).Error
	if err != nil {
		return nil, err
	}
	return &leave, nil
}

// GetByEmployeeID 獲取員工的請假記錄
func (r *LeaveRepository) GetByEmployeeID(employeeID uint) ([]models.Leave, error) {
	var leaves []models.Leave
	err := config.DB.Where("employee_id = ?", employeeID).Find(&leaves).Error
	if err != nil {
		return nil, err
	}
	return leaves, nil
}

// Update 更新請假記錄
func (r *LeaveRepository) Update(leave *models.Leave) error {
	return config.DB.Save(leave).Error
}

// Delete 刪除請假記錄
func (r *LeaveRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Leave{}, id).Error
}

// GetPendingLeaves 獲取待審批的請假記錄
func (r *LeaveRepository) GetPendingLeaves() ([]models.Leave, error) {
	var leaves []models.Leave
	err := config.DB.Where("status = ?", "pending").
		Preload("Employee").
		Find(&leaves).Error
	if err != nil {
		return nil, err
	}
	return leaves, nil
}

// GetAll 獲取所有請假記錄
func (r *LeaveRepository) GetAll() ([]models.Leave, error) {
	var leaves []models.Leave
	err := config.DB.Preload("Employee").Find(&leaves).Error
	if err != nil {
		return nil, err
	}
	return leaves, nil
}
