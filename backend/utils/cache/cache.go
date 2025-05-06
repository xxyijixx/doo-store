package cache

import (
	"context"
	"doo-store/backend/utils/redis"
	"encoding/json"
	"time"
	
	log "github.com/sirupsen/logrus"
)

// 默认缓存过期时间
const DefaultExpiration = time.Hour * 24

// 设置缓存，自动序列化对象为JSON
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redis.Set(ctx, key, string(data), expiration)
}

// 获取缓存，自动反序列化JSON为对象
func Get(ctx context.Context, key string, dest interface{}) error {
	data, err := redis.Get(ctx, key)
	if err != nil {
		return err
	}
	
	return json.Unmarshal([]byte(data), dest)
}

// 获取或设置缓存
// 如果缓存存在，则返回缓存的值
// 如果缓存不存在，则调用 getter 函数获取值，并将其设置到缓存中
func GetOrSet(ctx context.Context, key string, dest interface{}, expiration time.Duration, getter func() (interface{}, error)) error {
	// 尝试从缓存获取
	err := Get(ctx, key, dest)
	if err == nil {
		// 缓存命中
		return nil
	}
	
	// 缓存未命中，调用getter获取数据
	value, err := getter()
	if err != nil {
		return err
	}
	
	// 设置缓存
	if err := Set(ctx, key, value, expiration); err != nil {
		log.Warnf("设置缓存失败: %v", err)
	}
	
	// 将获取的值赋给dest
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, dest)
}

// 删除缓存
func Delete(ctx context.Context, keys ...string) error {
	return redis.Del(ctx, keys...)
}