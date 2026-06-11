package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// FoldersHandler 处理文件夹管理相关的 HTTP 请求。
type FoldersHandler struct {
	DB *sql.DB
}

// FolderResponse 是文件夹 API 的响应结构体。
type FolderResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Icon      string  `json:"icon"`
	ParentID  *string `json:"parent_id"`
	SortOrder int     `json:"sort_order"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// FolderCreateRequest 是创建文件夹的请求体。
type FolderCreateRequest struct {
	Name      string  `json:"name"`
	Icon      string  `json:"icon"`
	ParentID  *string `json:"parent_id"`
	SortOrder int     `json:"sort_order"`
}

// FolderUpdateRequest 是更新文件夹的请求体。
type FolderUpdateRequest struct {
	Name      *string `json:"name"`
	Icon      *string `json:"icon"`
	ParentID  *string `json:"parent_id"`
	SortOrder *int    `json:"sort_order"`
}

// List GET /api/folders — 获取文件夹列表
// 支持 ?parent_id=xxx 查询子文件夹，不传则返回根目录下的文件夹
func (h *FoldersHandler) List(w http.ResponseWriter, r *http.Request) {
	parentID := r.URL.Query().Get("parent_id")

	var rows *sql.Rows
	var err error
	if parentID == "" {
		rows, err = h.DB.Query(
			"SELECT id, name, icon, parent_id, sort_order, created_at, updated_at FROM folders WHERE parent_id IS NULL ORDER BY sort_order, name",
		)
	} else {
		rows, err = h.DB.Query(
			"SELECT id, name, icon, parent_id, sort_order, created_at, updated_at FROM folders WHERE parent_id = ? ORDER BY sort_order, name",
			parentID,
		)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询文件夹列表失败")
		return
	}
	defer rows.Close()

	folders := make([]FolderResponse, 0)
	for rows.Next() {
		var f FolderResponse
		if err := rows.Scan(&f.ID, &f.Name, &f.Icon, &f.ParentID, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "读取文件夹数据失败")
			return
		}
		folders = append(folders, f)
	}

	writeJSON(w, http.StatusOK, folders)
}

// Get GET /api/folders/{id} — 获取单个文件夹详情
func (h *FoldersHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var f FolderResponse
	err := h.DB.QueryRow(
		"SELECT id, name, icon, parent_id, sort_order, created_at, updated_at FROM folders WHERE id = ?",
		id,
	).Scan(&f.ID, &f.Name, &f.Icon, &f.ParentID, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "文件夹不存在")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询文件夹失败")
		return
	}

	writeJSON(w, http.StatusOK, f)
}

// Create POST /api/folders — 创建文件夹
func (h *FoldersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req FolderCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体格式无效")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "文件夹名称不能为空")
		return
	}

	icon := req.Icon
	if icon == "" {
		icon = "📁"
	}

	id := uuid.New().String()
	_, err := h.DB.Exec(
		"INSERT INTO folders (id, name, icon, parent_id, sort_order) VALUES (?, ?, ?, ?, ?)",
		id, req.Name, icon, req.ParentID, req.SortOrder,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "创建文件夹失败")
		return
	}

	var f FolderResponse
	h.DB.QueryRow(
		"SELECT id, name, icon, parent_id, sort_order, created_at, updated_at FROM folders WHERE id = ?",
		id,
	).Scan(&f.ID, &f.Name, &f.Icon, &f.ParentID, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt)

	writeJSON(w, http.StatusCreated, f)
}

// Update PUT /api/folders/{id} — 更新文件夹
func (h *FoldersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var exists int
	if err := h.DB.QueryRow("SELECT 1 FROM folders WHERE id = ?", id).Scan(&exists); err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "文件夹不存在")
		return
	}

	var req FolderUpdateRequest
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
	if req.Icon != nil {
		fields += "icon = ?, "
		args = append(args, *req.Icon)
	}
	if req.ParentID != nil {
		fields += "parent_id = ?, "
		args = append(args, *req.ParentID)
	}
	if req.SortOrder != nil {
		fields += "sort_order = ?, "
		args = append(args, *req.SortOrder)
	}

	if fields == "" {
		writeError(w, http.StatusBadRequest, "没有需要更新的字段")
		return
	}

	fields = fields[:len(fields)-2]
	args = append(args, id)

	_, err := h.DB.Exec("UPDATE folders SET "+fields+" WHERE id = ?", args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "更新文件夹失败")
		return
	}

	var f FolderResponse
	h.DB.QueryRow(
		"SELECT id, name, icon, parent_id, sort_order, created_at, updated_at FROM folders WHERE id = ?",
		id,
	).Scan(&f.ID, &f.Name, &f.Icon, &f.ParentID, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt)

	writeJSON(w, http.StatusOK, f)
}

// Delete DELETE /api/folders/{id} — 删除文件夹（级联删除子文件和子文件夹）
func (h *FoldersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := h.DB.Exec("DELETE FROM folders WHERE id = ?", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "删除文件夹失败")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "文件夹不存在")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}
