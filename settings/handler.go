package settings

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Handler struct{ DB *sql.DB }

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/settings/{key}", h.getSetting)
	mux.HandleFunc("PUT /api/settings/{key}", h.setSetting)
	mux.HandleFunc("GET /api/prefixes", h.listPrefixes)
	mux.HandleFunc("PUT /api/prefixes/{prefix}", h.updatePrefix)
}

func jsonResp(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) getSetting(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	var value string
	err := h.DB.QueryRow(`SELECT value FROM settings WHERE key=?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		jsonResp(w, 404, map[string]string{"error": "设置不存在"})
		return
	}
	if err != nil { jsonResp(w, 500, map[string]string{"error": err.Error()}); return }
	jsonResp(w, 200, map[string]string{"key": key, "value": value})
}

func (h *Handler) setSetting(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	var body struct{ Value string `json:"value"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonResp(w, 400, map[string]string{"error": "无效的请求体"}); return
	}
	_, err := h.DB.Exec(`INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value=?`,
		key, body.Value, body.Value)
	if err != nil { jsonResp(w, 500, map[string]string{"error": err.Error()}); return }
	jsonResp(w, 200, map[string]string{"status": "ok"})
}

type PrefixConfig struct {
	Prefix      string  `json:"prefix"`
	Type        string  `json:"type"`
	MapPath     *string `json:"map_path"`
	URLTemplate *string `json:"url_template"`
}

func (h *Handler) listPrefixes(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT prefix, type, map_path, url_template FROM prefix_config`)
	if err != nil { jsonResp(w, 500, map[string]string{"error": err.Error()}); return }
	defer rows.Close()
	var items []PrefixConfig
	for rows.Next() {
		var p PrefixConfig
		if err := rows.Scan(&p.Prefix, &p.Type, &p.MapPath, &p.URLTemplate); err != nil { continue }
		items = append(items, p)
	}
	jsonResp(w, 200, items)
}

func (h *Handler) updatePrefix(w http.ResponseWriter, r *http.Request) {
	prefix := r.PathValue("prefix")
	var body PrefixConfig
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonResp(w, 400, map[string]string{"error": "无效的请求体"}); return
	}
	_, err := h.DB.Exec(`UPDATE prefix_config SET type=?, map_path=?, url_template=? WHERE prefix=?`,
		body.Type, body.MapPath, body.URLTemplate, prefix)
	if err != nil { jsonResp(w, 500, map[string]string{"error": err.Error()}); return }
	jsonResp(w, 200, map[string]string{"status": "ok"})
}
