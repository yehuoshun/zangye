// Package model 定义数据模型
package model

import "time"

// Tag 标签模型
// 对应数据库 tags 表
type Tag struct {
	ID        string    `json:"id"`         // UUID 主键
	Name      string    `json:"name"`       // 标签名称
	Color     string    `json:"color"`      // 标签颜色（十六进制，如 #FF6B6B）
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// TagWithCount 标签模型（含文件数量统计）
type TagWithCount struct {
	Tag
	FileCount int64 `json:"file_count"` // 关联的文件数量
}
