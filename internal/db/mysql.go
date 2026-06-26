// Package db 提供数据库连接管理
// Go 的 database/sql 是标准库自带的数据库操作接口，类似 Java 的 JDBC
// 需要配合具体数据库驱动使用（这里用 go-sql-driver/mysql）
package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动，_ 表示只执行 init() 不直接使用
)

// Init 初始化数据库连接
// Go 函数返回 (*sql.DB, error) 是标准模式：返回结果 + 可能的错误
// 类比 Java：返回 Connection Pool + 抛出 SQLException
// dsn 格式：user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func Init(dsn string) (*sql.DB, error) {
	// sql.Open 只验证 DSN 格式，不会真正连接数据库
	// 类比 Java 的 DataSource 初始化
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %w", err)
	}

	// 设置连接池参数
	// Go 的 sql.DB 自带连接池，无需额外配置连接池框架
	db.SetMaxOpenConns(25)                 // 最大打开连接数
	db.SetMaxIdleConns(10)                 // 最大空闲连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最大存活时间
	db.SetConnMaxIdleTime(2 * time.Minute) // 空闲连接超时

	// db.Ping() 真正建立连接，验证 DSN 是否正确
	// 类比 Java 的 Connection.isValid()
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return db, nil
}
