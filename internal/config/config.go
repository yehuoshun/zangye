// Package config 提供应用配置管理
// Go 没有 class，struct + 方法接收器 ≈ Java 的类
// 零值 zero value：int=0, string="", 指针=nil，不会 NPE
package config

import (
	"os"
)

// Config 应用配置结构体
// Go 的 struct 相当于 Java 的 POJO/JavaBean，但不需要 getter/setter
// 字段名大写 = public，小写 = private（包内可见）
type Config struct {
	// Addr HTTP 服务监听地址，格式 "host:port"
	Addr string

	// DSN MySQL 数据源名称，格式 "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	DSN string
}

// Load 加载配置，优先读取环境变量，不存在则使用默认值
// Go 函数多返回值 ≈ Java 的异常，但更轻量
// 返回 *Config 指针而不是 Config 值，避免复制整个结构体
func Load() *Config {
	// Go 的 := 是短变量声明，自动推断类型，相当于 Java 的 var
	addr := os.Getenv("ZANGYE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:27138"
	}

	dsn := os.Getenv("ZANGYE_DSN")
	if dsn == "" {
		// 默认 DSN，生产环境请通过环境变量覆盖
		dsn = "root:root@tcp(127.0.0.1:3306)/zangye?charset=utf8mb4&parseTime=True&loc=Local"
	}

	return &Config{
		Addr: addr,
		DSN:  dsn,
	}
}
