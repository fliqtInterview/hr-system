.PHONY: build run test clean deploy stop logs ps

# 預設目標
all: build

# 確保 bin 目錄存在
bin:
	mkdir -p bin

# 編譯專案
build: bin
	go build -o bin/hr-system

# 執行專案
run:
	go run main.go

# 執行測試
test:
	go test ./... -v

# 清理編譯產物
clean:
	rm -rf bin/

# 部署（使用 docker-compose）
deploy:
	docker-compose build
	docker-compose up -d

# 停止部署
stop:
	docker-compose down

# 查看服務日誌
logs:
	docker-compose logs -f

# 查看服務狀態
ps:
	docker-compose ps 