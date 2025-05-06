package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hr-system/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockLeaveService 模擬請假服務
type MockLeaveService struct {
	mock.Mock
}

func (m *MockLeaveService) CreateLeave(leave *models.Leave) error {
	args := m.Called(leave)
	return args.Error(0)
}

func (m *MockLeaveService) GetLeave(id uint) (*models.Leave, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Leave), args.Error(1)
}

func (m *MockLeaveService) ListLeaves() ([]models.Leave, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	if args.Get(0) == nil {
		return []models.Leave{}, nil
	}
	return args.Get(0).([]models.Leave), nil
}

func (m *MockLeaveService) UpdateLeaveStatus(id uint, status string, remark string) error {
	args := m.Called(id, status, remark)
	return args.Error(0)
}

func (m *MockLeaveService) DeleteLeave(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// 確保 MockLeaveService 實現了 LeaveServiceInterface
var _ LeaveServiceInterface = (*MockLeaveService)(nil)

func setupLeaveTestRouter(handler *LeaveHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	api := r.Group("/api")
	{
		leaves := api.Group("/leaves")
		{
			leaves.POST("", handler.CreateLeave)
			leaves.GET("", handler.ListLeaves)
			leaves.GET("/:id", handler.GetLeave)
			leaves.PUT("/:id/status", handler.UpdateLeaveStatus)
			leaves.DELETE("/:id", handler.DeleteLeave)
		}
	}

	return r
}

func TestCreateLeave(t *testing.T) {
	mockService := &MockLeaveService{}
	handler := NewLeaveHandler(mockService)
	router := setupLeaveTestRouter(handler)

	startDate, _ := time.Parse(time.RFC3339, "2024-04-01T00:00:00Z")
	endDate, _ := time.Parse(time.RFC3339, "2024-04-02T00:00:00Z")

	tests := []struct {
		name       string
		payload    models.Leave
		mockSetup  func()
		wantStatus int
	}{
		{
			name: "成功創建請假記錄",
			payload: models.Leave{
				EmployeeID: 1,
				StartDate:  startDate,
				EndDate:    endDate,
				LeaveType:  "年假",
				Reason:     "休息",
			},
			mockSetup: func() {
				mockService.On("CreateLeave", mock.AnythingOfType("*models.Leave")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "無效的請假數據",
			payload: models.Leave{
				EmployeeID: 0, // 缺少必填字段
			},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/leaves", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetLeave(t *testing.T) {
	mockService := &MockLeaveService{}
	handler := NewLeaveHandler(mockService)
	router := setupLeaveTestRouter(handler)

	startDate, _ := time.Parse(time.RFC3339, "2024-04-01T00:00:00Z")
	endDate, _ := time.Parse(time.RFC3339, "2024-04-02T00:00:00Z")

	tests := []struct {
		name       string
		id         string
		mockSetup  func()
		wantStatus int
	}{
		{
			name: "成功獲取請假記錄",
			id:   "1",
			mockSetup: func() {
				leave := &models.Leave{}
				leave.ID = 1
				leave.EmployeeID = 1
				leave.StartDate = startDate
				leave.EndDate = endDate
				mockService.On("GetLeave", uint(1)).Return(leave, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "請假記錄不存在",
			id:   "999",
			mockSetup: func() {
				mockService.On("GetLeave", uint(999)).Return(nil, assert.AnError)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "無效的ID",
			id:         "invalid",
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/leaves/"+tt.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestListLeaves(t *testing.T) {
	mockService := &MockLeaveService{}
	handler := NewLeaveHandler(mockService)
	router := setupLeaveTestRouter(handler)

	tests := []struct {
		name       string
		mockSetup  func()
		wantStatus int
	}{
		{
			name: "成功獲取請假列表",
			mockSetup: func() {
				leaves := []models.Leave{
					{
						Model:      gorm.Model{ID: 1},
						EmployeeID: 1,
						StartDate:  time.Now(),
						EndDate:    time.Now().Add(24 * time.Hour),
					},
					{
						Model:      gorm.Model{ID: 2},
						EmployeeID: 2,
						StartDate:  time.Now(),
						EndDate:    time.Now().Add(48 * time.Hour),
					},
				}
				mockService.On("ListLeaves").Return(leaves, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "獲取請假列表失敗",
			mockSetup: func() {
				mockService.On("ListLeaves").Return(nil, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/leaves", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestUpdateLeaveStatus(t *testing.T) {
	mockService := &MockLeaveService{}
	handler := NewLeaveHandler(mockService)
	router := setupLeaveTestRouter(handler)

	tests := []struct {
		name       string
		id         string
		payload    map[string]string
		mockSetup  func()
		wantStatus int
	}{
		{
			name: "成功更新請假狀態",
			id:   "1",
			payload: map[string]string{
				"status": "approved",
				"remark": "同意",
			},
			mockSetup: func() {
				mockService.On("UpdateLeaveStatus", uint(1), "approved", "同意").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "無效的ID",
			id:         "invalid",
			payload:    map[string]string{},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPut, "/api/leaves/"+tt.id+"/status", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestDeleteLeave(t *testing.T) {
	mockService := &MockLeaveService{}
	handler := NewLeaveHandler(mockService)
	router := setupLeaveTestRouter(handler)

	tests := []struct {
		name       string
		id         string
		mockSetup  func()
		wantStatus int
	}{
		{
			name: "成功刪除請假記錄",
			id:   "1",
			mockSetup: func() {
				mockService.On("DeleteLeave", uint(1)).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "請假記錄不存在",
			id:   "999",
			mockSetup: func() {
				mockService.On("DeleteLeave", uint(999)).Return(assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "無效的ID",
			id:         "invalid",
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodDelete, "/api/leaves/"+tt.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
