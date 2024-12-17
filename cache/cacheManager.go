package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/muyi-zcy/my-golru/fetcher"
	"log"
)

type CacheManager struct {
	memCache    *MemoryCache
	redisCache  *RedisCache
	dataFetcher fetcher.DataFetcher
}

func InitCacheManager(code string, memCacheSize int, redisClient *redis.Client, redisCacheSize int, dataFetcher fetcher.DataFetcher) *CacheManager {
	// 创建内存缓存
	memCache, err := InitMemoryCache(memCacheSize)
	if err != nil {
		log.Fatalf("Error creating memory cache: %v", err)
	}

	// 创建 Redis 缓存（如果传入了 Redis 客户端）
	var redisCache *RedisCache
	if redisClient != nil {
		redisCache = NewRedisCache(code, redisClient, redisCacheSize)
	}

	return &CacheManager{
		memCache:    memCache,
		redisCache:  redisCache,
		dataFetcher: dataFetcher,
	}
}

func (cm *CacheManager) GetData(key string) (string, error) {
	// 1. 从内存缓存获取
	if val, found := cm.memCache.Get(key); found {
		return val, nil
	}

	// 2. 从 Redis 获取（如果启用了 Redis）
	if cm.redisCache != nil {
		if val, found := cm.redisCache.Get(key); found {
			cm.memCache.Set(key, val)
			return val, nil
		}
	}

	// 3. 从数据库获取
	val, err := cm.dataFetcher.FetchData(key)
	if err != nil {
		return "", err
	}

	// 4. 更新缓存
	cm.memCache.Set(key, val)
	if cm.redisCache != nil {
		cm.redisCache.Set(key, val)
	}

	return val, nil
}
