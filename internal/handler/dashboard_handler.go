// Package handler 提供 HTTP 请求处理器
package handler

import (
	"net/http"
	"zangye/internal/service"
)

// DashboardHandler 仪表盘 HTTP 处理器
type DashboardHandler struct {
	dashboardSvc *service.DashboardSvc
}

// NewDashboardHandler 创建 DashboardHandler 实例
func NewDashboardHandler(dashboardSvc *service.DashboardSvc) *DashboardHandler {
	return &DashboardHandler{dashboardSvc: dashboardSvc}
}

// HandleStats 处理 /api/dashboard/stats 路由
// GET: 获取仪表盘统计数据
func (h *DashboardHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
		return
	}

	stats, err := h.dashboardSvc.GetStats()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, stats)
}
