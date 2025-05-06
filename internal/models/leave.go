package models

import (
	"time"

	"gorm.io/gorm"
)

// Leave 請假記錄模型
type Leave struct {
	gorm.Model
	EmployeeID    uint       `gorm:"not null" json:"employee_id"`                      // 員工ID
	Employee      Employee   `gorm:"foreignKey:EmployeeID" json:"employee"`            // 關聯員工
	StartDate     time.Time  `json:"start_date"`                                       // 開始日期
	EndDate       time.Time  `json:"end_date"`                                         // 結束日期
	LeaveType     string     `gorm:"type:varchar(20);not null" json:"leave_type"`      // 請假類型（年假/病假/事假等）
	Reason        string     `gorm:"type:text" json:"reason"`                          // 請假原因
	Status        string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // 狀態（pending/approved/rejected）
	ApproverID    *uint      `json:"approver_id,omitempty"`                            // 審批人ID
	ApproveTime   *time.Time `json:"approve_time,omitempty"`                           // 審批時間
	ApproveRemark string     `gorm:"type:text" json:"approve_remark"`                  // 審批備註
}
