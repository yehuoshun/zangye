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
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	ParentID    *string `json:"parent_id"`
	Description *string `json:"description"`
	SortOrder   int     `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// FolderCreateRequest 是创建文件夹的请求体。
type FolderCreateRequest struct {
	Name        string  `json:"name"`
	Icon        string  `json:"icon"`
	ParentID    *string `json:"parent_id"`
	Description *string `json:"description"`
	SortOrder   int     `json:"sort_order"`
}

// FolderUpdateRequest 是更新文件夹的请求体。
type FolderUpdateRequest struct {
	Name        *string `json:"name"`
	Icon        *string `json:"icon"`
	ParentID    *string `json:"parent_id"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
}

// List GET /api/folders — 获取文件夹列表
// 支持 ?parent_id=xxx 查询子文件夹，不传则返回根目录下的文件夹
func (h *FoldersHandler) List(w http.ResponseWriter, r *http.Request) {
	parentID := r.URL.Query().Get("parent_id")

	var rows *sql.Rows
	var err error
	if parentID == "" {
		rows, err = h.DB.Query(
			"SELECT id, name, icon, parent_id, description, sort_order, created_at, updated_at FROM folders WHERE parent_id IS NULL ORDER BY sort_order, name",
		)
	} else {
		rows, err = h.DB.Query(
			"SELECT id, name, icon, parent_id, description, sort_order, created_at, updated_at FROM folders WHERE parent_id = ? ORDER BY sort_order, name",
			parentID,
		)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询文件夹列表失败")
		return
	}
	defer closeRows(rows)

	folders := make([]FolderResponse, 0)
	for rows.Next() {
		var f FolderResponse
		if err := rows.Scan(&f.ID, &f.Name, &f.Icon, &f.ParentID, &f.Description, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt); err != nil {
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
	if !scanRow(w, h.DB.QueryRow(
		"SELECT id, name, icon, parent_id, description, sort_order, created_at, updated_at FROM folders WHERE id = ?",
		id,
	), &f.ID, &f.Name, &f.Icon, &f.ParentID, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt) {
		writeError(w, http.StatusNotFound, "文件夹不存在")
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
		"INSERT INTO folders (id, name, icon, parent_id, description, sort_order) VALUES (?, ?, ?, ?, ?, ?)",
		id, req.Name, icon, req.ParentID, req.Description, req.SortOrder,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "创建文件夹失败")
		return
	}

	var f FolderResponse
	if !mustScanRow(w, h.DB.QueryRow(
		"SELECT id, name, icon, parent_id, description, sort_order, created_at, updated_at FROM folders WHERE id = ?",
		id,
	), &f.ID, &f.Name, &f.Icon, &f.ParentID, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt) {
		return
	}

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
	if req.Description != nil {
		fields += "description = ?, "
		args = append(args, *req.Description)
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
	if !mustScanRow(w, h.DB.QueryRow(
		"SELECT id, name, icon, parent_id, description, sort_order, created_at, updated_at FROM folders WHERE id = ?",
		id,
	), &f.ID, &f.Name, &f.Icon, &f.ParentID, &f.SortOrder, &f.CreatedAt, &f.UpdatedAt) {
		return
	}

	writeJSON(w, http.StatusOK, f)
}

// FolderStatsResponse 是文件夹统计信息的响应结构体。
type FolderStatsResponse struct {
	FolderCount  int   `json:"folder_count"`
	FileCount    int   `json:"file_count"`
	TotalSize    int64 `json:"total_size"`
}

// Stats GET /api/folders/{id}/stats — 获取文件夹统计（递归子文件夹）
func (h *FoldersHandler) Stats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// 确认文件夹存在
	var exists int
	if err := h.DB.QueryRow("SELECT 1 FROM folders WHERE id = ?", id).Scan(&exists); err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "文件夹不存在")
		return
	}

	// 递归收集所有子文件夹 ID
	allFolderIDs := []string{id}
	queue := []string{id}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		rows, err := h.DB.Query("SELECT id FROM folders WHERE parent_id = ?", current)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "查询子文件夹失败")
			return
		}
		for rows.Next() {
			var subID string
			err = rows.Scan(&subID)
			if err != nil {
				rows.Close()
				writeError(w, http.StatusInternalServerError, "扫描子文件夹ID失败")
				return
			}
			allFolderIDs = append(allFolderIDs, subID)
			queue = append(queue, subID)
		}
		closeRows(rows)
	}

	// 统计这些文件夹下的文件数和总大小
	var stats FolderStatsResponse
	stats.FolderCount = len(allFolderIDs) - 1 // 减去自身

	// 用 IN 查询所有文件
	if len(allFolderIDs) > 0 {
		query, args := buildInQuery(
			"SELECT COUNT(*), COALESCE(SUM(file_size), 0) FROM files WHERE folder_id IN",
			allFolderIDs,
		)
		if !scanRow(w, h.DB.QueryRow(query, args...), &stats.FileCount, &stats.TotalSize) {
			return
		}
	}

	writeJSON(w, http.StatusOK, stats)
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
