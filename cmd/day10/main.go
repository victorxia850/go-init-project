package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// --- 形式 2 所需的处理器结构体 ---
type StructHandler struct{}

func (h *StructHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "姿势 [2]: 标准库结构化版\n路径: %s\n时间: %v", r.URL.Path, time.Now().Format(time.Kitchen))
}

func main() {
	fmt.Println("🚀 正在同时启动三种 HTTP 服务...")

	// ---------------------------------------------------------
	// 姿势 1: 标准库函数式 (最简实现)
	// 监听端口: 8081
	// ---------------------------------------------------------
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "姿势 [1]: 标准库函数式 (Hello World)")
		})
		fmt.Println("✅ [1] 函数式已启动: http://localhost:8081")
		http.ListenAndServe(":8081", nil)
	}()

	// ---------------------------------------------------------
	// 姿势 2: 标准库结构化 (高度可控)
	// 监听端口: 8082
	// ---------------------------------------------------------
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/hello", &StructHandler{})

		server := &http.Server{
			Addr:         ":8082",
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		}
		fmt.Println("✅ [2] 结构化版已启动: http://localhost:8082/hello")
		server.ListenAndServe()
	}()

	// ---------------------------------------------------------
	// 姿势 3: Gin 框架版 (工业级 API 常用)
	// 监听端口: 8083
	// ---------------------------------------------------------
	go func() {
		// 设置为发布模式，减少不必要的日志输出
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "姿势 [3]: Gin 框架版 (Hello World)",
				"status":  "ok",
			})
		})
		fmt.Println("✅ [3] Gin 框架版已启动: http://localhost:8083/ping")
		r.Run(":8083")
	}()

	// 防止主协程直接退出
	select {}
}
