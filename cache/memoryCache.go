package cache

import (
	"errors"
	"sync"
)

type MemoryCache struct {
	store     map[string]string
	mu        sync.RWMutex
	maxLength int
}

func InitMemoryCache(maxLength int) (*MemoryCache, error) {
	if maxLength <= 0 {
		return nil, errors.New("maxLength must be greater than 0")
	}
	return &MemoryCache{
		store:     make(map[string]string),
		maxLength: maxLength,
	}, nil
}

func (m *MemoryCache) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, found := m.store[key]
	return val, found
}

func (m *MemoryCache) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果缓存超出最大长度，移除最旧的元素
	if len(m.store) >= m.maxLength {
		for k := range m.store {
			delete(m.store, k)
			break
		}
	}

	// 设置新的值
	m.store[key] = value
}
