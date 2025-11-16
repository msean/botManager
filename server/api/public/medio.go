package public

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msean/botmanager/server/global"
)

type MedioApi struct{}

func (api *MedioApi) UploadMedia(c *gin.Context) {

	// 允许上传的类型
	allowExt := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".mp4": true, ".mov": true, ".avi": true, ".mkv": true,
	}

	// 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传文件"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowExt[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持该文件类型: " + ext})
		return
	}

	// 本地存储目录（物理路径）
	saveDir := global.GVA_CONFIG.Local.StorePath
	if saveDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器未配置存储路径"})
		return
	}

	// URL 访问前缀
	viewDir := global.GVA_CONFIG.Local.Path
	if viewDir == "" {
		viewDir = "uploads/file"
	}

	// 创建目录
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败: " + err.Error()})
		return
	}

	if err := os.Chmod(saveDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "设置目录权限失败: " + err.Error()})
		return
	}

	// 新文件名，避免重复
	newName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(saveDir, newName)

	// 保存文件到本地
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	if err := os.Chmod(savePath, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "设置文件权限失败: " + err.Error()})
		return
	}

	// 访问 URL
	domain := strings.TrimRight(global.GVA_CONFIG.System.Domain, "/")
	fileURL := fmt.Sprintf("%s/%s/%s", domain, strings.Trim(viewDir, "/"), newName)

	c.JSON(http.StatusOK, gin.H{
		"msg": "上传成功",
		"url": fileURL,
	})
}
