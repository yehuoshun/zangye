package file
import "time"

type Collection struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Icon      string     `json:"icon"`
	ParentID  *string    `json:"parent_id"`
	SortOrder int        `json:"sort_order"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
type CollectionPath struct {
	ID           string `json:"id"`
	CollectionID string `json:"collection_id"`
	Path         string `json:"path"`
	AutoScan     bool   `json:"auto_scan"`
	SortOrder    int    `json:"sort_order"`
}
type VirtualFile struct {
	ID           string  `json:"id"`
	CollectionID string  `json:"collection_id"`
	Path         string  `json:"path"`
	DisplayName  *string `json:"display_name"`
	Size         int64   `json:"size"`
	SortOrder    int     `json:"sort_order"`
	CreatedAt    string  `json:"created_at"`
}
type ScanCache struct {
	ID               int64  `json:"id"`
	CollectionPathID string `json:"collection_path_id"`
	FilePath         string `json:"file_path"`
	FileName         string `json:"file_name"`
	FileSize         int64  `json:"file_size"`
	ModTime          string `json:"mod_time"`
	MimeType         string `json:"mime_type"`
	CachedAt         string `json:"cached_at"`
}
type BrowseItem struct {
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Size        int64   `json:"size"`
	ModTime     string  `json:"mod_time"`
	MimeType    string  `json:"mime_type"`
	Source      string  `json:"source"`
	DisplayName *string `json:"display_name"`
	IsDir       bool    `json:"is_dir"`
}
type Tag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
type OverviewStats struct {
	TotalFiles      int64           `json:"total_files"`
	TotalSize       int64           `json:"total_size"`
	Categories      []CategoryStats `json:"categories"`
	CollectionCount int             `json:"collection_count"`
}
type CategoryStats struct {
	Name    string  `json:"name"`
	Label   string  `json:"label"`
	Count   int64   `json:"count"`
	Size    int64   `json:"size"`
	Percent float64 `json:"percent"`
}
