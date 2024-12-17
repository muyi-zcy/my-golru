package cache

import (
	"github.com/go-redis/redis/v7"
	"log"
)

type RedisCache struct {
	code      string
	client    *redis.Client
	maxLength int // 最大缓存数量
}

func NewRedisCache(code string, client *redis.Client, maxLength int) *RedisCache {
	if client == nil {
		return nil // 返回 nil 表示不启用 Redis 二级缓存
	}
	return &RedisCache{
		code:      code,
		client:    client,
		maxLength: maxLength,
	}
}

func (r *RedisCache) Get(key string) (string, bool) {
	if r == nil {
		return "", false // 如果 Redis 未启用，直接返回
	}

	// 如果内存缓存没有，再从 Redis 获取
	val, err := r.client.HGet(r.code, key).Result()
	if err == redis.Nil {
		return "", false
	}
	if err != nil {
		log.Println("Redis error:", err)
		return "", false
	}
	return val, true
}
func (r *RedisCache) Set(key, value string) {
	if r == nil {
		return // 如果 Redis 未启用，不进行任何操作
	}

	// 检查 Redis 中存储的缓存数量
	cacheLength, err := r.client.HLen(r.code).Result()
	if err != nil {
		log.Printf("Error checking Redis cache length: %v", err)
		return
	}

	// 如果 Redis 缓存已满，移除最旧的元素
	if cacheLength >= int64(r.maxLength) {
		// 获取最旧的缓存 key
		oldestKey, err := r.client.LPop("cache_order").Result() // 假设我们使用一个链表来存储访问顺序
		if err != nil {
			log.Printf("Error getting the oldest key: %v", err)
			return
		}

		// 从 Redis 哈希表中删除该 key
		r.client.HDel("cache_map", oldestKey)
	}

	// 将新数据存入 Redis 哈希表
	r.client.HSet(r.code, key, value)

	// 更新顺序链表：将新键添加到链表的末尾
	r.client.RPush(r.code+"_order", key)
}
