package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// TagsHandler 处理标签管理相关的 HTTP 请求。
type TagsHandler struct {
	DB *sql.DB
}

// TagResponse 是标签 API 的响应结构体。
type TagResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	CreatedAt string `json:"created_at"`
}

// TagCreateRequest 是创建标签的请求体。
type TagCreateRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// TagUpdateRequest 是更新标签的请求体。
type TagUpdateRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

// List GET /api/tags — 获取所有标签
func (h *TagsHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name, color, created_at FROM tags ORDER BY name")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询标签列表失败")
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

// Create POST /api/tags — 创建标签
func (h *TagsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req TagCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式无效")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "标签名称不能为空")
		return
	}

	color := req.Color
	if color == "" {
		color = "gray"
	}

	id := uuid.New().String()
	_, err := h.DB.Exec(
		"INSERT INTO tags (id, name, color) VALUES (?, ?, ?)",
		id, req.Name, color,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "创建标签失败（名称可能重复）")
		return
	}

	var t TagResponse
	if !mustScanRow(w, h.DB.QueryRow(
		"SELECT id, name, color, created_at FROM tags WHERE id = ?", id,
	), &t.ID, &t.Name, &t.Color, &t.CreatedAt) {
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

// Update PUT /api/tags/{id} — 更新标签
func (h *TagsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var exists int
	if err := h.DB.QueryRow("SELECT 1 FROM tags WHERE id = ?", id).Scan(&exists); err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "标签不存在")
		return
	}

	var req TagUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式无效")
		return
	}

	fields := ""
	args := []any{}
	if req.Name != nil {
		fields += "name = ?, "
		args = append(args, *req.Name)
	}
	if req.Color != nil {
		fields += "color = ?, "
		args = append(args, *req.Color)
	}

	if fields == "" {
		writeError(w, http.StatusBadRequest, "没有需要更新的字段")
		return
	}

	fields = fields[:len(fields)-2]
	args = append(args, id)

	_, err := h.DB.Exec("UPDATE tags SET "+fields+" WHERE id = ?", args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "更新标签失败")
		return
	}

	var t TagResponse
	if !mustScanRow(w, h.DB.QueryRow(
		"SELECT id, name, color, created_at FROM tags WHERE id = ?", id,
	), &t.ID, &t.Name, &t.Color, &t.CreatedAt) {
		return
	}

	writeJSON(w, http.StatusOK, t)
}

// Delete DELETE /api/tags/{id} — 删除标签
func (h *TagsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := h.DB.Exec("DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "删除标签失败")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "标签不存在")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}