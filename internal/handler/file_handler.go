// Package handler 提供 HTTP 请求处理器
package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"zangye/internal/model"
	"zangye/internal/service"
	"zangye/internal/util"
)

// FileHandler 文件 HTTP 处理器
type FileHandler struct {
	fileSvc *service.FileSvc
}

// NewFileHandler 创建 FileHandler 实例
func NewFileHandler(fileSvc *service.FileSvc) *FileHandler {
	return &FileHandler{fileSvc: fileSvc}
}

// HandleFiles 处理 /api/files 路由
// GET: 查询文件列表（支持分页、搜索、筛选）
// POST: 创建文件
func (h *FileHandler) HandleFiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listFiles(w, r)
	case http.MethodPost:
		h.createFile(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

// HandleFile 处理 /api/files/{id} 路由
// GET: 获取文件详情
// PUT: 更新文件
// DELETE: 删除文件（软删除）
func (h *FileHandler) HandleFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "缺少文件 ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getFile(w, r, id)
	case http.MethodPut:
		h.updateFile(w, r, id)
	case http.MethodDelete:
		h.deleteFile(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

// HandleFileContent 处理 /api/files/{id}/content 路由
// 支持 Range 请求（视频拖动进度条）
func (h *FileHandler) HandleFileContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "缺少文件 ID")
		return
	}

	file, err := h.fileSvc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if file == nil {
		writeError(w, http.StatusNotFound, "文件不存在")
		return
	}

	// 获取实际文件路径
	filePath, err := h.fileSvc.GetContentPath(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "文件路径不存在: "+err.Error())
		return
	}

	// 根据文件类型处理
	ext := file.FileType
	mimeType := util.GetMimeType(ext)

	if util.IsTextFile(ext) {
		// 文本文件：读取内容返回
		http.ServeFile(w, r, filePath)
	} else if util.IsImageFile(ext) {
		// 图片文件：返回二进制 + Content-Type
		w.Header().Set("Content-Type", mimeType)
		http.ServeFile(w, r, filePath)
	} else if util.IsVideoFile(ext) || util.IsAudioFile(ext) {
		// 视频/音频文件：支持 Range 请求
		w.Header().Set("Accept-Ranges", "bytes")
		http.ServeFile(w, r, filePath)
	} else {
		// 其他文件：作为二进制流返回
		w.Header().Set("Content-Type", mimeType)
		http.ServeFile(w, r, filePath)
	}
}

// HandleFileTags 处理 /api/files/{id}/tags 路由
// GET: 获取文件的标签列表
// PUT: 设置文件的标签列表（全量替换）
// POST: 为文件添加单个标签
// DELETE: 移除文件的单个标签
func (h *FileHandler) HandleFileTags(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "缺少文件 ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getFileTags(w, r, id)
	case http.MethodPut:
		h.setFileTags(w, r, id)
	case http.MethodPost:
		h.addFileTag(w, r, id)
	case http.MethodDelete:
		h.removeFileTag(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "不支持的请求方法")
	}
}

// listFiles 查询文件列表
func (h *FileHandler) listFiles(w http.ResponseWriter, r *http.Request) {
	q := model.FileQuery{
		OrderBy:  r.URL.Query().Get("order_by"),
		OrderDir: r.URL.Query().Get("order_dir"),
	}

	// 解析可选参数
	if folderID := r.URL.Query().Get("folder_id"); folderID != "" {
		q.FolderID = &folderID
	}
	if keyword := r.URL.Query().Get("keyword"); keyword != "" {
		q.Keyword = &keyword
	}
	if fileType := r.URL.Query().Get("type"); fileType != "" {
		q.FileType = &fileType
	}
	if tagID := r.URL.Query().Get("tag_id"); tagID != "" {
		q.TagID = &tagID
	}

	// 分页参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 50
	}
	q.Page = page
	q.PageSize = pageSize

	result, err := h.fileSvc.Query(q)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONList(w, result.Files, result.Total, result.Page, result.Size)
}

// createFile 创建文件
func (h *FileHandler) createFile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FolderID    *string `json:"folder_id"`
		Name        string  `json:"name"`
		Paths       string  `json:"paths"`
		Description string  `json:"description"`
		FileSize    int64   `json:"file_size"`
	}
	if err := parseJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "文件名不能为空")
		return
	}

	file, err := h.fileSvc.Create(req.FolderID, req.Name, req.Paths, req.Description, req.FileSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, file)
}

// getFile 获取文件详情
func (h *FileHandler) getFile(w http.ResponseWriter, r *http.Request, id string) {
	file, err := h.fileSvc.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if file == nil {
		writeError(w, http.StatusNotFound, "文件不存在")
		return
	}

	writeJSON(w, file)
}

// updateFile 更新文件
func (h *FileHandler) updateFile(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		FolderID    *string `json:"folder_id"`
		Name        string  `json:"name"`
		Paths       string  `json:"paths"`
		Description string  `json:"description"`
		FileSize    int64   `json:"file_size"`
	}
	if err := parseJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "文件名不能为空")
		return
	}

	file, err := h.fileSvc.Update(id, req.FolderID, req.Name, req.Paths, req.Description, req.FileSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, file)
}

// deleteFile 软删除文件
func (h *FileHandler) deleteFile(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.fileSvc.SoftDelete(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeOK(w)
}

// getFileTags 获取文件的标签列表
func (h *FileHandler) getFileTags(w http.ResponseWriter, r *http.Request, id string) {
	tags, err := h.fileSvc.GetTags(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, tags)
}

// setFileTags 设置文件的标签列表（全量替换）
func (h *FileHandler) setFileTags(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		TagIDs []string `json:"tag_ids"`
	}
	if err := parseJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}

	if err := h.fileSvc.SetTags(id, req.TagIDs); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeOK(w)
}

// addFileTag 为文件添加单个标签
func (h *FileHandler) addFileTag(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		TagID string `json:"tag_id"`
	}
	if err := parseJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.TagID == "" {
		writeError(w, http.StatusBadRequest, "标签 ID 不能为空")
		return
	}

	if err := h.fileSvc.AddTag(id, req.TagID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeOK(w)
}

// removeFileTag 移除文件的单个标签
func (h *FileHandler) removeFileTag(w http.ResponseWriter, r *http.Request, id string) {
	// 从查询参数获取 tag_id
	tagID := r.URL.Query().Get("tag_id")
	if tagID == "" {
		// 从请求体获取
		var req struct {
			TagID string `json:"tag_id"`
		}
		if err := parseJSONBody(r, &req); err == nil {
			tagID = req.TagID
		}
	}
	if tagID == "" {
		writeError(w, http.StatusBadRequest, "标签 ID 不能为空")
		return
	}

	if err := h.fileSvc.RemoveTag(id, tagID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeOK(w)
}

// ensure unused imports are used
var _ = fmt.Sprintf
var _ = io.Copy
var _ = strings.TrimSpace
var _ = strconv.Itoa
