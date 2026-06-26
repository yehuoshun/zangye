// Package handler 提供 HTTP 请求处理器
package handler

import (
	"encoding/json"
	"net/http"
	"zangye/internal/service"
)

// FolderHandler 文件夹 HTTP 处理器
type FolderHandler struct {
	folderSvc *service.FolderSvc
}

// NewFolderHandler 创建 FolderHandler 实例
func NewFolderHandler(folderSvc *service.FolderSvc) *FolderHandler {
	return &FolderHandler{folderSvc: folderSvc}
}

// HandleFolders 处理 /api/folders 路由
// GET: 获取文件夹树形列表
// POST: 创建文件夹
func (h *FolderHandler) HandleFolders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listFolders(w, r)
	case http.MethodPost:
		h.createFolder(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

// HandleFolder 处理 /api/folders/{id} 路由
// GET: 获取文件夹详情
// PUT: 更新文件夹
// DELETE: 删除文件夹
func (h *FolderHandler) HandleFolder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "缺少文件夹 ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getFolder(w, r, id)
	case http.MethodPut:
		h.updateFolder(w, r, id)
	case http.MethodDelete:
		h.deleteFolder(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

// HandleFolderStats 处理 /api/folders/{id}/stats 路由
func (h *FolderHandler) HandleFolderStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "缺少文件夹 ID")
		return
	}

	stats, err := h.folderSvc.GetStats(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, stats)
}

// listFolders 获取文件夹树形列表
func (h *FolderHandler) listFolders(w http.ResponseWriter, r *http.Request) {
	tree, err := h.folderSvc.GetTree()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, tree)
}

// createFolder 创建文件夹
func (h *FolderHandler) createFolder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string  `json:"name"`
		ParentID *string `json:"parent_id"`
	}
	if err := parseJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "文件夹名称不能为空")
		return
	}

	folder, err := h.folderSvc.Create(req.Name, req.ParentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, folder)
}

// getFolder 获取文件夹详情
func (h *FolderHandler) getFolder(w http.ResponseWriter, r *http.Request, id string) {
	folder, err := h.folderSvc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if folder == nil {
		writeError(w, http.StatusNotFound, "文件夹不存在")
		return
	}

	writeJSON(w, folder)
}

// updateFolder 更新文件夹
func (h *FolderHandler) updateFolder(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Name     string  `json:"name"`
		ParentID *string `json:"parent_id"`
	}
	if err := parseJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "文件夹名称不能为空")
		return
	}

	folder, err := h.folderSvc.Update(id, req.Name, req.ParentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, folder)
}

// deleteFolder 删除文件夹
func (h *FolderHandler) deleteFolder(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.folderSvc.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeOK(w)
}

// ensure *FolderHandler implements http.Handler interfaces
var _ = json.Marshal
