// Package handler 提供 HTTP 请求处理器
package handler

import (
	"net/http"
	"zangye/internal/service"
)

// TagHandler 标签 HTTP 处理器
type TagHandler struct {
	tagSvc *service.TagSvc
}

// NewTagHandler 创建 TagHandler 实例
func NewTagHandler(tagSvc *service.TagSvc) *TagHandler {
	return &TagHandler{tagSvc: tagSvc}
}

// HandleTags 处理 /api/tags 路由
// GET: 获取所有标签
// POST: 创建标签
func (h *TagHandler) HandleTags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTags(w, r)
	case http.MethodPost:
		h.createTag(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

// HandleTag 处理 /api/tags/{id} 路由
// GET: 获取标签详情
// PUT: 更新标签
// DELETE: 删除标签
func (h *TagHandler) HandleTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "缺少标签 ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTag(w, r, id)
	case http.MethodPut:
		h.updateTag(w, r, id)
	case http.MethodDelete:
		h.deleteTag(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

// listTags 获取所有标签
func (h *TagHandler) listTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.tagSvc.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, tags)
}

// createTag 创建标签
func (h *TagHandler) createTag(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := parseJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "标签名称不能为空")
		return
	}

	tag, err := h.tagSvc.Create(req.Name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, tag)
}

// getTag 获取标签详情
func (h *TagHandler) getTag(w http.ResponseWriter, r *http.Request, id string) {
	tag, err := h.tagSvc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tag == nil {
		writeError(w, http.StatusNotFound, "标签不存在")
		return
	}

	writeJSON(w, tag)
}

// updateTag 更新标签
func (h *TagHandler) updateTag(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	if err := parseJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "标签名称不能为空")
		return
	}

	tag, err := h.tagSvc.Update(id, req.Name, req.Color)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tag == nil {
		writeError(w, http.StatusNotFound, "标签不存在")
		return
	}

	writeJSON(w, tag)
}

// deleteTag 删除标签
func (h *TagHandler) deleteTag(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.tagSvc.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeOK(w)
}
