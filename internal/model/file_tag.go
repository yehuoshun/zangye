// Package model 定义数据模型
package model

import "time"

// FileTag 文件-标签关联模型
// 对应数据库 file_tags 表
// 多对多关系的中间表
type FileTag struct {
	FileID    string    `json:"file_id"`    // 文件 ID
	TagID     string    `json:"tag_id"`     // 标签 ID
	CreatedAt time.Time `json:"created_at"` // 关联创建时间
}
