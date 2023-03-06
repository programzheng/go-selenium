package screenshot

import (
	"context"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

const screenshotPath = "./storage/screenshots"

func fullScreenshot(url string) (*[]byte, error) {
	// 建立新的上下文和超時
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// 建立新的Chrome實例
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// 導航到指定的URL，等待頁面載入完成
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	); err != nil {
		return nil, err
	}

	// 獲取整個網頁的截圖
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 100)); err != nil {
		return nil, err
	}

	return &buf, nil
}

func SaveScreenshotByURL(c *gin.Context) {
	// 從 POST 請求中讀取 URL 參數
	url := c.PostForm("url")

	buf, err := fullScreenshot(url)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 將截圖保存到本地文件
	filename := fmt.Sprintf("%s_%s.jpeg", getDateTime(), randomString(10))
	filepath := filepath.Join(screenshotPath, filename)
	if err := os.WriteFile(filepath, *buf, 0644); err != nil {
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

func DynamicPreviewScreenshot(c *gin.Context) {
	// 取得截圖的網址參數
	url := c.Param("url")

	// 將網址轉換為編碼過的字串
	decodedURL, err := base64.URLEncoding.DecodeString(url)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("解碼失敗：%v", err))
		return
	}

	// 產生網頁截圖
	buf, err := fullScreenshot(string(decodedURL))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 設定回應的 MIME 類型
	c.Header("Content-Type", "image/png")

	// 將圖片資料輸出到回應內容中
	c.Data(http.StatusOK, "image/png", *buf)
}

func PreviewByStaticScreenshotFile(c *gin.Context) {
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
	img, err := jpeg.Decode(file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 設置 HTTP 響應頭
	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Disposition", "inline")

	// 將圖片寫入響應中
	err = jpeg.Encode(c.Writer, img, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func DownloadStaticScreenshotFile(c *gin.Context) {
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
	img, err := jpeg.Decode(file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// 設置 HTTP 響應頭
	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// 將圖片寫入響應中
	err = jpeg.Encode(c.Writer, img, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
