package models

import (
	"time"

	"gorm.io/gorm"
)

// LeaveType 請假類型
type LeaveType string

const (
	LeaveTypeAnnual      LeaveType = "annual"      // 年假
	LeaveTypeSick        LeaveType = "sick"        // 病假
	LeaveTypePersonal    LeaveType = "personal"    // 事假
	LeaveTypeMaternity   LeaveType = "maternity"   // 產假
	LeaveTypeMarriage    LeaveType = "marriage"    // 婚假
	LeaveTypeBereavement LeaveType = "bereavement" // 喪假
)

// LeaveStatus 請假狀態
type LeaveStatus string

const (
	LeaveStatusPending   LeaveStatus = "pending"   // 待審核
	LeaveStatusApproved  LeaveStatus = "approved"  // 已批准
	LeaveStatusRejected  LeaveStatus = "rejected"  // 已拒絕
	LeaveStatusCancelled LeaveStatus = "cancelled" // 已取消
)

// Leave 請假記錄
type Leave struct {
	gorm.Model
	EmployeeID   uint        `gorm:"not null;index" json:"employee_id"`                                 // 員工ID
	Employee     Employee    `gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE" json:"employee"` // 員工信息
	Type         LeaveType   `gorm:"type:varchar(20);not null" json:"type"`                             // 請假類型
	StartDate    time.Time   `gorm:"not null" json:"start_date"`                                        // 開始日期
	EndDate      time.Time   `gorm:"not null" json:"end_date"`                                          // 結束日期
	Duration     float32     `gorm:"not null" json:"duration"`                                          // 請假天數
	Reason       string      `gorm:"type:varchar(500)" json:"reason"`                                   // 請假原因
	Status       LeaveStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`                  // 狀態
	ApprovedBy   *uint       `json:"approved_by"`                                                       // 審批人ID
	ApprovedAt   *time.Time  `json:"approved_at"`                                                       // 審批時間
	RejectReason string      `gorm:"type:varchar(500)" json:"reject_reason"`                            // 拒絕原因
}
