package redis

import (
	"context"
	"doo-store/backend/config"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var (
	// 全局Redis客户端
	client *redis.Client
	mu     sync.Mutex
)

// 初始化Redis客户端
func Init() error {
	redisConfig := config.EnvConfig.DooTaskRedis()
	addr := fmt.Sprintf("%s:%s", redisConfig.HOST, redisConfig.PORT)

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisConfig.PASSWORD,
		DB:       0, // 默认使用0号数据库
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Errorf("Redis连接失败: %v", err)
		return err
	}

	log.Info("Redis连接成功")
	return nil
}

// GetClient 获取Redis客户端实例
func GetClient() *redis.Client {
	if client == nil {
		mu.Lock()
		defer mu.Unlock()

		// 双重检查，防止在获取锁的过程中其他goroutine已经完成了初始化
		if client == nil {
			if err := Init(); err != nil {
				log.Errorf("Redis初始化失败: %v", err)
				return nil
			}
		}
	}
	return client
}

// Set 设置键值对
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return GetClient().Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(ctx context.Context, key string) (string, error) {
	return GetClient().Get(ctx, key).Result()
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	return GetClient().Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, keys ...string) (bool, error) {
	result, err := GetClient().Exists(ctx, keys...).Result()
	return result > 0, err
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return GetClient().Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return GetClient().TTL(ctx, key).Result()
}

// HSet 设置哈希表字段
func HSet(ctx context.Context, key string, values ...interface{}) error {
	return GetClient().HSet(ctx, key, values...).Err()
}

// HGet 获取哈希表字段
func HGet(ctx context.Context, key, field string) (string, error) {
	return GetClient().HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希表所有字段
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return GetClient().HGetAll(ctx, key).Result()
}

// HDel 删除哈希表字段
func HDel(ctx context.Context, key string, fields ...string) error {
	return GetClient().HDel(ctx, key, fields...).Err()
}

// LPush 将元素推入列表左侧
func LPush(ctx context.Context, key string, values ...interface{}) error {
	return GetClient().LPush(ctx, key, values...).Err()
}

// RPush 将元素推入列表右侧
func RPush(ctx context.Context, key string, values ...interface{}) error {
	return GetClient().RPush(ctx, key, values...).Err()
}

// LPop 从列表左侧弹出元素
func LPop(ctx context.Context, key string) (string, error) {
	return GetClient().LPop(ctx, key).Result()
}

// RPop 从列表右侧弹出元素
func RPop(ctx context.Context, key string) (string, error) {
	return GetClient().RPop(ctx, key).Result()
}

// LRange 获取列表指定范围的元素
func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetClient().LRange(ctx, key, start, stop).Result()
}

// SAdd 添加集合成员
func SAdd(ctx context.Context, key string, members ...interface{}) error {
	return GetClient().SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func SMembers(ctx context.Context, key string) ([]string, error) {
	return GetClient().SMembers(ctx, key).Result()
}

// SRem 移除集合成员
func SRem(ctx context.Context, key string, members ...interface{}) error {
	return GetClient().SRem(ctx, key, members...).Err()
}

// ZAdd 添加有序集合成员
func ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return GetClient().ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合指定范围的成员
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return GetClient().ZRange(ctx, key, start, stop).Result()
}

// ZRem 移除有序集合成员
func ZRem(ctx context.Context, key string, members ...interface{}) error {
	return GetClient().ZRem(ctx, key, members...).Err()
}

// Incr 将键的整数值加1
func Incr(ctx context.Context, key string) (int64, error) {
	return GetClient().Incr(ctx, key).Result()
}

// IncrBy 将键的整数值加上指定的增量
func IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return GetClient().IncrBy(ctx, key, value).Result()
}

// Decr 将键的整数值减1
func Decr(ctx context.Context, key string) (int64, error) {
	return GetClient().Decr(ctx, key).Result()
}

// DecrBy 将键的整数值减去指定的减量
func DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return GetClient().DecrBy(ctx, key, value).Result()
}

// 设置带前缀的键名
func KeyWithPrefix(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}
