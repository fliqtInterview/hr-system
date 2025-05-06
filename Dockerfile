FROM golang:1.21-alpine

WORKDIR /app

# 安裝依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製源代碼
COPY . .

# 編譯
RUN go build -o main .

# 暴露端口
EXPOSE 8080

# 運行
CMD ["./main"] 