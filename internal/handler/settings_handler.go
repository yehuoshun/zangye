// Package handler 提供 HTTP 请求处理器
package handler

import (
	"net/http"
	"zangye/internal/repository"
)

// SettingsHandler 设置 HTTP 处理器
type SettingsHandler struct {
	settingsRepo *repository.SettingsRepo
}

// NewSettingsHandler 创建 SettingsHandler 实例
func NewSettingsHandler(settingsRepo *repository.SettingsRepo) *SettingsHandler {
	return &SettingsHandler{settingsRepo: settingsRepo}
}

// HandleSettings 处理 /api/settings 路由
// GET: 获取所有设置
// PUT: 更新设置
func (h *SettingsHandler) HandleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getSettings(w, r)
	case http.MethodPut:
		h.updateSettings(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

// getSettings 获取所有设置
func (h *SettingsHandler) getSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settingsRepo.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, settings)
}

// updateSettings 更新设置
// 接收 JSON 对象，批量更新
func (h *SettingsHandler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var settings map[string]string
	if err := parseJSONBody(r, &settings); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}

	for key, value := range settings {
		if err := h.settingsRepo.Set(key, value); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	writeOK(w)
}
