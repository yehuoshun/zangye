// Package main 是藏叶（Zangye）个人文件管理器的主入口。
//
// 藏叶是一个单二进制文件管理工具，后端使用 Go + MySQL，
// 前端使用 Vue 3 + Vite 构建后通过 embed 嵌入到二进制中。
// 启动后提供一个完整的 Web 界面，用于管理本地文件集合。
//
// 核心功能：
//   - 文件集合（Collection）管理：树形结构组织文件
//   - 标签（Tag）系统：为文件添加标签便于检索
//   - 仪表盘（Dashboard）：展示文件统计概览
//
// 启动方式：
//   go build -o zangye .
//   ./zangye
//   # 或在开发时
//   go run .
//
// 默认监听 127.0.0.1:27138，可通过环境变量 ZANGYE_ADDR 自定义。
package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yehuoshun/zangye/internal/db"
	"github.com/yehuoshun/zangye/internal/handler"
)

//go:embed frontend/dist
// frontendEmbed 嵌入前端构建产物，作为后备。
// 生产模式：go build 打包时会将 dist 嵌入二进制。
// 开发模式：如果磁盘上存在 frontend/dist/，优先从磁盘读取。
var frontendEmbed embed.FS

func main() {
	// ========================================
	// 1. 初始化数据库连接
	// ========================================
	database, err := db.New(db.DefaultConfig())
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.Close()

	// ========================================
	// 2. 注册 HTTP 路由
	// ========================================

	dashboardH := &handler.DashboardHandler{DB: database}
	settingsH := &handler.SettingsHandler{DB: database}
	filesH := &handler.FilesHandler{DB: database}

	mux := http.NewServeMux()

	// 仪表盘
	mux.HandleFunc("GET /api/dashboard/stats", dashboardH.Stats)
	// 设置
	mux.HandleFunc("GET /api/settings", settingsH.GetAll)
	mux.HandleFunc("PUT /api/settings", settingsH.Update)
	// 文件管理
	mux.HandleFunc("GET /api/files", filesH.List)
	mux.HandleFunc("POST /api/files/preview", filesH.Preview)
	mux.HandleFunc("GET /api/files/{id}", filesH.Get)
	mux.HandleFunc("POST /api/files", filesH.Create)
	mux.HandleFunc("PUT /api/files/{id}", filesH.Update)
	mux.HandleFunc("DELETE /api/files/{id}", filesH.Delete)

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := database.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"db_error"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// ========================================
	// 3. 挂载前端 SPA
	// ========================================
	// 优先使用磁盘上的 frontend/dist/（开发模式），
	// 不存在时回退到 embed（生产模式）。
	var staticFS fs.FS
	if _, err := os.Stat("frontend/dist"); err == nil {
		log.Println("📁 使用磁盘前端文件: frontend/dist/")
		staticFS = os.DirFS("frontend/dist")
	} else if sub, err := fs.Sub(frontendEmbed, "frontend/dist"); err == nil {
		log.Println("📦 使用嵌入前端文件")
		staticFS = sub
	} else {
		log.Println("⚠️  前端未编译，仅 API 可用")
	}

	if staticFS != nil {
		mux.Handle("/", &spaHandler{staticFS: staticFS})
	}

	// ========================================
	// 4. 配置并启动 HTTP 服务器
	// ========================================

	addr := "127.0.0.1:27138"
	if env := os.Getenv("ZANGYE_ADDR"); env != "" {
		addr = env
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("正在关闭…")
		server.Shutdown(context.Background())
	}()

	log.Printf("🦞 藏叶 启动 → http://%s", addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// spaHandler 是 SPA 回退处理器。
type spaHandler struct {
	staticFS fs.FS
}

// ServeHTTP 实现 http.Handler 接口。
func (s *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	if path == "" {
		path = "index.html"
	}

	// 尝试读取文件
	data, err := fs.ReadFile(s.staticFS, path)
	if err == nil {
		// 根据扩展名设置 MIME 类型
		if len(path) >= 3 && path[len(path)-3:] == ".js" {
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		} else if len(path) >= 4 && path[len(path)-4:] == ".css" {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		} else if len(path) >= 5 && path[len(path)-5:] == ".html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		w.Write(data)
		return
	}

	// 回退到 index.html（SPA 路由）
	index, err := fs.ReadFile(s.staticFS, "index.html")
	if err != nil {
		http.Error(w, "前端未编译", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(index)
}