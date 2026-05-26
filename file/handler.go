package file

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Handler struct{ DB *sql.DB }

func newID() string { b := make([]byte, 8); rand.Read(b); return hex.EncodeToString(b) }

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/collections", h.listRootCollections)
	mux.HandleFunc("GET /api/collections/{id}", h.getCollection)
	mux.HandleFunc("GET /api/collections/{id}/children", h.listChildCollections)
	mux.HandleFunc("POST /api/collections", h.createCollection)
	mux.HandleFunc("PUT /api/collections/{id}", h.updateCollection)
	mux.HandleFunc("DELETE /api/collections/{id}", h.deleteCollection)
	mux.HandleFunc("PUT /api/collections/reorder", h.reorderCollections)
	mux.HandleFunc("GET /api/collections/{id}/paths", h.listPaths)
	mux.HandleFunc("POST /api/collections/{id}/paths", h.createPath)
	mux.HandleFunc("PUT /api/paths/{id}", h.updatePath)
	mux.HandleFunc("DELETE /api/paths/{id}", h.deletePath)
	mux.HandleFunc("POST /api/paths/{id}/scan", h.scanPath)
	mux.HandleFunc("GET /api/collections/{id}/vfiles", h.listVirtualFiles)
	mux.HandleFunc("POST /api/collections/{id}/vfiles", h.createVirtualFile)
	mux.HandleFunc("PUT /api/vfiles/{id}", h.updateVirtualFile)
	mux.HandleFunc("DELETE /api/vfiles/{id}", h.deleteVirtualFile)
	mux.HandleFunc("GET /api/collections/{id}/browse", h.browseCollection)
	mux.HandleFunc("GET /api/preview/content", h.previewContent)
	mux.HandleFunc("GET /api/preview/thumbnail", h.previewThumbnail)
	mux.HandleFunc("POST /api/open-external", h.openExternal)
	mux.HandleFunc("GET /api/tags/search", h.searchTags)
	mux.HandleFunc("POST /api/tags", h.createTag)
	mux.HandleFunc("DELETE /api/tags/{id}", h.deleteTag)
	mux.HandleFunc("GET /api/files/{id}/tags", h.getFileTags)
	mux.HandleFunc("PUT /api/files/{id}/tags", h.updateFileTags)
	mux.HandleFunc("GET /api/collections/{id}/tags", h.getCollectionTags)
	mux.HandleFunc("PUT /api/collections/{id}/tags", h.updateCollectionTags)
	mux.HandleFunc("GET /api/overview/stats", h.overviewStats)
}

func jsonResp(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func jsonErr(w http.ResponseWriter, status int, msg string) {
	jsonResp(w, status, map[string]string{"error": msg})
}

func (h *Handler) listRootCollections(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, name, icon, parent_id, sort_order, created_at, updated_at FROM collections WHERE parent_id IS NULL ORDER BY sort_order, created_at`)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	defer rows.Close()
	var items []Collection
	for rows.Next() {
		var c Collection; var pid sql.NullString; var ca, ua string
		if err := rows.Scan(&c.ID, &c.Name, &c.Icon, &pid, &c.SortOrder, &ca, &ua); err != nil { jsonErr(w, 500, err.Error()); return }
		if pid.Valid { c.ParentID = &pid.String }
		c.CreatedAt, _ = time.Parse(time.RFC3339, ca); c.UpdatedAt, _ = time.Parse(time.RFC3339, ua)
		items = append(items, c)
	}
	jsonResp(w, 200, items)
}

func (h *Handler) getCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var c Collection; var pid sql.NullString; var ca, ua string
	err := h.DB.QueryRow(`SELECT id, name, icon, parent_id, sort_order, created_at, updated_at FROM collections WHERE id=?`, id).Scan(&c.ID, &c.Name, &c.Icon, &pid, &c.SortOrder, &ca, &ua)
	if err == sql.ErrNoRows { jsonErr(w, 404, "文件夹不存在"); return }
	if err != nil { jsonErr(w, 500, err.Error()); return }
	if pid.Valid { c.ParentID = &pid.String }
	c.CreatedAt, _ = time.Parse(time.RFC3339, ca); c.UpdatedAt, _ = time.Parse(time.RFC3339, ua)
	jsonResp(w, 200, c)
}

func (h *Handler) listChildCollections(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rows, err := h.DB.Query(`SELECT id, name, icon, parent_id, sort_order, created_at, updated_at FROM collections WHERE parent_id=? ORDER BY sort_order, created_at`, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	defer rows.Close()
	var items []Collection
	for rows.Next() {
		var c Collection; var pid sql.NullString; var ca, ua string
		if err := rows.Scan(&c.ID, &c.Name, &c.Icon, &pid, &c.SortOrder, &ca, &ua); err != nil { jsonErr(w, 500, err.Error()); return }
		if pid.Valid { c.ParentID = &pid.String }
		c.CreatedAt, _ = time.Parse(time.RFC3339, ca); c.UpdatedAt, _ = time.Parse(time.RFC3339, ua)
		items = append(items, c)
	}
	jsonResp(w, 200, items)
}

func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	var c Collection
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	c.ID = newID(); if c.Icon == "" { c.Icon = "📁" }
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := h.DB.Exec(`INSERT INTO collections (id, name, icon, parent_id, sort_order, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, c.ID, c.Name, c.Icon, c.ParentID, c.SortOrder, now, now)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	c.CreatedAt, _ = time.Parse(time.RFC3339, now); c.UpdatedAt, _ = time.Parse(time.RFC3339, now)
	jsonResp(w, 201, c)
}

func (h *Handler) updateCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var c Collection
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := h.DB.Exec(`UPDATE collections SET name=?, icon=?, parent_id=?, sort_order=?, updated_at=? WHERE id=?`, c.Name, c.Icon, c.ParentID, c.SortOrder, now, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) deleteCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := h.DB.Exec(`DELETE FROM collections WHERE id=?`, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) reorderCollections(w http.ResponseWriter, r *http.Request) {
	var order []struct {
		ID        string `json:"id"`
		SortOrder int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	tx, err := h.DB.Begin()
	if err != nil { jsonErr(w, 500, err.Error()); return }
	for _, o := range order {
		if _, err := tx.Exec(`UPDATE collections SET sort_order=? WHERE id=?`, o.SortOrder, o.ID); err != nil { tx.Rollback(); jsonErr(w, 500, err.Error()); return }
	}
	tx.Commit()
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) listPaths(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rows, err := h.DB.Query(`SELECT id, collection_id, path, auto_scan, sort_order FROM collection_paths WHERE collection_id=? ORDER BY sort_order`, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	defer rows.Close()
	var items []CollectionPath
	for rows.Next() {
		var p CollectionPath; var autoScan int
		if err := rows.Scan(&p.ID, &p.CollectionID, &p.Path, &autoScan, &p.SortOrder); err != nil { jsonErr(w, 500, err.Error()); return }
		p.AutoScan = autoScan == 1; items = append(items, p)
	}
	jsonResp(w, 200, items)
}

func (h *Handler) createPath(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var p CollectionPath
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	p.ID = newID(); p.CollectionID = id
	autoScan := 0; if p.AutoScan { autoScan = 1 }
	_, err := h.DB.Exec(`INSERT INTO collection_paths (id, collection_id, path, auto_scan, sort_order) VALUES (?, ?, ?, ?, ?)`, p.ID, p.CollectionID, p.Path, autoScan, p.SortOrder)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	jsonResp(w, 201, p)
}

func (h *Handler) updatePath(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var p CollectionPath
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	autoScan := 0; if p.AutoScan { autoScan = 1 }
	_, err := h.DB.Exec(`UPDATE collection_paths SET path=?, auto_scan=?, sort_order=? WHERE id=?`, p.Path, autoScan, p.SortOrder, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) deletePath(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := h.DB.Exec(`DELETE FROM collection_paths WHERE id=?`, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) scanPath(w http.ResponseWriter, r *http.Request) {
	pathID := r.PathValue("id")
	var cp CollectionPath; var autoScan int
	err := h.DB.QueryRow(`SELECT id, collection_id, path, auto_scan FROM collection_paths WHERE id=?`, pathID).Scan(&cp.ID, &cp.CollectionID, &cp.Path, &autoScan)
	if err != nil { jsonErr(w, 404, "路径不存在"); return }
	realPath := h.resolvePath(cp.Path)
	if realPath == "" { jsonErr(w, 400, "无法解析路径前缀"); return }
	info, err := os.Stat(realPath)
	if err != nil { jsonErr(w, 404, fmt.Sprintf("路径不可访问: %v", err)); return }
	if !info.IsDir() { jsonErr(w, 400, "路径不是目录"); return }
	h.DB.Exec(`DELETE FROM scan_cache WHERE collection_path_id=?`, pathID)
	count := 0
	filepath.WalkDir(realPath, func(fp string, d fs.DirEntry, err error) error {
		if err != nil { return nil }
		if d.IsDir() && fp != realPath { return filepath.SkipDir }
		if d.IsDir() { return nil }
		info, err := d.Info()
		if err != nil { return nil }
		rel, _ := filepath.Rel(realPath, fp)
		mime := detectMime(fp)
		_, err = h.DB.Exec(`INSERT OR REPLACE INTO scan_cache (collection_path_id, file_path, file_name, file_size, mod_time, mime_type) VALUES (?, ?, ?, ?, ?, ?)`, pathID, rel, d.Name(), info.Size(), info.ModTime().UTC().Format(time.RFC3339), mime)
		if err == nil { count++ }
		return nil
	})
	jsonResp(w, 200, map[string]any{"status": "ok", "count": count})
}

func (h *Handler) listVirtualFiles(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rows, err := h.DB.Query(`SELECT id, collection_id, path, display_name, size, sort_order, created_at FROM virtual_files WHERE collection_id=? ORDER BY sort_order, created_at`, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	defer rows.Close()
	var items []VirtualFile
	for rows.Next() {
		var v VirtualFile
		if err := rows.Scan(&v.ID, &v.CollectionID, &v.Path, &v.DisplayName, &v.Size, &v.SortOrder, &v.CreatedAt); err != nil { jsonErr(w, 500, err.Error()); return }
		items = append(items, v)
	}
	jsonResp(w, 200, items)
}

func (h *Handler) createVirtualFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var v VirtualFile
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	v.ID = newID(); v.CollectionID = id
	if v.CreatedAt == "" { v.CreatedAt = time.Now().UTC().Format(time.RFC3339) }
	_, err := h.DB.Exec(`INSERT INTO virtual_files (id, collection_id, path, display_name, size, sort_order, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, v.ID, v.CollectionID, v.Path, v.DisplayName, v.Size, v.SortOrder, v.CreatedAt)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	jsonResp(w, 201, v)
}

func (h *Handler) updateVirtualFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var v VirtualFile
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	_, err := h.DB.Exec(`UPDATE virtual_files SET path=?, display_name=?, size=?, sort_order=? WHERE id=?`, v.Path, v.DisplayName, v.Size, v.SortOrder, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) deleteVirtualFile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := h.DB.Exec(`DELETE FROM virtual_files WHERE id=?`, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) browseCollection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort")
	var items []BrowseItem
	paths, _ := h.listPathsInternal(id)
	for _, p := range paths {
		source := extractSource(p.Path)
		rows, err := h.DB.Query(`SELECT file_path, file_name, file_size, mod_time, mime_type FROM scan_cache WHERE collection_path_id=?`, p.ID)
		if err != nil { continue }
		for rows.Next() {
			var item BrowseItem
			if err := rows.Scan(&item.Path, &item.Name, &item.Size, &item.ModTime, &item.MimeType); err != nil { continue }
			item.Type = "file"; item.Source = source; item.IsDir = false; items = append(items, item)
		}
		rows.Close()
	}
	vfiles, _ := h.listVirtualFilesInternal(id)
	for _, v := range vfiles {
		item := BrowseItem{Type: "virtual", Name: v.Path, Path: v.Path, Size: v.Size, Source: extractSource(v.Path), DisplayName: v.DisplayName, IsDir: false}
		if v.DisplayName != nil && *v.DisplayName != "" { item.Name = *v.DisplayName }
		items = append(items, item)
	}
	for _, p := range paths {
		source := extractSource(p.Path)
		prefix := extractPrefix(p.Path)
		if prefixType(prefix, h.DB) == "web" {
			items = append(items, BrowseItem{Type: "link", Name: p.Path, Path: p.Path, Source: source, IsDir: false})
		}
	}
	if search != "" {
		filtered := make([]BrowseItem, 0)
		for _, item := range items {
			if strings.Contains(strings.ToLower(item.Name), strings.ToLower(search)) { filtered = append(filtered, item) }
		}
		items = filtered
	}
	switch sortBy {
	case "size": sort.Slice(items, func(i, j int) bool { return items[i].Size > items[j].Size })
	case "time": sort.Slice(items, func(i, j int) bool { return items[i].ModTime > items[j].ModTime })
	default: sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	}
	jsonResp(w, 200, items)
}

func (h *Handler) listPathsInternal(cid string) ([]CollectionPath, error) {
	rows, err := h.DB.Query(`SELECT id, collection_id, path, auto_scan, sort_order FROM collection_paths WHERE collection_id=? ORDER BY sort_order`, cid)
	if err != nil { return nil, err }
	defer rows.Close()
	var items []CollectionPath
	for rows.Next() {
		var p CollectionPath; var autoScan int
		if err := rows.Scan(&p.ID, &p.CollectionID, &p.Path, &autoScan, &p.SortOrder); err != nil { continue }
		p.AutoScan = autoScan == 1; items = append(items, p)
	}
	return items, nil
}

func (h *Handler) listVirtualFilesInternal(cid string) ([]VirtualFile, error) {
	rows, err := h.DB.Query(`SELECT id, collection_id, path, display_name, size, sort_order, created_at FROM virtual_files WHERE collection_id=? ORDER BY sort_order`, cid)
	if err != nil { return nil, err }
	defer rows.Close()
	var items []VirtualFile
	for rows.Next() {
		var v VirtualFile
		if err := rows.Scan(&v.ID, &v.CollectionID, &v.Path, &v.DisplayName, &v.Size, &v.SortOrder, &v.CreatedAt); err != nil { continue }
		items = append(items, v)
	}
	return items, nil
}

func (h *Handler) previewContent(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" { jsonErr(w, 400, "缺少 path 参数"); return }
	realPath := h.resolvePath(path)
	if realPath == "" { jsonErr(w, 400, "无法解析路径"); return }
	ext := strings.ToLower(filepath.Ext(realPath))
	textExts := map[string]bool{".txt": true, ".md": true, ".go": true, ".java": true, ".py": true, ".js": true, ".ts": true, ".vue": true, ".html": true, ".css": true, ".json": true, ".xml": true, ".yaml": true, ".yml": true, ".toml": true, ".sql": true, ".sh": true, ".bat": true, ".ini": true, ".cfg": true, ".log": true, ".csv": true, ".kt": true, ".rs": true, ".c": true, ".h": true}
	if !textExts[ext] { jsonErr(w, 415, "不支持预览此文件类型"); return }
	data, err := os.ReadFile(realPath)
	if err != nil { jsonErr(w, 500, fmt.Sprintf("读取失败: %v", err)); return }
	if len(data) > 1024*1024 { data = data[:1024*1024] }
	mime := detectMime(realPath)
	jsonResp(w, 200, map[string]any{"content": string(data), "mime": mime, "size": len(data)})
}

func (h *Handler) previewThumbnail(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" { jsonErr(w, 400, "缺少 path 参数"); return }
	realPath := h.resolvePath(path)
	if realPath == "" { jsonErr(w, 400, "无法解析路径"); return }
	ext := strings.ToLower(filepath.Ext(realPath))
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".webp" || ext == ".svg" { http.ServeFile(w, r, realPath); return }
	jsonErr(w, 415, "不支持生成此类型的缩略图")
}

func (h *Handler) openExternal(w http.ResponseWriter, r *http.Request) {
	var body struct{ Path string `json:"path"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	realPath := h.resolvePath(body.Path)
	if realPath != "" { openWithSystem(realPath); jsonResp(w, 200, map[string]any{"type": "local", "path": realPath}); return }
	prefix := extractPrefix(body.Path)
	if prefixType(prefix, h.DB) == "web" {
		url := buildURL(body.Path, h.DB)
		if url != "" { openWithSystem(url); jsonResp(w, 200, map[string]any{"type": "web", "url": url}); return }
	}
	jsonErr(w, 400, "无法解析路径")
}

func (h *Handler) searchTags(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	rows, err := h.DB.Query(`SELECT id, name, color FROM tags WHERE name LIKE ? LIMIT 20`, "%"+q+"%")
	if err != nil { jsonErr(w, 500, err.Error()); return }
	defer rows.Close()
	var items []Tag
	for rows.Next() { var t Tag; rows.Scan(&t.ID, &t.Name, &t.Color); items = append(items, t) }
	jsonResp(w, 200, items)
}

func (h *Handler) createTag(w http.ResponseWriter, r *http.Request) {
	var t Tag
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	t.ID = newID(); if t.Color == "" { t.Color = "gray" }
	_, err := h.DB.Exec(`INSERT INTO tags (id, name, color) VALUES (?, ?, ?)`, t.ID, t.Name, t.Color)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	jsonResp(w, 201, t)
}

func (h *Handler) deleteTag(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_, err := h.DB.Exec(`DELETE FROM tags WHERE id=?`, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) getFileTags(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rows, err := h.DB.Query(`SELECT t.id, t.name, t.color FROM tags t JOIN file_tags ft ON t.id = ft.tag_id WHERE ft.file_id=?`, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	defer rows.Close()
	var items []Tag
	for rows.Next() { var t Tag; rows.Scan(&t.ID, &t.Name, &t.Color); items = append(items, t) }
	jsonResp(w, 200, items)
}

func (h *Handler) updateFileTags(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var tagIDs []string
	if err := json.NewDecoder(r.Body).Decode(&tagIDs); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	tx, _ := h.DB.Begin()
	tx.Exec(`DELETE FROM file_tags WHERE file_id=?`, id)
	for _, tid := range tagIDs { tx.Exec(`INSERT OR IGNORE INTO file_tags (file_id, tag_id) VALUES (?, ?)`, id, tid) }
	tx.Commit()
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) getCollectionTags(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rows, err := h.DB.Query(`SELECT t.id, t.name, t.color FROM tags t JOIN collection_tags ct ON t.id = ct.tag_id WHERE ct.collection_id=?`, id)
	if err != nil { jsonErr(w, 500, err.Error()); return }
	defer rows.Close()
	var items []Tag
	for rows.Next() { var t Tag; rows.Scan(&t.ID, &t.Name, &t.Color); items = append(items, t) }
	jsonResp(w, 200, items)
}

func (h *Handler) updateCollectionTags(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var tagIDs []string
	if err := json.NewDecoder(r.Body).Decode(&tagIDs); err != nil { jsonErr(w, 400, "无效的请求体"); return }
	tx, _ := h.DB.Begin()
	tx.Exec(`DELETE FROM collection_tags WHERE collection_id=?`, id)
	for _, tid := range tagIDs { tx.Exec(`INSERT OR IGNORE INTO collection_tags (collection_id, tag_id) VALUES (?, ?)`, id, tid) }
	tx.Commit()
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) overviewStats(w http.ResponseWriter, r *http.Request) {
	stats := OverviewStats{}
	h.DB.QueryRow(`SELECT COUNT(*) FROM scan_cache`).Scan(&stats.TotalFiles)
	h.DB.QueryRow(`SELECT COALESCE(SUM(file_size), 0) FROM scan_cache`).Scan(&stats.TotalSize)
	h.DB.QueryRow(`SELECT COUNT(*) FROM collections`).Scan(&stats.CollectionCount)
	rows, _ := h.DB.Query(`SELECT CASE WHEN mime_type LIKE 'video/%' THEN 'video' WHEN mime_type LIKE 'audio/%' THEN 'audio' WHEN mime_type LIKE 'image/%' THEN 'image' WHEN mime_type LIKE 'text/%' OR mime_type LIKE 'application/pdf' OR mime_type LIKE 'application/msword' OR mime_type LIKE 'application/vnd.openxmlformats%' THEN 'document' ELSE 'other' END as category, COUNT(*) as cnt, COALESCE(SUM(file_size), 0) as total_size FROM scan_cache GROUP BY category`)
	defer rows.Close()
	categoryLabels := map[string]string{"video": "🎬 视频", "audio": "🎵 音频", "image": "🖼️ 图片", "document": "📄 文档", "other": "📦 其他"}
	for rows.Next() {
		var cs CategoryStats; rows.Scan(&cs.Name, &cs.Count, &cs.Size)
		cs.Label = categoryLabels[cs.Name]
		if stats.TotalSize > 0 { cs.Percent = float64(cs.Size) / float64(stats.TotalSize) * 100 }
		stats.Categories = append(stats.Categories, cs)
	}
	jsonResp(w, 200, stats)
}

// resolvePath 从 DB prefix_config 表查前缀映射
func (h *Handler) resolvePath(path string) string {
	rows, err := h.DB.Query(`SELECT prefix, map_path FROM prefix_config WHERE type='local' AND map_path IS NOT NULL AND map_path != ''`)
	if err != nil { return path }
	defer rows.Close()
	for rows.Next() {
		var prefix, mapPath string
		if err := rows.Scan(&prefix, &mapPath); err != nil { continue }
		if strings.HasPrefix(path, prefix) {
			rel := strings.TrimPrefix(path, prefix)
			rel = strings.TrimPrefix(rel, "\\"); rel = strings.TrimPrefix(rel, "/")
			return filepath.Join(mapPath, rel)
		}
	}
	if strings.HasPrefix(path, "/") || (len(path) >= 2 && path[1] == ':') { return path }
	return path
}

func extractSource(path string) string {
	if strings.HasPrefix(path, "115:\\") { return "115" }
	if strings.HasPrefix(path, "tg:\\") { return "tg" }
	if strings.HasPrefix(path, "notion:\\") { return "notion" }
	return "local"
}

func extractPrefix(path string) string {
	for _, p := range []string{"115:\\", "tg:\\", "notion:\\"} {
		if strings.HasPrefix(path, p) { return p }
	}
	return ""
}

func prefixType(prefix string, db *sql.DB) string {
	var t string; db.QueryRow(`SELECT type FROM prefix_config WHERE prefix=?`, prefix).Scan(&t); return t
}

func buildURL(path string, db *sql.DB) string {
	prefix := extractPrefix(path)
	var template string; db.QueryRow(`SELECT url_template FROM prefix_config WHERE prefix=?`, prefix).Scan(&template)
	if template == "" { return "" }
	rel := strings.TrimPrefix(path, prefix)
	return strings.Replace(template, "{path}", rel, 1)
}

func openWithSystem(target string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows": cmd = exec.Command("cmd", "/c", "start", "", target)
	case "darwin": cmd = exec.Command("open", target)
	default: cmd = exec.Command("xdg-open", target)
	}
	cmd.Start()
}

func detectMime(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	mimeMap := map[string]string{
		".mp4": "video/mp4", ".mkv": "video/x-matroska", ".avi": "video/x-msvideo", ".mov": "video/quicktime",
		".webm": "video/webm", ".flv": "video/x-flv", ".mp3": "audio/mpeg", ".flac": "audio/flac",
		".wav": "audio/wav", ".aac": "audio/aac", ".ogg": "audio/ogg", ".jpg": "image/jpeg", ".jpeg": "image/jpeg",
		".png": "image/png", ".gif": "image/gif", ".webp": "image/webp", ".svg": "image/svg+xml", ".bmp": "image/bmp",
		".pdf": "application/pdf", ".doc": "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls": "application/vnd.ms-excel", ".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt": "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt": "text/plain", ".md": "text/markdown", ".html": "text/html", ".css": "text/css",
		".js": "application/javascript", ".json": "application/json", ".xml": "application/xml",
		".zip": "application/zip", ".rar": "application/x-rar-compressed", ".7z": "application/x-7z-compressed",
		".tar": "application/x-tar", ".gz": "application/gzip", ".go": "text/x-go", ".java": "text/x-java",
		".py": "text/x-python", ".ts": "text/typescript", ".vue": "text/x-vue",
	}
	if m, ok := mimeMap[ext]; ok { return m }
	return "application/octet-stream"
}

var _ = strconv.Itoa
var _ = sort.Ints