package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/muyi-zcy/my-golru/cache"
	"log"
	"time"
)

type MySQLFetcher struct {
	db *sql.DB
}

func InitMySQLFetcher(dsn string) (*MySQLFetcher, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &MySQLFetcher{db: db}, nil
}

func (f *MySQLFetcher) FetchData(key string) (string, error) {
	var name string
	query := "SELECT Path FROM t_file_info WHERE id = ?"
	err := f.db.QueryRow(query, key).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("error fetching data from MySQL: %v", err)
	}
	fmt.Println("查询了")
	return name, nil
}

func main() {
	redisAddr := "192.168.0.109:6379"                                              // Redis 地址
	redisDB := 0                                                                   // Redis DB
	mysqlDSN := "root:devMysqlPasswd@tcp(192.168.0.109:3306)/my-ideaistudio-guixu" // MySQL 数据源名称

	var rdb *redis.Client
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "devRedisPasswd",
		DB:       redisDB,
	})

	// 初始化 MySQL 数据获取器
	mysqlFetcher, err := InitMySQLFetcher(mysqlDSN)
	if err != nil {
		log.Fatalf("Error initializing MySQLFetcher: %v", err)
	}

	// 创建缓存管理器，假设内存缓存大小为 5，使用 Redis 作为二级缓存
	cacheManager := cache.InitCacheManager("test", 5, rdb, 5, mysqlFetcher)

	// 测试：获取数据
	key := "237336134413742080"
	startTime := time.Now()
	fmt.Print(cacheManager.GetData(key))
	fmt.Printf("Execution Time: %v\n", time.Since(startTime))
	startTime = time.Now()
	fmt.Println(cacheManager.GetData("237286465042149376"))
	fmt.Printf("Execution Time: %v\n", time.Since(startTime))
	startTime = time.Now()
	fmt.Println(cacheManager.GetData("237287795580235776"))
	fmt.Printf("Execution Time: %v\n", time.Since(startTime))

	startTime = time.Now()
	fmt.Println(cacheManager.GetData("237287795580235776"))
	fmt.Printf("Execution Time: %v\n", time.Since(startTime))
	startTime = time.Now()
	fmt.Println(cacheManager.GetData("237343968048214016"))
	fmt.Printf("Execution Time: %v\n", time.Since(startTime))

	//
	// 如果不使用 Redis，可以将 rdb 传入 nil
	cacheManagerNoRedis := cache.InitCacheManager("test_noredis", 5, nil, 0, mysqlFetcher)

	// 测试：获取数据
	resultNoRedis, err := cacheManagerNoRedis.GetData(key)
	if err != nil {
		log.Fatalf("Error getting data: %v", err)
	}
	fmt.Println("Result without Redis:", resultNoRedis)
}
