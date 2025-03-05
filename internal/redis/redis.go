package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	client *redis.Client
	logger *zap.Logger
)

// Config Redis配置结构
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// InitRedis 初始化Redis客户端
func InitRedis(cfg *Config) error {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to create logger: %v", err)
	}
	
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info("Creating Redis client",
		zap.String("addr", addr),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Int("db", cfg.DB))

	opts := &redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			logger.Info("Redis OnConnect callback triggered",
				zap.String("addr", addr))
			return nil
		},
	}

	client = redis.NewClient(opts)
	logger.Info("Redis client created, attempting to ping")

	// 测试连接，使用带超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to ping Redis",
			zap.Error(err),
			zap.String("addr", addr))
		return fmt.Errorf("failed to ping redis: %v", err)
	}

	// 尝试设置一个测试值
	testKey := "test_connection"
	testValue := "ok"
	err = client.Set(ctx, testKey, testValue, 1*time.Minute).Err()
	if err != nil {
		logger.Error("Failed to set test value",
			zap.Error(err),
			zap.String("key", testKey))
		return fmt.Errorf("failed to set test value: %v", err)
	}

	// 尝试获取测试值
	val, err := client.Get(ctx, testKey).Result()
	if err != nil {
		logger.Error("Failed to get test value",
			zap.Error(err),
			zap.String("key", testKey))
		return fmt.Errorf("failed to get test value: %v", err)
	}

	if val != testValue {
		logger.Error("Test value mismatch",
			zap.String("expected", testValue),
			zap.String("got", val))
		return fmt.Errorf("test value mismatch: expected %s, got %s", testValue, val)
	}

	logger.Info("Redis connection test successful",
		zap.String("addr", addr))
	return nil
}

// GetClient 获取Redis客户端实例
func GetClient() *redis.Client {
	if client == nil {
		logger.Error("Redis client is nil")
		return nil
	}
	return client
}

// Set 设置缓存
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		logger.Error("Failed to set cache",
			zap.Error(err),
			zap.String("key", key))
		return err
	}

	logger.Info("Cache set successfully",
		zap.String("key", key),
		zap.Int("data_size", len(data)))
	return nil
}

// Get 获取缓存
func Get(ctx context.Context, key string, value interface{}) error {
	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, value)
}

// Delete 删除缓存
func Delete(ctx context.Context, key string) error {
	return client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	n, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// SetNX 设置缓存（如果不存在）
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	return client.SetNX(ctx, key, data, expiration).Result()
}

// GetOrSet 获取缓存，如果不存在则设置
func GetOrSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	exists, err := Exists(ctx, key)
	if err != nil {
		return err
	}

	if !exists {
		return Set(ctx, key, value, expiration)
	}

	return Get(ctx, key, value)
} 