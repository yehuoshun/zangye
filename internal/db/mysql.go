// Package db 提供藏叶的 MySQL 数据库连接和自动建表功能。
//
// 核心职责：
//   - 创建 MySQL 数据库连接池
//   - 自动执行 schema.sql 建表语句（幂等，使用 IF NOT EXISTS）
//   - 管理连接池参数（最大连接数、空闲连接数、连接生命周期）
//
// 数据库配置通过 Config 结构体传入，默认连接本地 MySQL。
// 建表语句通过 embed 嵌入 schema.sql，编译时打包进二进制文件。
package db

import (
	"database/sql"
	// 用于嵌入 schema.sql 文件
	_ "embed"
	"fmt"
	"time"

	// MySQL 驱动，init() 中注册，仅需副作用导入
	_ "github.com/go-sql-driver/mysql"
)

//go:embed schema.sql
// schemaSQL 包含所有建表语句，在 New() 中自动执行。
// 所有 CREATE TABLE 都使用 IF NOT EXISTS，确保幂等。
var schemaSQL string

// Config 定义 MySQL 数据库连接参数。
type Config struct {
	// 数据库用户名
	User string
	// 数据库密码
	Password string
	// 数据库主机地址
	Host string
	// 数据库端口
	Port int
	// 数据库名
	DBName string
}

// DefaultConfig 返回默认的数据库配置。
// 默认连接本地 MySQL（127.0.0.1:3306），数据库名为 zang_ye。
// 生产环境应通过环境变量或配置文件覆盖。
func DefaultConfig() Config {
	return Config{
		User:     "root",
		Password: "root",
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "zang_ye",
	}
}

// New 创建 MySQL 数据库连接并自动初始化表结构。
//
// 流程：
//  1. 根据 Config 构建 DSN（数据源名称）
//  2. 打开数据库连接并配置连接池参数
//  3. 执行 Ping 测试连接
//  4. 执行 schema.sql 自动建表（幂等）
//
// 连接池参数：
//   - MaxOpenConns: 10（最大打开连接数）
//   - MaxIdleConns: 5（最大空闲连接数）
//   - ConnMaxLifetime: 5 分钟（连接最大存活时间）
//
// 返回 *sql.DB 和可能的错误。调用方负责在程序退出时调用 Close()。
func New(cfg Config) (*sql.DB, error) {
	// 构建 DSN（Data Source Name）
	// charset=utf8mb4：支持 emoji 等 4 字节 UTF-8 字符
	// parseTime=true：自动将 MySQL 时间类型转为 Go time.Time
	// loc=Local：使用本地时区
	// multiStatements=true：允许在一个 Exec 中执行多条 SQL（建表需要）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	// 打开数据库连接（此时还未实际连接，只是创建连接池）
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接: %w", err)
	}

	// 配置连接池
	// 最大打开连接数
	database.SetMaxOpenConns(10)
	// 最大空闲连接数
	database.SetMaxIdleConns(5)
	// 连接最大存活时间
	database.SetConnMaxLifetime(5 * time.Minute)

	// 测试数据库连接是否可用
	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	// 自动执行建表语句（IF NOT EXISTS 保证幂等）
	if _, err := database.Exec(schemaSQL); err != nil {
		database.Close()
		return nil, fmt.Errorf("建表失败: %w", err)
	}

	return database, nil
}