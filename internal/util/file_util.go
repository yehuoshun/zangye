// Package util 提供通用工具函数
package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 文件类型分类映射
var mimeTypes = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
	"bmp":  "image/bmp",
	"svg":  "image/svg+xml",
	"ico":  "image/x-icon",
	"mp4":  "video/mp4",
	"avi":  "video/x-msvideo",
	"mov":  "video/quicktime",
	"wmv":  "video/x-ms-wmv",
	"flv":  "video/x-flv",
	"mkv":  "video/x-matroska",
	"webm": "video/webm",
	"mp3":  "audio/mpeg",
	"wav":  "audio/wav",
	"ogg":  "audio/ogg",
	"flac": "audio/flac",
	"aac":  "audio/aac",
	"wma":  "audio/x-ms-wma",
	"pdf":  "application/pdf",
	"doc":  "application/msword",
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"xls":  "application/vnd.ms-excel",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"ppt":  "application/vnd.ms-powerpoint",
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"zip":  "application/zip",
	"rar":  "application/vnd.rar",
	"7z":   "application/x-7z-compressed",
	"tar":  "application/x-tar",
	"gz":   "application/gzip",
	"txt":  "text/plain",
	"md":   "text/markdown",
	"json": "application/json",
	"xml":  "application/xml",
	"html": "text/html",
	"htm":  "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"ts":   "application/typescript",
	"go":   "text/x-go",
	"py":   "text/x-python",
	"java": "text/x-java",
	"sql":  "text/x-sql",
	"yaml": "text/yaml",
	"yml":  "text/yaml",
}

// 文本文件扩展名集合
var textExtensions = map[string]bool{
	"txt": true, "md": true, "json": true, "xml": true, "html": true,
	"htm": true, "css": true, "js": true, "ts": true, "go": true,
	"py": true, "java": true, "sql": true, "yaml": true, "yml": true,
	"log": true, "cfg": true, "ini": true, "conf": true, "sh": true,
	"bat": true, "ps1": true, "env": true, "gitignore": true,
}

// 图片文件扩展名集合
var imageExtensions = map[string]bool{
	"jpg": true, "jpeg": true, "png": true, "gif": true,
	"webp": true, "bmp": true, "svg": true, "ico": true,
}

// 视频文件扩展名集合
var videoExtensions = map[string]bool{
	"mp4": true, "avi": true, "mov": true, "wmv": true,
	"flv": true, "mkv": true, "webm": true,
}

// 音频文件扩展名集合
var audioExtensions = map[string]bool{
	"mp3": true, "wav": true, "ogg": true, "flac": true,
	"aac": true, "wma": true,
}

// GetFileType 从文件名推断文件类型（扩展名）
// 例如："photo.jpg" -> "jpg"，"document.pdf" -> "pdf"
// 如果没有扩展名，返回 "unknown"
func GetFileType(name string) string {
	ext := strings.TrimPrefix(filepath.Ext(name), ".")
	if ext == "" {
		return "unknown"
	}
	return strings.ToLower(ext)
}

// GetMimeType 根据文件扩展名获取 MIME 类型
// 类比 Java 的 MimetypesFileTypeMap
func GetMimeType(ext string) string {
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")
	if mime, ok := mimeTypes[ext]; ok {
		return mime
	}
	return "application/octet-stream" // 未知类型默认二进制流
}

// IsTextFile 判断是否为文本文件
func IsTextFile(ext string) bool {
	return textExtensions[strings.TrimPrefix(strings.ToLower(ext), ".")]
}

// IsImageFile 判断是否为图片文件
func IsImageFile(ext string) bool {
	return imageExtensions[strings.TrimPrefix(strings.ToLower(ext), ".")]
}

// IsVideoFile 判断是否为视频文件
func IsVideoFile(ext string) bool {
	return videoExtensions[strings.TrimPrefix(strings.ToLower(ext), ".")]
}

// IsAudioFile 判断是否为音频文件
func IsAudioFile(ext string) bool {
	return audioExtensions[strings.TrimPrefix(strings.ToLower(ext), ".")]
}

// FileExists 检查文件是否存在
// Go 的 os.Stat 返回 error，用 os.IsNotExist 判断文件不存在
// 类比 Java 的 File.exists()
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// FormatFileSize 格式化文件大小为人类可读字符串
// 自动选择 B/KB/MB/GB/TB 单位
// 例如：1024 -> "1.00 KB", 1048576 -> "1.00 MB"
func FormatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	} else if size < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	}
	return fmt.Sprintf("%.2f TB", float64(size)/(1024*1024*1024*1024))
}

// GetCategory 根据扩展名返回文件大类
// 用于仪表盘统计分类：image/video/audio/other
func GetCategory(ext string) string {
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")
	if IsImageFile(ext) {
		return "image"
	}
	if IsVideoFile(ext) {
		return "video"
	}
	if IsAudioFile(ext) {
		return "audio"
	}
	return "other"
}
