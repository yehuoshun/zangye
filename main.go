// 藏叶 (ZangYe) 虚拟文件管理器 - 主入口
//
// 项目定位：虚拟文件管理器（纯手工录入，非磁盘扫描）
// 管理用户在互联网上的备份文件清单
//
// 技术栈：
//   - 后端：Go 1.22+ (net/http 标准库，无第三方框架)
//   - 数据库：MySQL 8.0 (go-sql-driver/mysql)
//   - 前端：Vue 3 + TypeScript + Vite 6 (通过 go:embed 嵌入)
//   - 部署：单 exe（go:embed 前端 dist）
//
// Go 启动要点：
//   - main 函数是程序入口，类比 Java 的 public static void main
//   - init() 函数在 main 之前自动执行，类比 Java 的静态初始化块
//   - defer 在函数返回前执行，类比 Java 的 finally
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
	"zangye/internal/config"
	"zangye/internal/db"
	"zangye/internal/handler"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库连接
	database, err := db.Init(cfg.DSN)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	// defer ≈ Java 的 finally，确保程序退出时关闭数据库连接
	defer database.Close()

	// 设置路由
	apiHandler := handler.SetupRoutes(database)

	// 设置静态文件服务（go:embed 嵌入的前端 dist）
	// 尝试加载嵌入的前端文件
	staticFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("前端静态文件未嵌入（开发模式），仅提供 API 服务: %v", err)
		staticFS = nil
	}

	// 创建主 HTTP mux
	mux := http.NewServeMux()

	// API 路由
	mux.Handle("/api/", apiHandler)

	// 静态文件路由（如果前端已构建）
	if staticFS != nil {
		fileServer := http.FileServer(http.FS(staticFS))
		mux.Handle("/", fileServer)
	} else {
		// 开发模式：返回提示
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(`<h1>藏叶虚拟文件管理器</h1>
<p>API 服务已启动。</p>
<p>前端开发模式：请在 frontend/ 目录下运行 <code>npm run dev</code></p>`))
		})
	}

	// 创建 HTTP 服务器
	// Go 的 http.Server 类比 Java 的 Tomcat/Undertow
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 在 goroutine 中启动服务器
	// Go 的 go 关键字 ≈ Java 的 new Thread().start()，但更轻量
	go func() {
		log.Printf("藏叶虚拟文件管理器启动中...")
		log.Printf("监听地址: http://%s", cfg.Addr)
		log.Printf("API 文档: http://%s/api/health", cfg.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

	// 自动打开浏览器（Windows 平台）
	// 使用 exec.Command 调用 cmd /c start
	// 类比 Java 的 Desktop.getDesktop().browse()
	// 注意：在 Linux 服务器上不会执行
	// exec.Command("cmd", "/c", "start", "http://"+cfg.Addr).Start()

	// 等待中断信号实现优雅关闭
	// Go 的 signal.Notify 捕获系统信号
	// SIGINT: Ctrl+C, SIGTERM: kill 命令
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞直到收到信号

	log.Println("正在关闭服务器...")

	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 优雅关闭：停止接受新请求，等待正在处理的请求完成
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}

	log.Println("服务器已安全关闭")
}
