package models

import (
	"time"

	"gorm.io/gorm"
)

// Employee 員工模型
type Employee struct {
	gorm.Model
	Name             string    `gorm:"type:varchar(100);not null" json:"name"`              // 姓名
	Email            string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"` // 電子郵件
	Phone            string    `gorm:"type:varchar(20)" json:"phone"`                       // 電話
	Position         string    `gorm:"type:varchar(50)" json:"position"`                    // 職位
	Department       string    `gorm:"type:varchar(50)" json:"department"`                  // 部門
	Level            int       `json:"level"`                                               // 職等
	Salary           float64   `json:"salary"`                                              // 薪資
	HireDate         time.Time `json:"hire_date"`                                           // 入職日期
	Address          string    `gorm:"type:varchar(200)" json:"address"`                    // 地址
	EmergencyContact string    `gorm:"type:varchar(100)" json:"emergency_contact"`          // 緊急聯絡人
	Status           string    `gorm:"type:varchar(20);default:'active'" json:"status"`     // 狀態（active/inactive）
}
