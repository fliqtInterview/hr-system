package services

import (
	"errors"
	"time"

	"hr-system/internal/models"
	"hr-system/internal/repositories"
)

type LeaveService struct {
	leaveRepo    *repositories.LeaveRepository
	employeeRepo *repositories.EmployeeRepository
}

func NewLeaveService(leaveRepo *repositories.LeaveRepository, employeeRepo *repositories.EmployeeRepository) *LeaveService {
	return &LeaveService{
		leaveRepo:    leaveRepo,
		employeeRepo: employeeRepo,
	}
}

// CreateLeave 創建請假申請
func (s *LeaveService) CreateLeave(leave *models.Leave) error {
	// 檢查員工是否存在
	employee, err := s.employeeRepo.GetByID(leave.EmployeeID)
	if err != nil {
		return errors.New("employee not found")
	}

	// 檢查請假日期是否合理
	if leave.StartDate.After(leave.EndDate) {
		return errors.New("start date cannot be after end date")
	}

	// 檢查是否與現有請假時間重疊
	existingLeaves, err := s.leaveRepo.GetEmployeeLeavesByDateRange(
		leave.EmployeeID,
		leave.StartDate.Format("2006-01-02"),
		leave.EndDate.Format("2006-01-02"),
	)
	if err != nil {
		return err
	}
	if len(existingLeaves) > 0 {
		return errors.New("leave period overlaps with existing leave")
	}

	// 計算請假天數
	leave.Duration = s.calculateLeaveDuration(leave.StartDate, leave.EndDate)

	// 檢查年假額度（如果是年假）
	if leave.Type == models.LeaveTypeAnnual {
		summary, err := s.leaveRepo.GetEmployeeLeaveSummary(leave.EmployeeID, time.Now().Year())
		if err != nil {
			return err
		}
		// 假設每年有15天年假
		if summary[models.LeaveTypeAnnual]+leave.Duration > 15 {
			return errors.New("annual leave quota exceeded")
		}
	}

	// 設置初始狀態
	leave.Status = models.LeaveStatusPending
	leave.Employee = *employee

	return s.leaveRepo.Create(leave)
}

// GetLeave 獲取請假記錄
func (s *LeaveService) GetLeave(id uint) (*models.Leave, error) {
	return s.leaveRepo.GetByID(id)
}

// UpdateLeaveStatus 更新請假狀態
func (s *LeaveService) UpdateLeaveStatus(id uint, status models.LeaveStatus, approverID uint, rejectReason string) error {
	leave, err := s.leaveRepo.GetByID(id)
	if err != nil {
		return err
	}

	if leave.Status != models.LeaveStatusPending {
		return errors.New("can only update pending leave requests")
	}

	leave.Status = status
	now := time.Now()
	leave.ApprovedAt = &now
	leave.ApprovedBy = &approverID

	if status == models.LeaveStatusRejected {
		if rejectReason == "" {
			return errors.New("reject reason is required")
		}
		leave.RejectReason = rejectReason
	}

	return s.leaveRepo.Update(leave)
}

// ListLeaves 獲取請假記錄列表
func (s *LeaveService) ListLeaves(page, pageSize int, employeeID *uint) ([]models.Leave, int64, error) {
	return s.leaveRepo.List(page, pageSize, employeeID)
}

// CancelLeave 取消請假申請
func (s *LeaveService) CancelLeave(id uint, employeeID uint) error {
	leave, err := s.leaveRepo.GetByID(id)
	if err != nil {
		return err
	}

	if leave.EmployeeID != employeeID {
		return errors.New("unauthorized to cancel this leave")
	}

	if leave.Status != models.LeaveStatusPending {
		return errors.New("can only cancel pending leave requests")
	}

	leave.Status = models.LeaveStatusCancelled
	return s.leaveRepo.Update(leave)
}

// calculateLeaveDuration 計算請假天數
func (s *LeaveService) calculateLeaveDuration(startDate, endDate time.Time) float32 {
	days := endDate.Sub(startDate).Hours() / 24
	return float32(days) + 1 // 包含開始和結束日
}

// DeleteLeave 刪除請假記錄
func (s *LeaveService) DeleteLeave(id uint) error {
	return s.leaveRepo.Delete(id)
}
