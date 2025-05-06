-- 創建數據庫
CREATE DATABASE IF NOT EXISTS hr_system;
USE hr_system;

-- 創建員工表
CREATE TABLE IF NOT EXISTS employees (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6),
    updated_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at TIMESTAMP(6) NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    position VARCHAR(50),
    department VARCHAR(50),
    level INT,
    salary DOUBLE,
    hire_date TIMESTAMP(6),
    address VARCHAR(200),
    emergency_contact VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    UNIQUE KEY idx_employees_email (email)
);

-- 創建請假表
CREATE TABLE IF NOT EXISTS leaves (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6),
    updated_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at TIMESTAMP(6) NULL,
    employee_id BIGINT UNSIGNED NOT NULL,
    type VARCHAR(50) NOT NULL,
    start_date TIMESTAMP(6) NOT NULL,
    end_date TIMESTAMP(6) NOT NULL,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    FOREIGN KEY (employee_id) REFERENCES employees(id)
); 