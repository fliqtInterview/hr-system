package handlers

import (
	"net/http"
	"strconv"

	"hr-system/internal/models"

	"github.com/gin-gonic/gin"
)

// LeaveServiceInterface 定義請假服務接口
type LeaveServiceInterface interface {
	CreateLeave(leave *models.Leave) error
	GetLeave(id uint) (*models.Leave, error)
	ListLeaves() ([]models.Leave, error)
	UpdateLeaveStatus(id uint, status string, remark string) error
	DeleteLeave(id uint) error
}

type LeaveHandler struct {
	leaveService LeaveServiceInterface
}

func NewLeaveHandler(leaveService LeaveServiceInterface) *LeaveHandler {
	return &LeaveHandler{
		leaveService: leaveService,
	}
}

// CreateLeave 創建請假記錄
func (h *LeaveHandler) CreateLeave(c *gin.Context) {
	var leave models.Leave
	if err := c.ShouldBindJSON(&leave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 驗證必填字段
	if leave.EmployeeID == 0 || leave.LeaveType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee ID and leave type are required"})
		return
	}

	if err := h.leaveService.CreateLeave(&leave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, leave)
}

// GetLeave 獲取請假記錄
func (h *LeaveHandler) GetLeave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	leave, err := h.leaveService.GetLeave(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave record not found"})
		return
	}

	c.JSON(http.StatusOK, leave)
}

// ListLeaves 獲取請假記錄列表
func (h *LeaveHandler) ListLeaves(c *gin.Context) {
	leaves, err := h.leaveService.ListLeaves()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if leaves == nil {
		leaves = []models.Leave{}
	}
	c.JSON(http.StatusOK, leaves)
}

// UpdateLeaveStatus 更新請假狀態
func (h *LeaveHandler) UpdateLeaveStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var status struct {
		Status string `json:"status"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.leaveService.UpdateLeaveStatus(uint(id), status.Status, status.Remark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leave status updated successfully"})
}

// DeleteLeave 刪除請假記錄
func (h *LeaveHandler) DeleteLeave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.leaveService.DeleteLeave(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Leave record deleted successfully"})
}
