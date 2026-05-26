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
var frontendFS embed.FS

func main() {
	database, err := db.New(db.DefaultConfig())
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.Close()

	dashboardH := &handler.DashboardHandler{DB: database}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/dashboard/stats", dashboardH.Stats)
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := database.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"db_error"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	frontendSub, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("前端未嵌入，仅 API 可用: %v", err)
	} else {
		mux.Handle("/", &spaHandler{staticFS: frontendSub})
	}

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

type spaHandler struct {
	staticFS fs.FS
}

func (s *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := s.staticFS.Open(r.URL.Path)
	if err == nil {
		f.Close()
		http.FileServer(http.FS(s.staticFS)).ServeHTTP(w, r)
		return
	}

	index, err := s.staticFS.Open("index.html")
	if err != nil {
		http.Error(w, "前端未编译", http.StatusNotFound)
		return
	}
	defer index.Close()

	stat, _ := index.Stat()
	http.ServeContent(w, r, "index.html", stat.ModTime(), index.(io.ReadSeeker))
}