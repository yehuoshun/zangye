package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// FileTagsHandler 处理文件-标签关联相关的 HTTP 请求。
type FileTagsHandler struct {
	DB *sql.DB
}

// FileTagResponse 是文件-标签关联的响应结构体。
type FileTagResponse struct {
	FileID string `json:"file_id"`
	TagID  string `json:"tag_id"`
}

// SetFileTagsRequest 是设置文件标签的请求体。
type SetFileTagsRequest struct {
	TagIDs []string `json:"tag_ids"`
}

// GetTags GET /api/files/{id}/tags — 获取文件的标签列表
func (h *FileTagsHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("id")

	rows, err := h.DB.Query(
		"SELECT t.id, t.name, t.color, t.created_at FROM tags t INNER JOIN file_tags ft ON t.id = ft.tag_id WHERE ft.file_id = ? ORDER BY t.name",
		fileID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询文件标签失败")
		return
	}
	defer closeRows(rows)

	tags := make([]TagResponse, 0)
	for rows.Next() {
		var t TagResponse
		if err := rows.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "读取标签数据失败")
			return
		}
		tags = append(tags, t)
	}

	writeJSON(w, http.StatusOK, tags)
}

// SetTags PUT /api/files/{id}/tags — 设置文件的标签（全量替换）
func (h *FileTagsHandler) SetTags(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("id")

	// 检查文件是否存在
	var exists int
	if err := h.DB.QueryRow("SELECT 1 FROM files WHERE id = ?", fileID).Scan(&exists); err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "文件不存在")
		return
	}

	var req SetFileTagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式无效")
		return
	}

	// 事务：先删后插
	tx, err := h.DB.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "开启事务失败")
		return
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM file_tags WHERE file_id = ?", fileID); err != nil {
		writeError(w, http.StatusInternalServerError, "清除旧标签失败")
		return
	}

	for _, tagID := range req.TagIDs {
		if _, err := tx.Exec("INSERT IGNORE INTO file_tags (file_id, tag_id) VALUES (?, ?)", fileID, tagID); err != nil {
			writeError(w, http.StatusInternalServerError, "关联标签失败")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, "提交事务失败")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// AddTag POST /api/files/{id}/tags — 给文件添加单个标签
func (h *FileTagsHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("id")

	var req struct {
		TagID string `json:"tag_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式无效")
		return
	}
	if req.TagID == "" {
		writeError(w, http.StatusBadRequest, "tag_id 不能为空")
		return
	}

	_, err := h.DB.Exec("INSERT IGNORE INTO file_tags (file_id, tag_id) VALUES (?, ?)", fileID, req.TagID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "添加标签失败")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

// RemoveTag DELETE /api/files/{id}/tags/{tagId} — 移除文件的单个标签
func (h *FileTagsHandler) RemoveTag(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("id")
	tagID := r.PathValue("tagId")

	result, err := h.DB.Exec("DELETE FROM file_tags WHERE file_id = ? AND tag_id = ?", fileID, tagID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "移除标签失败")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "关联不存在")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}