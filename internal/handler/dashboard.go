package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type DashboardHandler struct {
	DB *sql.DB
}

type DashboardStats struct {
	FileCount       int64  `json:"file_count"`
	CollectionCount int64  `json:"collection_count"`
	TagCount        int64  `json:"tag_count"`
	StorageBytes    int64  `json:"storage_bytes"`
	StorageDisplay  string `json:"storage_display"`
}

func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {
	var stats DashboardStats

	h.DB.QueryRow("SELECT COUNT(*) FROM files").Scan(&stats.FileCount)
	h.DB.QueryRow("SELECT COUNT(*) FROM collections").Scan(&stats.CollectionCount)
	h.DB.QueryRow("SELECT COUNT(*) FROM tags").Scan(&stats.TagCount)
	h.DB.QueryRow("SELECT COALESCE(SUM(file_size), 0) FROM files").Scan(&stats.StorageBytes)
	stats.StorageDisplay = formatSize(stats.StorageBytes)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func formatSize(bytes int64) string {
	switch {
	case bytes >= 1<<40:
		return fmt.Sprintf("%.1f TB", float64(bytes)/(1<<40))
	case bytes >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(bytes)/(1<<30))
	case bytes >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(bytes)/(1<<10))
	default:
		return "0 B"
	}
}