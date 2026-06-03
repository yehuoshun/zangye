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
	"io"
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
// frontendFS 嵌入前端构建产物（dist 目录），编译后无需外部文件。
// 前端通过 npm run build 构建到 frontend/dist 后，
// Go 编译器会将整个目录打包进二进制文件。
var frontendFS embed.FS

func main() {
	// ========================================
	// 1. 初始化数据库连接
	// ========================================
	// 使用默认配置连接 MySQL，自动执行建表语句。
	// 如果数据库无法连接，程序将直接退出（Fatal）。
	database, err := db.New(db.DefaultConfig())
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.Close() // 确保程序退出时关闭数据库连接

	// ========================================
	// 2. 注册 HTTP 路由
	// ========================================

	// 创建仪表盘处理器，注入数据库连接
	dashboardH := &handler.DashboardHandler{DB: database}

	mux := http.NewServeMux()

	// GET /api/dashboard/stats — 获取仪表盘统计数据（文件数、集合数、标签数、存储空间）
	mux.HandleFunc("GET /api/dashboard/stats", dashboardH.Stats)

	// GET /api/health — 健康检查端点，返回服务状态
	// 会实际 ping 数据库以确认连接正常
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
	// 3. 挂载前端 SPA（如果已编译）
	// ========================================
	// 从嵌入的 frontend/dist 中提取前端文件系统。
	// 如果前端未编译（dist 不存在），则仅提供 API 服务。
	frontendSub, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("前端未嵌入，仅 API 可用: %v", err)
	} else {
		// SPA 处理器：所有非 API 请求回退到 index.html
		// 支持 Vue Router 的 history 模式
		mux.Handle("/", &spaHandler{staticFS: frontendSub})
	}

	// ========================================
	// 4. 配置并启动 HTTP 服务器
	// ========================================

	// 监听地址：默认 127.0.0.1:27138，可通过 ZANGYE_ADDR 环境变量覆盖
	addr := "127.0.0.1:27138"
	if env := os.Getenv("ZANGYE_ADDR"); env != "" {
		addr = env
	}

	// 配置服务器超时参数，防止慢客户端占用资源
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second, // 读取请求超时
		WriteTimeout: 10 * time.Second, // 写入响应超时
	}

	// ========================================
	// 5. 优雅关闭
	// ========================================
	// 监听 SIGINT（Ctrl+C）和 SIGTERM 信号，
	// 收到信号后优雅关闭服务器，等待当前请求处理完毕。
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh                       // 阻塞等待信号
		log.Println("正在关闭…")
		server.Shutdown(context.Background()) // 优雅关闭
	}()

	log.Printf("🦞 藏叶 启动 → http://%s", addr)
	// ListenAndServe 在 Shutdown 后返回 ErrServerClosed，
	// 这是正常关闭，不应视为错误
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// spaHandler 是 SPA（单页应用）回退处理器。
// 对于静态文件请求，直接返回对应文件；
// 对于不存在的路径（如 Vue Router 路由），回退到 index.html。
// 这确保了前端路由刷新后不会出现 404。
type spaHandler struct {
	staticFS fs.FS // 嵌入的前端文件系统
}

// ServeHTTP 实现 http.Handler 接口。
// 逻辑：
//  1. 尝试打开请求路径对应的静态文件，如果存在则直接返回
//  2. 如果不存在（如 /dashboard 这样的前端路由），回退到 index.html
//     让 Vue Router 接管路由解析
func (s *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 尝试打开请求路径对应的文件
	f, err := s.staticFS.Open(r.URL.Path)
	if err == nil {
		// 文件存在，使用标准文件服务器返回
		f.Close()
		http.FileServer(http.FS(s.staticFS)).ServeHTTP(w, r)
		return
	}

	// 文件不存在，回退到 index.html（SPA 入口）
	index, err := s.staticFS.Open("index.html")
	if err != nil {
		http.Error(w, "前端未编译", http.StatusNotFound)
		return
	}
	defer index.Close()

	// 使用 ServeContent 返回 index.html，支持 Range 请求和缓存
	stat, _ := index.Stat()
	http.ServeContent(w, r, "index.html", stat.ModTime(), index.(io.ReadSeeker))
}