package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// FilesHandler 处理文件管理相关的 HTTP 请求。
type FilesHandler struct {
	DB *sql.DB
}

// FileResponse 是文件 API 的响应结构体。
type FileResponse struct {
	ID           string  `json:"id"`
	FolderID string  `json:"collection_id"`
	Path         string  `json:"path"`
	DisplayName  *string `json:"display_name"`
	FileSize     int64   `json:"file_size"`
	MimeType     *string `json:"mime_type"`
	SortOrder    int     `json:"sort_order"`
	CreatedAt    string  `json:"created_at"`
}

// FileCreateRequest 是创建文件的请求体。
type FileCreateRequest struct {
	FolderID string  `json:"collection_id"`
	Path         string  `json:"path"`
	DisplayName  *string `json:"display_name"`
	FileSize     int64   `json:"file_size"`
	MimeType     *string `json:"mime_type"`
	SortOrder    int     `json:"sort_order"`
}

// FileUpdateRequest 是更新文件的请求体。
type FileUpdateRequest struct {
	FolderID *string `json:"collection_id"`
	Path         *string `json:"path"`
	DisplayName  *string `json:"display_name"`
	FileSize     *int64  `json:"file_size"`
	MimeType     *string `json:"mime_type"`
	SortOrder    *int    `json:"sort_order"`
}

// List GET /api/files — 获取文件列表
func (h *FilesHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(
		"SELECT id, folder_id, path, display_name, file_size, mime_type, sort_order, created_at FROM files ORDER BY sort_order, created_at DESC",
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询文件列表失败")
		return
	}
	defer rows.Close()

	files := make([]FileResponse, 0)
	for rows.Next() {
		var f FileResponse
		if err := rows.Scan(&f.ID, &f.FolderID, &f.Path, &f.DisplayName, &f.FileSize, &f.MimeType, &f.SortOrder, &f.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "读取文件数据失败")
			return
		}
		files = append(files, f)
	}

	writeJSON(w, http.StatusOK, files)
}

// Get GET /api/files/{id} — 获取单个文件详情
func (h *FilesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var f FileResponse
	err := h.DB.QueryRow(
		"SELECT id, folder_id, path, display_name, file_size, mime_type, sort_order, created_at FROM files WHERE id = ?",
		id,
	).Scan(&f.ID, &f.FolderID, &f.Path, &f.DisplayName, &f.FileSize, &f.MimeType, &f.SortOrder, &f.CreatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "文件不存在")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询文件失败")
		return
	}

	writeJSON(w, http.StatusOK, f)
}

// Create POST /api/files — 创建文件
func (h *FilesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req FileCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式无效")
		return
	}

	if req.Path == "" {
		writeError(w, http.StatusBadRequest, "文件路径不能为空")
		return
	}
	if req.FolderID == "" {
		writeError(w, http.StatusBadRequest, "所属文件夹不能为空")
		return
	}

	id := uuid.New().String()
	_, err := h.DB.Exec(
		"INSERT INTO files (id, folder_id, path, display_name, file_size, mime_type, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, req.FolderID, req.Path, req.DisplayName, req.FileSize, req.MimeType, req.SortOrder,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "创建文件失败")
		return
	}

	var f FileResponse
	h.DB.QueryRow(
		"SELECT id, folder_id, path, display_name, file_size, mime_type, sort_order, created_at FROM files WHERE id = ?",
		id,
	).Scan(&f.ID, &f.FolderID, &f.Path, &f.DisplayName, &f.FileSize, &f.MimeType, &f.SortOrder, &f.CreatedAt)

	writeJSON(w, http.StatusCreated, f)
}

// Update PUT /api/files/{id} — 更新文件
func (h *FilesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// 检查文件是否存在
	var exists int
	if err := h.DB.QueryRow("SELECT 1 FROM files WHERE id = ?", id).Scan(&exists); err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "文件不存在")
		return
	}

	var req FileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式无效")
		return
	}

	// 构建动态更新 SQL
	fields := ""
	args := []any{}
	if req.FolderID != nil {
		fields += "folder_id = ?, "
		args = append(args, *req.FolderID)
	}
	if req.Path != nil {
		fields += "path = ?, "
		args = append(args, *req.Path)
	}
	if req.DisplayName != nil {
		fields += "display_name = ?, "
		args = append(args, *req.DisplayName)
	}
	if req.FileSize != nil {
		fields += "file_size = ?, "
		args = append(args, *req.FileSize)
	}
	if req.MimeType != nil {
		fields += "mime_type = ?, "
		args = append(args, *req.MimeType)
	}
	if req.SortOrder != nil {
		fields += "sort_order = ?, "
		args = append(args, *req.SortOrder)
	}

	if fields == "" {
		writeError(w, http.StatusBadRequest, "没有需要更新的字段")
		return
	}

	// 去掉末尾的 ", "
	fields = fields[:len(fields)-2]
	args = append(args, id)

	_, err := h.DB.Exec("UPDATE files SET "+fields+" WHERE id = ?", args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "更新文件失败")
		return
	}

	// 返回更新后的文件
	var f FileResponse
	h.DB.QueryRow(
		"SELECT id, folder_id, path, display_name, file_size, mime_type, sort_order, created_at FROM files WHERE id = ?",
		id,
	).Scan(&f.ID, &f.FolderID, &f.Path, &f.DisplayName, &f.FileSize, &f.MimeType, &f.SortOrder, &f.CreatedAt)

	writeJSON(w, http.StatusOK, f)
}

// Delete DELETE /api/files/{id} — 删除文件
func (h *FilesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := h.DB.Exec("DELETE FROM files WHERE id = ?", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "删除文件失败")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "文件不存在")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}
