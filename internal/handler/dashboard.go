// Package handler 提供藏叶的 HTTP 请求处理器。
//
// 每个 Handler 负责一个功能域，通过依赖注入获取数据库连接。
// 本包使用 Go 标准库 net/http，不依赖第三方 Web 框架。
package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// DashboardHandler 处理仪表盘相关的 HTTP 请求。
// 通过 DB 字段注入数据库连接。
type DashboardHandler struct {
	DB *sql.DB // 数据库连接池
}

// DashboardStats 是仪表盘统计数据的响应结构体。
// 通过 json tag 映射到 JSON 字段名（snake_case）。
type DashboardStats struct {
	FileCount      int64  `json:"file_count"`      // 文件总数
	ImageCount     int64  `json:"image_count"`     // 图片数
	VideoCount     int64  `json:"video_count"`     // 视频数
	AudioCount     int64  `json:"audio_count"`     // 音频数
	OtherCount     int64  `json:"other_count"`     // 其他文件数
	StorageBytes   int64  `json:"storage_bytes"`   // 存储空间总字节数
	StorageDisplay string `json:"storage_display"` // 存储空间的人类可读格式（如 "1.5 GB"）
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

	// 查询文件总数
	h.DB.QueryRow("SELECT COUNT(*) FROM files").Scan(&stats.FileCount)
	// 查询图片数
	h.DB.QueryRow("SELECT COUNT(*) FROM files WHERE mime_type LIKE 'image/%'").Scan(&stats.ImageCount)
	// 查询视频数
	h.DB.QueryRow("SELECT COUNT(*) FROM files WHERE mime_type LIKE 'video/%'").Scan(&stats.VideoCount)
	// 查询音频数
	h.DB.QueryRow("SELECT COUNT(*) FROM files WHERE mime_type LIKE 'audio/%'").Scan(&stats.AudioCount)
	// 查询其他文件数（不匹配 image/video/audio）
	h.DB.QueryRow("SELECT COUNT(*) FROM files WHERE mime_type NOT LIKE 'image/%' AND mime_type NOT LIKE 'video/%' AND mime_type NOT LIKE 'audio/%'").Scan(&stats.OtherCount)
	// 查询存储空间总大小（COALESCE 处理空表返回 NULL 的情况）
	h.DB.QueryRow("SELECT COALESCE(SUM(file_size), 0) FROM files").Scan(&stats.StorageBytes)
	// 将字节数转换为人类可读格式
	stats.StorageDisplay = formatSize(stats.StorageBytes)

	// 返回 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
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
	switch {
	case bytes >= 1<<40: // 1 TB = 2^40
		return fmt.Sprintf("%.1f TB", float64(bytes)/(1<<40))
	case bytes >= 1<<30: // 1 GB = 2^30
		return fmt.Sprintf("%.1f GB", float64(bytes)/(1<<30))
	case bytes >= 1<<20: // 1 MB = 2^20
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1<<20))
	case bytes >= 1<<10: // 1 KB = 2^10
		return fmt.Sprintf("%.1f KB", float64(bytes)/(1<<10))
	default: // 小于 1 KB
		return "0 B"
	}
}