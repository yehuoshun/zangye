// Package handler 提供藏叶的 HTTP 请求处理器。
//
// 每个 Handler 负责一个功能域，通过依赖注入获取数据库连接。
// 本包使用 Go 标准库 net/http，不依赖第三方 Web 框架。
package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

// scanRow 执行 QueryRow 并将结果写入 dest。
// 返回 true 表示找到行，false 表示没有行（sql.ErrNoRows）。
// 其他错误通过 HTTP 响应返回。
func scanRow(w http.ResponseWriter, row *sql.Row, dest ...any) bool {
	if err := row.Scan(dest...); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		writeError(w, http.StatusInternalServerError, "数据库查询失败")
		return false
	}
	return true
}

// buildInQuery 构建 SQL IN 查询，如 "SELECT ... WHERE id IN (?, ?, ?)"。
func buildInQuery(prefix string, ids []string) (string, []any) {
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	return prefix + " (" + strings.Join(placeholders, ",") + ")", args
}
