// Package middleware 提供 HTTP 中间件
// Go 的中间件模式：接收 http.Handler，返回 http.Handler
// 类比 Java Servlet Filter 或 Spring 的 HandlerInterceptor
package middleware

import "net/http"

// CORS 返回跨域中间件
// 允许前端开发服务器（如 Vite 的 5173 端口）跨域访问后端 API
// 生产环境 go:embed 模式下不需要跨域，但保留此中间件以便开发调试
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 响应头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Range")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Range, Accept-Ranges, Content-Length")

		// 处理预检请求（OPTIONS）
		// 浏览器在跨域请求前会先发 OPTIONS 请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
