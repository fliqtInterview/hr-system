# HR System API

本專案提供一個簡單的人事管理系統 API，包含員工與請假管理功能。

## 系統需求

- Go 1.21 或以上
- Docker & Docker Compose
- Make

## 快速開始

### 1. 建置專案

```bash
# 編譯專案
make build

# 執行單元測試
make test
```

### 2. 部署服務

```bash
# 部署所有服務（API、MySQL、Redis）
make deploy

# 檢查服務狀態
make ps

# 查看服務日誌
make logs

# 停止所有服務
make stop
```

### 3. 服務端口

- API 服務：http://localhost:8080
- MySQL：localhost:3306
- Redis：localhost:6379

## API 使用說明

所有 API 都可以使用 curl 或其他 HTTP 客戶端（如 Postman）進行調用。

### 健康檢查

```bash
# 請求
curl http://localhost:8080/ping

# 回應
{
  "message": "pong"
}
```

### 員工管理 API

#### 1. 新增員工

```bash
# 請求
curl -X POST http://localhost:8080/api/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "王小明",
    "email": "xiaoming.wang@example.com",
    "phone": "0912345678",
    "position": "工程師",
    "department": "研發部",
    "level": 1,
    "salary": 60000,
    "hire_date": "2022-01-10T00:00:00Z",
    "address": "台北市信義區",
    "emergency_contact": "王媽媽 0911222333",
    "status": "active"
  }'

# 回應
{
  "id": 1,
  "created_at": "2024-05-06T11:43:23Z",
  "updated_at": "2024-05-06T11:43:23Z",
  "deleted_at": null,
  "name": "王小明",
  "email": "xiaoming.wang@example.com",
  "phone": "0912345678",
  "position": "工程師",
  "department": "研發部",
  "level": 1,
  "salary": 60000,
  "hire_date": "2022-01-10T00:00:00Z",
  "address": "台北市信義區",
  "emergency_contact": "王媽媽 0911222333",
  "status": "active"
}
```

#### 2. 查詢員工列表

```bash
# 請求
curl http://localhost:8080/api/employees

# 回應
[
  {
    "id": 1,
    "created_at": "2024-05-06T11:43:23Z",
    "updated_at": "2024-05-06T11:43:23Z",
    "deleted_at": null,
    "name": "王小明",
    "email": "xiaoming.wang@example.com",
    "phone": "0912345678",
    "position": "工程師",
    "department": "研發部",
    "level": 1,
    "salary": 60000,
    "hire_date": "2022-01-10T00:00:00Z",
    "address": "台北市信義區",
    "emergency_contact": "王媽媽 0911222333",
    "status": "active"
  }
]
```

#### 3. 查詢單一員工

```bash
# 請求
curl http://localhost:8080/api/employees/1

# 回應
{
  "id": 1,
  "created_at": "2024-05-06T11:43:23Z",
  "updated_at": "2024-05-06T11:43:23Z",
  "deleted_at": null,
  "name": "王小明",
  "email": "xiaoming.wang@example.com",
  "phone": "0912345678",
  "position": "工程師",
  "department": "研發部",
  "level": 1,
  "salary": 60000,
  "hire_date": "2022-01-10T00:00:00Z",
  "address": "台北市信義區",
  "emergency_contact": "王媽媽 0911222333",
  "status": "active"
}
```

#### 4. 更新員工

```bash
# 請求
curl -X PUT http://localhost:8080/api/employees/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "王小明",
    "email": "xiaoming.wang@example.com",
    "phone": "0912345678",
    "position": "資深工程師",
    "department": "研發部",
    "level": 2,
    "salary": 70000,
    "hire_date": "2022-01-10T00:00:00Z",
    "address": "台北市信義區",
    "emergency_contact": "王媽媽 0911222333",
    "status": "active"
  }'

# 回應
{
  "id": 1,
  "created_at": "2024-05-06T11:43:23Z",
  "updated_at": "2024-05-06T11:43:23Z",
  "deleted_at": null,
  "name": "王小明",
  "email": "xiaoming.wang@example.com",
  "phone": "0912345678",
  "position": "資深工程師",
  "department": "研發部",
  "level": 2,
  "salary": 70000,
  "hire_date": "2022-01-10T00:00:00Z",
  "address": "台北市信義區",
  "emergency_contact": "王媽媽 0911222333",
  "status": "active"
}
```

#### 5. 刪除員工

```bash
# 請求
curl -X DELETE http://localhost:8080/api/employees/1

# 回應
{
  "message": "Employee deleted successfully"
}
```

### 請假管理 API

#### 1. 新增請假

```bash
# 請求
curl -X POST http://localhost:8080/api/leaves \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "start_date": "2024-06-01T00:00:00Z",
    "end_date": "2024-06-03T00:00:00Z",
    "leave_type": "年假",
    "reason": "家庭旅遊",
    "status": "pending"
  }'

# 回應
{
  "id": 1,
  "created_at": "2024-05-06T11:43:23Z",
  "updated_at": "2024-05-06T11:43:23Z",
  "deleted_at": null,
  "employee_id": 1,
  "start_date": "2024-06-01T00:00:00Z",
  "end_date": "2024-06-03T00:00:00Z",
  "leave_type": "年假",
  "reason": "家庭旅遊",
  "status": "pending"
}
```

#### 2. 查詢請假列表

```bash
# 請求
curl http://localhost:8080/api/leaves

# 回應
[
  {
    "id": 1,
    "created_at": "2024-05-06T11:43:23Z",
    "updated_at": "2024-05-06T11:43:23Z",
    "deleted_at": null,
    "employee_id": 1,
    "start_date": "2024-06-01T00:00:00Z",
    "end_date": "2024-06-03T00:00:00Z",
    "leave_type": "年假",
    "reason": "家庭旅遊",
    "status": "pending"
  }
]
```

#### 3. 查詢單一請假記錄

```bash
# 請求
curl http://localhost:8080/api/leaves/1

# 回應
{
  "id": 1,
  "created_at": "2024-05-06T11:43:23Z",
  "updated_at": "2024-05-06T11:43:23Z",
  "deleted_at": null,
  "employee_id": 1,
  "start_date": "2024-06-01T00:00:00Z",
  "end_date": "2024-06-03T00:00:00Z",
  "leave_type": "年假",
  "reason": "家庭旅遊",
  "status": "pending"
}
```

#### 4. 更新請假狀態

```bash
# 請求
curl -X PUT http://localhost:8080/api/leaves/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved",
    "remark": "已核准"
  }'

# 回應
{
  "message": "Leave status updated successfully"
}
```

#### 5. 刪除請假記錄

```bash
# 請求
curl -X DELETE http://localhost:8080/api/leaves/1

# 回應
{
  "message": "Leave record deleted successfully"
}
```

## 資料結構

### 員工（Employee）

```json
{
  "id": "整數，自動生成",
  "name": "字串，必填，員工姓名",
  "email": "字串，必填，員工郵箱（唯一）",
  "phone": "字串，員工電話",
  "position": "字串，職位",
  "department": "字串，部門",
  "level": "整數，職等",
  "salary": "浮點數，薪資",
  "hire_date": "日期時間，入職日期",
  "address": "字串，地址",
  "emergency_contact": "字串，緊急聯絡人",
  "status": "字串，狀態（active/inactive）"
}
```

### 請假（Leave）

```json
{
  "id": "整數，自動生成",
  "employee_id": "整數，必填，關聯員工ID",
  "start_date": "日期時間，必填，開始日期",
  "end_date": "日期時間，必填，結束日期",
  "leave_type": "字串，必填，請假類型（年假/病假/事假等）",
  "reason": "字串，請假原因",
  "status": "字串，狀態（pending/approved/rejected）",
  "approver_id": "整數，審批人ID",
  "approve_time": "日期時間，審批時間",
  "approve_remark": "字串，審批備註"
}
```

## 錯誤處理

所有 API 在發生錯誤時會返回適當的 HTTP 狀態碼和錯誤訊息：

- 400 Bad Request：請求格式錯誤
- 404 Not Found：資源不存在
- 500 Internal Server Error：服務器內部錯誤

錯誤回應格式：
```json
{
  "error": "錯誤訊息"
}
``` 