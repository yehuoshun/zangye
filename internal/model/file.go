// Package model 定义数据模型
package model

import (
	"database/sql"
	"time"
)

// File 文件模型
// 对应数据库 files 表
// Go 的 sql.NullString / sql.NullTime 用于处理数据库 NULL 值
// 类比 Java 的 Optional 或可为空的字段类型
type File struct {
	ID          string         `json:"id"`           // UUID 主键
	FolderID    *string        `json:"folder_id"`    // 所属文件夹 ID，*string 表示可为 NULL
	Name        string         `json:"name"`         // 文件名（含扩展名）
	Paths       sql.NullString `json:"paths"`        // 文件路径 JSON 数组，sql.NullString 处理 NULL
	FileType    string         `json:"file_type"`    // 文件类型（扩展名）
	FileSize    int64          `json:"file_size"`    // 文件大小（字节）
	Description sql.NullString `json:"description"`  // 文件描述
	DeletedAt   *time.Time     `json:"deleted_at"`   // 软删除时间，*time.Time 表示可为 NULL
	CreatedAt   time.Time      `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`   // 更新时间
	Tags        []Tag          `json:"tags,omitempty"` // 关联的标签列表（查询时可选加载）
}

// FileQuery 文件查询参数
type FileQuery struct {
	FolderID *string `json:"folder_id"` // 按文件夹筛选
	Keyword  *string `json:"keyword"`   // 按关键字搜索（匹配 name/paths/description）
	FileType *string `json:"type"`      // 按文件类型筛选
	TagID    *string `json:"tag_id"`    // 按标签筛选
	Trash    bool    `json:"trash"`     // 是否查询回收站
	OrderBy  string  `json:"order_by"`  // 排序字段：name/file_type/file_size/created_at
	OrderDir string  `json:"order_dir"` // 排序方向：asc/desc
	Page     int     `json:"page"`      // 页码，从 1 开始
	PageSize int     `json:"page_size"` // 每页数量
}

// FileListResult 文件列表查询结果
type FileListResult struct {
	Files []File `json:"files"`       // 文件列表
	Total int64  `json:"total"`       // 总记录数
	Page  int    `json:"page"`        // 当前页码
	Size  int    `json:"size"`        // 每页数量
}
