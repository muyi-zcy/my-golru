# Redis 缓存管理器

这是一个用于实现 **内存缓存 -> Redis 缓存 -> MySQL 数据库** 三级缓存机制的 Golang 库。该项目实现了一个高效的缓存管理器，通过内存缓存和 Redis 缓存相结合的方式，减少数据库的访问次数，提高数据访问效率。

## 项目结构

- **内存缓存**：基于内存实现的简单缓存，支持容量限制。
- **Redis 缓存**：在 Redis 中存储缓存数据
- **MySQL 数据库**：使用数据获取器从 MySQL 数据库获取数据。

该项目支持两级缓存：内存缓存（一级缓存）和 Redis 缓存（二级缓存）。当数据在内存和 Redis 中都不存在时，它将从 MySQL 数据库加载。

## 特性

- 支持内存缓存、Redis 缓存和 MySQL 数据库三层缓存策略。
- Redis 缓存采用哈希表存储，并使用链表记录访问顺序来管理缓存淘汰。
- 为每个缓存项设置唯一的 `code`，可以通过 `code` 查找和删除缓存项。
- 提供缓存的添加、更新和删除功能。

## 安装

你可以通过 `go get` 命令来安装本项目：

```bash
go get github.com/muyi-zcy/my-golru
