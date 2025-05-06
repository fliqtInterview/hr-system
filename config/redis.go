package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// 緩存過期時間
const (
	EmployeeCacheExpiration = 30 * time.Minute
	LeaveCacheExpiration    = 15 * time.Minute
)

// 緩存鍵前綴
const (
	EmployeeKeyPrefix = "employee:"
	LeaveKeyPrefix    = "leave:"
)

// InitRedis 初始化 Redis 連接
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
		Password: getEnv("REDIS_PASSWORD", ""), // 如果有密碼的話
		DB:       0,                            // 使用默認的 DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 測試連接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Successfully connected to Redis")
}
