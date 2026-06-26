// Package util 提供通用工具函数
package util

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateUUID 生成 UUID v4 字符串
// 使用 google/uuid 库，类比 Java 的 UUID.randomUUID().toString()
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateUUIDShort 生成短 UUID（取前 8 位）
// 用于需要短 ID 的场景
func GenerateUUIDShort() string {
	id := uuid.New().String()
	return id[:8]
}

// ParseUUID 解析 UUID 字符串，验证其合法性
// 返回 error 表示格式无效
func ParseUUID(s string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("无效的 UUID 格式: %w", err)
	}
	return parsed, nil
}
