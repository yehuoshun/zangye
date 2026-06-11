// Package handler 提供藏叶的 HTTP 请求处理器。
//
// 每个 Handler 负责一个功能域，通过依赖注入获取数据库连接。
// 本包使用 Go 标准库 net/http，不依赖第三方 Web 框架。
package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSON 将数据序列化为 JSON 并写入响应。
// 序列化失败时返回 500 并记录日志。
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON 编码失败: %v", err)
	}
}

// writeError 返回统一的错误 JSON 响应。
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
