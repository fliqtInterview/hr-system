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

// MockEmployeeService 模擬員工服務
type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) CreateEmployee(employee *models.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockEmployeeService) GetEmployee(id uint) (*models.Employee, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeService) ListEmployees() ([]models.Employee, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	if args.Get(0) == nil {
		return []models.Employee{}, nil
	}
	return args.Get(0).([]models.Employee), nil
}

func (m *MockEmployeeService) UpdateEmployee(employee *models.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockEmployeeService) DeleteEmployee(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// 確保 MockEmployeeService 實現了 EmployeeServiceInterface
var _ EmployeeServiceInterface = (*MockEmployeeService)(nil)

func setupTestRouter(handler *EmployeeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	api := r.Group("/api")
	{
		employees := api.Group("/employees")
		{
			employees.POST("", handler.CreateEmployee)
			employees.GET("", handler.ListEmployees)
			employees.GET("/:id", handler.GetEmployee)
			employees.PUT("/:id", handler.UpdateEmployee)
			employees.DELETE("/:id", handler.DeleteEmployee)
		}
	}

	return r
}

func TestCreateEmployee(t *testing.T) {
	mockService := &MockEmployeeService{}
	handler := NewEmployeeHandler(mockService)
	router := setupTestRouter(handler)

	hireDate, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")

	tests := []struct {
		name       string
		payload    models.Employee
		mockSetup  func()
		wantStatus int
	}{
		{
			name: "成功創建員工",
			payload: models.Employee{
				Name:             "測試員工",
				Email:            "test@example.com",
				Phone:            "13800138000",
				Position:         "工程師",
				Department:       "技術部",
				Level:            1,
				Salary:           10000,
				HireDate:         hireDate,
				Address:          "北京市",
				EmergencyContact: "緊急聯繫人",
			},
			mockSetup: func() {
				mockService.On("CreateEmployee", mock.AnythingOfType("*models.Employee")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "無效的員工數據",
			payload: models.Employee{
				Name: "", // 缺少必填字段
			},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/employees", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetEmployee(t *testing.T) {
	mockService := &MockEmployeeService{}
	handler := NewEmployeeHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name       string
		id         string
		mockSetup  func()
		wantStatus int
	}{
		{
			name: "成功獲取員工",
			id:   "1",
			mockSetup: func() {
				employee := &models.Employee{}
				employee.ID = 1
				employee.Name = "測試員工"
				mockService.On("GetEmployee", uint(1)).Return(employee, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "員工不存在",
			id:   "999",
			mockSetup: func() {
				mockService.On("GetEmployee", uint(999)).Return(nil, assert.AnError)
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

			req := httptest.NewRequest(http.MethodGet, "/api/employees/"+tt.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestListEmployees(t *testing.T) {
	mockService := &MockEmployeeService{}
	handler := NewEmployeeHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name       string
		mockSetup  func()
		wantStatus int
	}{
		{
			name: "成功獲取員工列表",
			mockSetup: func() {
				employees := []models.Employee{
					{
						Model: gorm.Model{ID: 1},
						Name:  "員工1",
					},
					{
						Model: gorm.Model{ID: 2},
						Name:  "員工2",
					},
				}
				mockService.On("ListEmployees").Return(employees, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "獲取員工列表失敗",
			mockSetup: func() {
				mockService.On("ListEmployees").Return(nil, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodGet, "/api/employees", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestUpdateEmployee(t *testing.T) {
	mockService := &MockEmployeeService{}
	handler := NewEmployeeHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name       string
		id         string
		payload    models.Employee
		mockSetup  func()
		wantStatus int
	}{
		{
			name: "成功更新員工",
			id:   "1",
			payload: models.Employee{
				Name:     "更新後的員工",
				Email:    "updated@example.com",
				Position: "高級工程師",
			},
			mockSetup: func() {
				mockService.On("UpdateEmployee", mock.AnythingOfType("*models.Employee")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "無效的ID",
			id:         "invalid",
			payload:    models.Employee{},
			mockSetup:  func() {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPut, "/api/employees/"+tt.id, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestDeleteEmployee(t *testing.T) {
	mockService := &MockEmployeeService{}
	handler := NewEmployeeHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name       string
		id         string
		mockSetup  func()
		wantStatus int
	}{
		{
			name: "成功刪除員工",
			id:   "1",
			mockSetup: func() {
				mockService.On("DeleteEmployee", uint(1)).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "員工不存在",
			id:   "999",
			mockSetup: func() {
				mockService.On("DeleteEmployee", uint(999)).Return(assert.AnError)
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

			req := httptest.NewRequest(http.MethodDelete, "/api/employees/"+tt.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
