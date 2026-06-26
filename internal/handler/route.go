// Package handler 提供 HTTP 请求处理器
package handler

import (
	"database/sql"
	"net/http"
	"zangye/internal/middleware"
	"zangye/internal/repository"
	"zangye/internal/service"
)

// SetupRoutes 注册所有 API 路由
// Go 1.22 的 net/http 支持路径参数 {id}，类比 Spring 的 @PathVariable
// 所有路由前缀 /api
func SetupRoutes(db *sql.DB) http.Handler {
	// 初始化各层
	// Go 的依赖注入：手动组装各层依赖
	// 类比 Spring 的 @Autowired，但更显式

	// Repository 层
	folderRepo := repository.NewFolderRepo(db)
	fileRepo := repository.NewFileRepo(db)
	tagRepo := repository.NewTagRepo(db)
	fileTagRepo := repository.NewFileTagRepo(db)
	settingsRepo := repository.NewSettingsRepo(db)

	// Service 层
	folderSvc := service.NewFolderSvc(folderRepo, fileRepo, db)
	fileSvc := service.NewFileSvc(fileRepo, folderRepo, tagRepo, fileTagRepo, db)
	tagSvc := service.NewTagSvc(tagRepo, fileTagRepo)
	dashboardSvc := service.NewDashboardSvc(folderRepo, fileRepo)

	// Handler 层
	folderHandler := NewFolderHandler(folderSvc)
	fileHandler := NewFileHandler(fileSvc)
	tagHandler := NewTagHandler(tagSvc)
	settingsHandler := NewSettingsHandler(settingsRepo)
	dashboardHandler := NewDashboardHandler(dashboardSvc)

	// 创建路由 mux
	// Go 1.22 的 http.NewServeMux() 支持方法前缀和路径参数
	mux := http.NewServeMux()

	// 健康检查
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok", "service": "zangye"})
	})

	// 仪表盘
	mux.HandleFunc("GET /api/dashboard/stats", dashboardHandler.HandleStats)

	// 设置
	mux.HandleFunc("GET /api/settings", settingsHandler.HandleSettings)
	mux.HandleFunc("PUT /api/settings", settingsHandler.HandleSettings)

	// 文件夹
	mux.HandleFunc("GET /api/folders", folderHandler.HandleFolders)
	mux.HandleFunc("POST /api/folders", folderHandler.HandleFolders)
	mux.HandleFunc("GET /api/folders/{id}", folderHandler.HandleFolder)
	mux.HandleFunc("PUT /api/folders/{id}", folderHandler.HandleFolder)
	mux.HandleFunc("DELETE /api/folders/{id}", folderHandler.HandleFolder)
	mux.HandleFunc("GET /api/folders/{id}/stats", folderHandler.HandleFolderStats)

	// 文件
	mux.HandleFunc("GET /api/files", fileHandler.HandleFiles)
	mux.HandleFunc("POST /api/files", fileHandler.HandleFiles)
	mux.HandleFunc("GET /api/files/{id}", fileHandler.HandleFile)
	mux.HandleFunc("PUT /api/files/{id}", fileHandler.HandleFile)
	mux.HandleFunc("DELETE /api/files/{id}", fileHandler.HandleFile)
	mux.HandleFunc("GET /api/files/{id}/content", fileHandler.HandleFileContent)
	mux.HandleFunc("GET /api/files/{id}/tags", fileHandler.HandleFileTags)
	mux.HandleFunc("PUT /api/files/{id}/tags", fileHandler.HandleFileTags)
	mux.HandleFunc("POST /api/files/{id}/tags", fileHandler.HandleFileTags)
	mux.HandleFunc("DELETE /api/files/{id}/tags", fileHandler.HandleFileTags)

	// 标签
	mux.HandleFunc("GET /api/tags", tagHandler.HandleTags)
	mux.HandleFunc("POST /api/tags", tagHandler.HandleTags)
	mux.HandleFunc("GET /api/tags/{id}", tagHandler.HandleTag)
	mux.HandleFunc("PUT /api/tags/{id}", tagHandler.HandleTag)
	mux.HandleFunc("DELETE /api/tags/{id}", tagHandler.HandleTag)

	// 回收站
	mux.HandleFunc("GET /api/trash/files", fileHandler.HandleFiles) // 复用文件列表，通过 trash 参数区分
	mux.HandleFunc("POST /api/trash/files/{id}/restore", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			writeError(w, http.StatusBadRequest, "缺少文件 ID")
			return
		}
		if err := fileSvc.Restore(id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeOK(w)
	})
	mux.HandleFunc("DELETE /api/trash/files/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			writeError(w, http.StatusBadRequest, "缺少文件 ID")
			return
		}
		if err := fileSvc.HardDelete(id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeOK(w)
	})

	// 应用 CORS 中间件
	return middleware.CORS(mux)
}
