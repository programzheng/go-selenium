package main

import (
	"github.com/gin-gonic/gin"
	"github.com/programzheng/go-selenium/internal/screenshot"
)

func main() {
	// 初始化 Gin 引擎
	router := gin.Default()

	// 設置API路由
	apiGroup := router.Group("/api/v1")
	{
		apiGroup.POST("/screenshot", screenshot.ScreenshotByURL)
	}

	router.GET("/screenshot/preview/:filename", screenshot.PreviewScreenshot)
	router.GET("/screenshot/download/:filename", screenshot.DownloadScreenshot)

	// 監聽端口
	router.Run(":80")
}
