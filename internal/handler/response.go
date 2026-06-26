// Package handler 提供 HTTP 请求处理器
// 类比 Java 的 Controller 层
// Go 的 http.HandlerFunc 是函数类型，满足 http.Handler 接口
package handler

import (
	"encoding/json"
	"net/http"
)

// writeJSON 写入 JSON 响应
// 统一响应格式，类比 Spring 的 @ResponseBody + ResponseEntity
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"序列化响应失败"}`, http.StatusInternalServerError)
	}
}

// writeError 写入错误响应
// 统一错误格式：{"error": "描述"}
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// writeJSONList 写入带总数的列表响应
// 格式：{"data": [...], "total": 100, "page": 1, "size": 50}
func writeJSONList(w http.ResponseWriter, data interface{}, total int64, page, size int) {
	writeJSON(w, map[string]interface{}{
		"data":  data,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// writeOK 写入成功响应
func writeOK(w http.ResponseWriter) {
	writeJSON(w, map[string]string{"status": "ok"})
}

// parseJSONBody 解析请求体 JSON 到指定结构体
// 类比 Spring 的 @RequestBody
func parseJSONBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
