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

// mustScanRow 执行 QueryRow 扫描——找不到行时返回 500。
// 用于 Create/Update 后回查，此时数据一定存在。
func mustScanRow(w http.ResponseWriter, row *sql.Row, dest ...any) bool {
	if err := row.Scan(dest...); err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusInternalServerError, "创建/更新后回查失败: 记录不存在")
		} else {
			writeError(w, http.StatusInternalServerError, "数据库查询失败")
		}
		return false
	}
	return true
}

// closeRows 安全关闭 sql.Rows，记录关闭错误。
// 用于替换 defer rows.Close()，确保错误被捕获。
func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("关闭 rows 失败: %v", err)
	}
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
