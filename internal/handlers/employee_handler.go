package handlers

import (
	"net/http"
	"strconv"

	"hr-system/internal/models"

	"github.com/gin-gonic/gin"
)

// EmployeeServiceInterface 定義員工服務接口
type EmployeeServiceInterface interface {
	CreateEmployee(employee *models.Employee) error
	GetEmployee(id uint) (*models.Employee, error)
	ListEmployees() ([]models.Employee, error)
	UpdateEmployee(employee *models.Employee) error
	DeleteEmployee(id uint) error
}

type EmployeeHandler struct {
	employeeService EmployeeServiceInterface
}

func NewEmployeeHandler(employeeService EmployeeServiceInterface) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

// CreateEmployee 創建員工
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 驗證必填字段
	if employee.Name == "" || employee.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and email are required"})
		return
	}

	if err := h.employeeService.CreateEmployee(&employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

// GetEmployee 獲取員工信息
func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	employee, err := h.employeeService.GetEmployee(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// ListEmployees 獲取員工列表
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	employees, err := h.employeeService.ListEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if employees == nil {
		employees = []models.Employee{}
	}
	c.JSON(http.StatusOK, employees)
}

// UpdateEmployee 更新員工信息
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee.ID = uint(id)
	if err := h.employeeService.UpdateEmployee(&employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// DeleteEmployee 刪除員工
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.employeeService.DeleteEmployee(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
