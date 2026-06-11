// Package handler 提供藏叶的 HTTP 请求处理器。
//
// 每个 Handler 负责一个功能域，通过依赖注入获取数据库连接。
// 本包使用 Go 标准库 net/http，不依赖第三方 Web 框架。
package handler

import (
	"database/sql"
	"fmt"
	"net/http"
)

// DashboardHandler 处理仪表盘相关的 HTTP 请求。
// 通过 DB 字段注入数据库连接。
type DashboardHandler struct {
	// 数据库连接池
	DB *sql.DB
}

// DashboardStats 是仪表盘统计数据的响应结构体。
// 通过 json tag 映射到 JSON 字段名（snake_case）。
type DashboardStats struct {
	// 文件夹总数
	FolderCount int64 `json:"folder_count"`
	// 文件总数
	FileCount int64 `json:"file_count"`
	// 图片数
	ImageCount int64 `json:"image_count"`
	// 视频数
	VideoCount int64 `json:"video_count"`
	// 音频数
	AudioCount int64 `json:"audio_count"`
	// 其他文件数
	OtherCount int64 `json:"other_count"`
	// 存储空间总字节数
	StorageBytes int64 `json:"storage_bytes"`
	// 存储空间的人类可读格式（如 "1.5 GB"）
	StorageDisplay string `json:"storage_display"`
}

// Stats 处理 GET /api/dashboard/stats 请求。
//
// 从数据库查询六项统计数据：
//   - 文件总数
//   - 图片数（mime_type 匹配 image/*）
//   - 视频数（mime_type 匹配 video/*）
//   - 音频数（mime_type 匹配 audio/*）
//   - 其他文件数（不匹配以上类型）
//   - 存储空间总大小（字节数 + 人类可读格式）
func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {
	var stats DashboardStats

	// 查询文件夹总数
	if !scanRow(w, h.DB.QueryRow("SELECT COUNT(*) FROM folders"), &stats.FolderCount) {
		return
	}
	// 查询文件总数
	if !scanRow(w, h.DB.QueryRow("SELECT COUNT(*) FROM files"), &stats.FileCount) {
		return
	}
	// 查询图片数
	scanRow(w, h.DB.QueryRow("SELECT COUNT(*) FROM files WHERE mime_type LIKE 'image/%'"), &stats.ImageCount)
	// 查询视频数
	scanRow(w, h.DB.QueryRow("SELECT COUNT(*) FROM files WHERE mime_type LIKE 'video/%'"), &stats.VideoCount)
	// 查询音频数
	scanRow(w, h.DB.QueryRow("SELECT COUNT(*) FROM files WHERE mime_type LIKE 'audio/%'"), &stats.AudioCount)
	// 查询其他文件数
	scanRow(w, h.DB.QueryRow("SELECT COUNT(*) FROM files WHERE mime_type NOT LIKE 'image/%' AND mime_type NOT LIKE 'video/%' AND mime_type NOT LIKE 'audio/%'"), &stats.OtherCount)
	// 查询存储空间总大小
	scanRow(w, h.DB.QueryRow("SELECT COALESCE(SUM(file_size), 0) FROM files"), &stats.StorageBytes)

	stats.StorageDisplay = formatSize(stats.StorageBytes)
	writeJSON(w, http.StatusOK, stats)
}

// formatSize 将字节数转换为人类可读的存储大小格式。
//
// 使用 1024 进制（二进制前缀）：
//   - 1 KB = 1024 B
//   - 1 MB = 1024 KB
//   - 1 GB = 1024 MB
//   - 1 TB = 1024 GB
//
// 返回值保留一位小数，如 "1.5 GB"。
// 0 字节返回 "0 B"。
func formatSize(bytes int64) string {
	const (
		TB = 1 << 40
		GB = 1 << 30
		MB = 1 << 20
		KB = 1 << 10
	)
	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.1f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}