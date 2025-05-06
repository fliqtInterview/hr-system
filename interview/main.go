package main

import (
	"log"

	"hr-system/config"
	"hr-system/internal/handlers"
	"hr-system/internal/models"
	"hr-system/internal/repositories"
	"hr-system/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化數據庫連接
	config.InitDB()

	// 分步驟遷移數據庫結構
	if err := config.DB.AutoMigrate(&models.Employee{}); err != nil {
		log.Fatal("Failed to migrate employee table:", err)
	}
	if err := config.DB.AutoMigrate(&models.Leave{}); err != nil {
		log.Fatal("Failed to migrate leave table:", err)
	}

	// 初始化依賴
	employeeRepo := repositories.NewEmployeeRepository(config.DB)
	leaveRepo := repositories.NewLeaveRepository(config.DB)

	employeeService := services.NewEmployeeService(employeeRepo)
	leaveService := services.NewLeaveService(leaveRepo, employeeRepo)

	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	leaveHandler := handlers.NewLeaveHandler(leaveService)

	// 設置路由
	r := gin.Default()

	// 健康檢查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API 路由組
	api := r.Group("/api")
	{
		// 員工相關路由
		employees := api.Group("/employees")
		{
			employees.POST("", employeeHandler.CreateEmployee)
			employees.GET("", employeeHandler.ListEmployees)
			employees.GET("/:id", employeeHandler.GetEmployee)
			employees.PUT("/:id", employeeHandler.UpdateEmployee)
			employees.DELETE("/:id", employeeHandler.DeleteEmployee)
		}

		// 請假相關路由
		leaves := api.Group("/leaves")
		{
			leaves.POST("", leaveHandler.CreateLeave)
			leaves.GET("", leaveHandler.ListLeaves)
			leaves.GET("/:id", leaveHandler.GetLeave)
			leaves.PUT("/:id/status", leaveHandler.UpdateLeaveStatus)
			leaves.DELETE("/:id", leaveHandler.DeleteLeave)
		}
	}

	// 啟動服務器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
