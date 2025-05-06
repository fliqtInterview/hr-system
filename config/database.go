package config

import (
	"fmt"
	"log"
	"time"

	"hr-system/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化數據庫連接
func InitDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "hr_system"),
	)

	// 嘗試連接數據庫，最多重試5次
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database, retrying in 5 seconds... (attempt %d/5)", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after 5 attempts:", err)
	}

	// 自動遷移數據庫結構
	err = DB.AutoMigrate(
		&models.Employee{},
		&models.Leave{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// ====== SEED DATA START ======
	var count int64
	DB.Model(&models.Employee{}).Count(&count)
	if count == 0 {
		employees := []models.Employee{
			{
				Name:             "王小明",
				Email:            "xiaoming.wang@example.com",
				Phone:            "0912345678",
				Position:         "工程師",
				Department:       "研發部",
				Level:            1,
				Salary:           60000,
				HireDate:         time.Date(2022, 1, 10, 0, 0, 0, 0, time.Local),
				Address:          "台北市信義區",
				EmergencyContact: "王媽媽 0911222333",
				Status:           "active",
			},
			{
				Name:             "陳美麗",
				Email:            "meili.chen@example.com",
				Phone:            "0922333444",
				Position:         "人資專員",
				Department:       "人資部",
				Level:            2,
				Salary:           50000,
				HireDate:         time.Date(2021, 7, 1, 0, 0, 0, 0, time.Local),
				Address:          "新北市板橋區",
				EmergencyContact: "陳爸爸 0933444555",
				Status:           "active",
			},
		}
		DB.Create(&employees)

		leaves := []models.Leave{
			{
				EmployeeID: 1,
				StartDate:  time.Date(2024, 6, 1, 0, 0, 0, 0, time.Local),
				EndDate:    time.Date(2024, 6, 3, 0, 0, 0, 0, time.Local),
				LeaveType:  "年假",
				Reason:     "家庭旅遊",
				Status:     "approved",
			},
			{
				EmployeeID: 2,
				StartDate:  time.Date(2024, 6, 5, 0, 0, 0, 0, time.Local),
				EndDate:    time.Date(2024, 6, 5, 0, 0, 0, 0, time.Local),
				LeaveType:  "病假",
				Reason:     "感冒請假",
				Status:     "pending",
			},
		}
		DB.Create(&leaves)
	}
	// ====== SEED DATA END ======

	log.Println("Successfully connected to database and migrated schemas")
}
