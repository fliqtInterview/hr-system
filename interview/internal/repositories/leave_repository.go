package repositories

import (
	"hr-system/internal/models"

	"gorm.io/gorm"
)

type LeaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) *LeaveRepository {
	return &LeaveRepository{db: db}
}

// Create 創建請假記錄
func (r *LeaveRepository) Create(leave *models.Leave) error {
	return r.db.Create(leave).Error
}

// GetByID 根據ID獲取請假記錄
func (r *LeaveRepository) GetByID(id uint) (*models.Leave, error) {
	var leave models.Leave
	err := r.db.Preload("Employee").First(&leave, id).Error
	if err != nil {
		return nil, err
	}
	return &leave, nil
}

// Update 更新請假記錄
func (r *LeaveRepository) Update(leave *models.Leave) error {
	return r.db.Save(leave).Error
}

// Delete 刪除請假記錄
func (r *LeaveRepository) Delete(id uint) error {
	return r.db.Delete(&models.Leave{}, id).Error
}

// List 獲取請假記錄列表
func (r *LeaveRepository) List(page, pageSize int, employeeID *uint) ([]models.Leave, int64, error) {
	var leaves []models.Leave
	var total int64
	query := r.db.Model(&models.Leave{})

	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Employee").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&leaves).Error
	if err != nil {
		return nil, 0, err
	}

	return leaves, total, nil
}

// GetEmployeeLeavesByDateRange 獲取員工在指定日期範圍內的請假記錄
func (r *LeaveRepository) GetEmployeeLeavesByDateRange(employeeID uint, startDate, endDate string) ([]models.Leave, error) {
	var leaves []models.Leave
	err := r.db.Where("employee_id = ? AND start_date <= ? AND end_date >= ?", employeeID, endDate, startDate).
		Find(&leaves).Error
	return leaves, err
}

// GetEmployeeLeaveSummary 獲取員工請假統計
func (r *LeaveRepository) GetEmployeeLeaveSummary(employeeID uint, year int) (map[models.LeaveType]float32, error) {
	var leaves []models.Leave
	summary := make(map[models.LeaveType]float32)

	err := r.db.Where("employee_id = ? AND YEAR(start_date) = ? AND status = ?",
		employeeID, year, models.LeaveStatusApproved).
		Find(&leaves).Error
	if err != nil {
		return nil, err
	}

	for _, leave := range leaves {
		summary[leave.Type] += leave.Duration
	}

	return summary, nil
}
