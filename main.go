package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/yehuoshun/zangye/db"
	"github.com/yehuoshun/zangye/file"
	"github.com/yehuoshun/zangye/settings"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	listenAddr := "127.0.0.1:27138"
	if addr := os.Getenv("ZANGYE_ADDR"); addr != "" {
		listenAddr = addr
	}
	dataDir := dataDirPath()
	dbPath := filepath.Join(dataDir, "zangye.db")
	database, err := db.New(dbPath)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.Close()
	mux := http.NewServeMux()
	fileHandler := &file.Handler{DB: database}
	fileHandler.RegisterRoutes(mux)
	settingsHandler := &settings.Handler{DB: database}
	settingsHandler.RegisterRoutes(mux)
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	frontendSub, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("前端文件未编译，仅 API 可用: %v", err)
	} else {
		spa := &spaHandler{staticFS: frontendSub}
		mux.Handle("/", spa)
	}
	server := &http.Server{Addr: listenAddr, Handler: corsMiddleware(mux)}
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("正在关闭…")
		server.Close()
	}()
	log.Printf("🦞 藏叶 启动 → http://%s", listenAddr)
	log.Printf("   数据目录: %s", dataDir)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func dataDirPath() string {
	if dir := os.Getenv("ZANGYE_DATA_DIR"); dir != "" {
		return dir
	}
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "zangye")
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", "zangye")
	default:
		home, _ := os.UserHomeDir()
		if home == "" { return "./data" }
		return filepath.Join(home, ".local", "share", "zangye")
	}
}

type spaHandler struct{ staticFS fs.FS }

func (s *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := s.staticFS.Open(r.URL.Path)
	if err == nil { f.Close(); http.FileServer(http.FS(s.staticFS)).ServeHTTP(w, r); return }
	index, err := s.staticFS.Open("index.html")
	if err != nil { http.Error(w, "前端未编译", http.StatusNotFound); return }
	defer index.Close()
	stat, _ := index.Stat()
	http.ServeContent(w, r, "index.html", stat.ModTime(), index.(io.ReadSeeker))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" { w.WriteHeader(204); return }
		next.ServeHTTP(w, r)
	})
}

var _ = embed.FS{}