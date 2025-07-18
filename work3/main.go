package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	router *gin.Engine
	file   *os.File
)

func init() {
	router = gin.Default()

	// 中间件顺序优化：gin.Recovery 放在最前
	router.Use(gin.Recovery(), RequestIDMiddleware(), LoggingMiddleware(), CustomErrorMiddleware())

	// 设置日志级别和格式
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// 打开日志文件
	file, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Info("Failed to log to file, using default stderr")
	}
	logrus.SetOutput(file)
}

func main() {
	defer func() {
		if file != nil {
			if err := file.Close(); err != nil {
				return
			}
		}
	}()
	setupRoutes()
	err := router.Run(":8080")
	if err != nil {
		return
	}
}

// 将路由注册抽离为单独函数，提高可读性
func setupRoutes() {
	router.POST("/register", Register)
	router.POST("/login", Login)
	router.GET("/fileList", FileList)
	router.GET("/fileRead/:postId", FileRead)
	router.GET("/readComment/:postId", ReadComment)

	// 需要鉴权的路由组
	routerGroup := router.Group("/api")
	routerGroup.Use(AuthMiddleware())
	{
		routerGroup.POST("/fileCreated", FileCreated)
		routerGroup.POST("/fileUpdated/:postId", AuthorOnlyMiddleware(), FileUpdated)
		routerGroup.POST("/fileDeleted/:postId", AuthorOnlyMiddleware(), FileDeleted)
		routerGroup.POST("/createComment", CreateComment)
	}
}
