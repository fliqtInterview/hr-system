package handlers

import (
	"net/http"
	"strconv"

	"hr-system/internal/models"
	"hr-system/internal/services"

	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	service *services.LeaveService
}

func NewLeaveHandler(service *services.LeaveService) *LeaveHandler {
	return &LeaveHandler{service: service}
}

// CreateLeave 創建請假申請
func (h *LeaveHandler) CreateLeave(c *gin.Context) {
	var leave models.Leave
	if err := c.ShouldBindJSON(&leave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateLeave(&leave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, leave)
}

// GetLeave 獲取請假記錄
func (h *LeaveHandler) GetLeave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	leave, err := h.service.GetLeave(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "leave record not found"})
		return
	}

	c.JSON(http.StatusOK, leave)
}

// UpdateLeaveStatus 更新請假狀態
func (h *LeaveHandler) UpdateLeaveStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var request struct {
		Status       models.LeaveStatus `json:"status"`
		ApproverID   uint               `json:"approver_id"`
		RejectReason string             `json:"reject_reason"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateLeaveStatus(uint(id), request.Status, request.ApproverID, request.RejectReason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "leave status updated successfully"})
}

// ListLeaves 獲取請假記錄列表
func (h *LeaveHandler) ListLeaves(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var employeeID *uint
	if empID := c.Query("employee_id"); empID != "" {
		id, err := strconv.ParseUint(empID, 10, 32)
		if err == nil {
			uintID := uint(id)
			employeeID = &uintID
		}
	}

	leaves, total, err := h.service.ListLeaves(page, pageSize, employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  leaves,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// CancelLeave 取消請假申請
func (h *LeaveHandler) CancelLeave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	employeeID, err := strconv.ParseUint(c.GetHeader("X-Employee-ID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	if err := h.service.CancelLeave(uint(id), uint(employeeID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "leave cancelled successfully"})
}

// DeleteLeave 刪除請假記錄
func (h *LeaveHandler) DeleteLeave(c *gin.Context) {
	id := c.Param("id")
	leaveID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid leave ID"})
		return
	}

	err = h.service.DeleteLeave(uint(leaveID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "leave deleted successfully"})
}
