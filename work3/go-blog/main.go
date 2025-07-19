package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-blog/handlers"
	"go-blog/middle"
	"go-blog/utils"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	utils.ConnectDB()
	log.Println("main")
	// 确保数据库连接成功
	if utils.GetDB() == nil {
		log.Fatal("Database connection is nil")
	}
	// 设置 Gin 为调试模式
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.Use(utils.LogRequest)

	r.POST("/register", handlers.Register)

	r.POST("/login", handlers.Login)

	authenticated := r.Group("/")

	authenticated.Use(middle.AuthMiddleware)

	authenticated.POST("/posts", handlers.CreatePost)

	authenticated.GET("/posts", handlers.GetPosts)

	authenticated.GET("/posts/:id", handlers.GetPostByID)

	authenticated.PUT("/posts/:id", handlers.UpdatePost)

	authenticated.DELETE("/posts/:id", handlers.DeletePost)

	authenticated.POST("/comments", handlers.CreateComment)

	authenticated.GET("/posts/:id/comments", handlers.GetCommentsByPostID)

	// 创建HTTP服务器
	serverAddr := fmt.Sprintf(":%s", "8080")
	server := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}
	// 启动服务器
	go func() {
		log.Printf("Server starting on %s", serverAddr)
		utils.LogInfo("Server starting", zap.String("port", "8080"))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	log.Println("Shutting down server...")
	utils.LogInfo("Server shutting down")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// 关闭数据库连接
	if db := utils.GetDB(); db != nil {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}
		}
	}

	log.Println("Server exited")
	utils.LogInfo("Server exited successfully")
}
