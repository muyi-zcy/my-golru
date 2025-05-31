# My-GoLRU 

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.22-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

一个高性能的多级缓存系统，用Go语言实现，支持**内存缓存 -> Redis缓存 -> 数据库**的三级缓存架构。

## 🚀 特性

- **多级缓存架构**: 内存缓存(L1) -> Redis缓存(L2) -> 数据库(L3)
- **高性能**: 基于LRU算法的内存缓存，最大化命中率
- **灵活配置**: 支持仅内存缓存或内存+Redis双级缓存
- **类型安全**: 接口化设计，支持多种数据源适配器
- **生产就绪**: 完善的错误处理和日志记录

## 📦 安装

```bash
go get github.com/muyi-zcy/my-golru
```

## 🏗️ 架构设计

```
┌─────────────────┐    Cache Miss    ┌──────────────────┐    Cache Miss    ┌──────────────────┐
│   Memory Cache  │ ────────────────▶ │   Redis Cache    │ ────────────────▶ │    Database      │
│      (L1)       │                  │       (L2)       │                  │       (L3)       │
└─────────────────┘                  └──────────────────┘                  └──────────────────┘
        ▲                                       ▲                                       ▲
        │                                       │                                       │
        └─── Cache Hit (最快) ───────────────────┴─── Cache Hit (快) ────────────────────┘
```

### 核心组件

- **CacheManager**: 缓存管理器，协调多级缓存的访问
- **MemoryCache**: 基于LRU算法的内存缓存
- **RedisCache**: 基于Redis的分布式缓存
- **DataFetcher**: 数据获取器接口，支持多种数据源

## 🔧 快速开始

### 基本使用

```go
package main

import (
    "log"
    "github.com/go-redis/redis/v7"
    "github.com/muyi-zcy/my-golru/cache"
)

// 实现数据获取器接口
type MyDataFetcher struct {}

func (f *MyDataFetcher) FetchData(key string) (string, error) {
    // 从数据库或其他数据源获取数据
    return "data_for_" + key, nil
}

func main() {
    // 创建Redis客户端
    redisClient := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    // 创建数据获取器
    dataFetcher := &MyDataFetcher{}

    // 创建缓存管理器
    cacheManager := cache.InitCacheManager(
        "myapp",      // 缓存前缀
        100,          // 内存缓存大小
        redisClient,  // Redis客户端
        1000,         // Redis缓存大小
        dataFetcher,  // 数据获取器
    )

    // 获取数据
    data, err := cacheManager.GetData("user:123")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Data:", data)
}
```

### 仅使用内存缓存

```go
// 创建仅内存缓存的管理器
cacheManager := cache.InitCacheManager(
    "myapp",     // 缓存前缀
    100,         // 内存缓存大小
    nil,         // 不使用Redis
    0,           // Redis缓存大小设为0
    dataFetcher, // 数据获取器
)
```

### MySQL数据源示例

```go
type MySQLFetcher struct {
    db *sql.DB
}

func (f *MySQLFetcher) FetchData(key string) (string, error) {
    var result string
    query := "SELECT data FROM cache_table WHERE key = ?"
    err := f.db.QueryRow(query, key).Scan(&result)
    if err != nil {
        return "", fmt.Errorf("error fetching from MySQL: %v", err)
    }
    return result, nil
}

// 使用MySQL获取器
mysqlFetcher := &MySQLFetcher{db: yourDBConnection}
cacheManager := cache.InitCacheManager("app", 100, redisClient, 1000, mysqlFetcher)
```

## 📊 性能优势

| 缓存级别 | 平均响应时间 | 适用场景 |
|---------|-------------|----------|
| 内存缓存 | < 1μs | 热点数据访问 |
| Redis缓存 | < 1ms | 温数据访问 |
| 数据库 | 10-100ms | 冷数据访问 |

## 🛠️ API 参考

### CacheManager

#### InitCacheManager
```go
func InitCacheManager(code string, memCacheSize int, redisClient *redis.Client, redisCacheSize int, dataFetcher fetcher.DataFetcher) *CacheManager
```

**参数说明:**
- `code`: 缓存键前缀
- `memCacheSize`: 内存缓存容量
- `redisClient`: Redis客户端实例 (可为nil)
- `redisCacheSize`: Redis缓存容量
- `dataFetcher`: 数据获取器实现

#### GetData
```go
func (cm *CacheManager) GetData(key string) (string, error)
```

从多级缓存中获取数据，按L1->L2->L3的顺序查找。

### DataFetcher接口

```go
type DataFetcher interface {
    FetchData(key string) (string, error)
}
```

实现此接口以支持不同的数据源。

## 📁 项目结构

```
my-golru/
├── cache/
│   ├── cacheManager.go    # 缓存管理器
│   ├── memoryCache.go     # 内存缓存实现
│   └── redisCache.go      # Redis缓存实现
├── fetcher/
│   └── mysqlFetcher.go    # 数据获取器接口
├── main.go                # 示例代码
├── go.mod                 # 依赖管理
└── README.md              # 项目文档
```

## 🔗 依赖

- `github.com/go-redis/redis/v7` - Redis客户端
- `github.com/go-sql-driver/mysql` - MySQL驱动

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 📄 许可证

本项目采用MIT许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者。
