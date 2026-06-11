package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// SettingsHandler 处理系统设置相关的 HTTP 请求。
type SettingsHandler struct {
	DB *sql.DB
}

// GetAll GET /api/settings — 获取所有设置项（键值对）
func (h *SettingsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT `key`, `value` FROM settings")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询设置失败")
		return
	}
	defer closeRows(rows)

	settings := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			writeError(w, http.StatusInternalServerError, "读取设置数据失败")
			return
		}
		settings[k] = v
	}

	writeJSON(w, http.StatusOK, settings)
}

// Update PUT /api/settings — 批量更新设置项
// 请求体：{"key1": "value1", "key2": "value2"}
func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式无效")
		return
	}

	for k, v := range data {
		_, err := h.DB.Exec("INSERT INTO settings (`key`, `value`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `value` = ?", k, v, v)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "保存设置失败")
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}