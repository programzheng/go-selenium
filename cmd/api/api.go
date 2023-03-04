package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/programzheng/go-selenium/config"
	"github.com/programzheng/go-selenium/internal/screenshot"
)

func main() {
	// 初始化 Gin 引擎
	router := gin.Default()

	// 設置API路由
	apiGroup := router.Group("/api/v1")
	{
		apiGroup.POST("/screenshot", screenshot.SaveScreenshotByURL)
	}

	router.GET("/screenshot/dynamic/:url", screenshot.DynamicPreviewScreenshot)
	router.GET("/screenshot/static/:filename", screenshot.PreviewByStaticScreenshotFile)
	router.GET("/screenshot/download/:filename", screenshot.DownloadStaticScreenshotFile)

	// 監聽端口
	port := config.Cfg.GetString("PORT")
	if port == "" {
		port = "80"
	}

	router.Run(fmt.Sprintf(":%s", port))
}
