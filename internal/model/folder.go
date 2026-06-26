// Package model 定义数据模型
// Go 没有 class，struct + 方法接收器 ≈ Java 的类
// 零值 zero value：int=0, string="", 指针=nil，不会 NPE
package model

import "time"

// Folder 文件夹模型
// 对应数据库 folders 表
// Go 的 struct 字段标签 `json:"xxx"` 用于序列化/反序列化，类似 Jackson 注解
type Folder struct {
	ID        string     `json:"id"`         // UUID 主键
	Name      string     `json:"name"`       // 文件夹名称
	ParentID  *string    `json:"parent_id"`  // 父文件夹 ID，*string 表示可为 NULL（对应数据库 NULL）
	SortOrder int        `json:"sort_order"` // 排序序号
	CreatedAt time.Time  `json:"created_at"` // 创建时间
	UpdatedAt time.Time  `json:"updated_at"` // 更新时间
	Children  []*Folder  `json:"children,omitempty"`  // 子文件夹列表（树形结构用，omitempty 表示空时不序列化）
}

// FolderTreeItem 文件夹树节点（含文件数量统计）
type FolderTreeItem struct {
	Folder
	FileCount int64 `json:"file_count"` // 该文件夹下直接文件数量
}

// FolderStats 文件夹统计信息
type FolderStats struct {
	TotalFolders int64 `json:"total_folders"` // 递归子文件夹总数
	TotalFiles   int64 `json:"total_files"`   // 递归文件总数
	TotalSize    int64 `json:"total_size"`    // 递归文件总大小（字节）
}
