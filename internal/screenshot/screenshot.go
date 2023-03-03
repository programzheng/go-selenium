package screenshot

import (
	"context"
	"fmt"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

const screenshotPath = "./storage/screenshots"

func ScreenshotByURL(c *gin.Context) {
	// 從 POST 請求中讀取 URL 參數
	url := c.PostForm("url")

	// 建立新的上下文和超時
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 建立新的Chrome實例
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// 導航到指定的URL，等待頁面載入完成
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	); err != nil {
		log.Fatal(err)
	}

	// 獲取整個網頁的截圖
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 100)); err != nil {
		log.Fatal(err)
	}

	// 將截圖保存到本地文件
	filename := fmt.Sprintf("%s_%s.png", getDateTime(), randomString(10))
	filepath := filepath.Join(screenshotPath, filename)
	if err := os.WriteFile(filepath, buf, 0644); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
	})
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.New(rand.NewSource(time.Now().UnixNano()))

	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = charset[rand.Intn(len(charset))]
	}

	return string(bytes)
}

func getDateTime() string {
	return time.Now().Format("20060102150405")
}

func PreviewScreenshot(c *gin.Context) {
	// 從 URL 中讀取圖片文件名
	filename := c.Param("filename")

	// 構造文件路徑
	filepath := filepath.Join(screenshotPath, filename)

	// 打開文件
	file, err := os.Open(filepath)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer file.Close()

	// 讀取文件
	img, err := png.Decode(file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 設置 HTTP 響應頭
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "inline")

	// 將圖片寫入響應中
	err = png.Encode(c.Writer, img)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func DownloadScreenshot(c *gin.Context) {
	// 從 URL 中讀取圖片文件名
	filename := c.Param("filename")

	// 構造文件路徑
	filepath := filepath.Join(screenshotPath, filename)

	// 打開文件
	file, err := os.Open(filepath)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer file.Close()

	// 讀取文件
	img, err := png.Decode(file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 設置 HTTP 響應頭
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// 將圖片寫入響應中
	err = png.Encode(c.Writer, img)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
