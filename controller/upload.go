package controller

import (
	"api/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const basePath = "E:\\docker\\nginx-image-server\\images\\"

// 检查文件扩展名是否允许
func isAllowedExtension(file *multipart.FileHeader) bool {
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".ico"}
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	for _, ext := range allowedExtensions {
		if fileExtension == ext {
			return true
		}
	}
	return false
}

func UploadPic(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取上传文件
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: err.Error()})
		}

		// 检查文件大小是否超过 5MB
		const maxSize = 5 * 1024 * 1024
		if file.Size > maxSize {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(types.Error{Error: "File too large"})
		}

		// 限制扩展名
		if !isAllowedExtension(file) {
			return c.Status(fiber.StatusBadRequest).JSON(types.Error{Error: "File type not allowed"})
		}

		// 获取当前日期并格式化为 YYYYMMDD
		currentDate := time.Now().Format("20060102")

		// 创建保存文件的目录
		uploadDir := filepath.Join(basePath, currentDate)
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}

		// 使用 UUID 重命名文件
		newFileName := uuid.New().String() + filepath.Ext(file.Filename)
		filePath := filepath.Join(uploadDir, newFileName)

		// 保存文件到指定目录
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(types.Error{Error: err.Error()})
		}
		// 适配 Windows 和 Linux 路径
		relativePath := strings.Join(strings.Split(filePath, "\\")[4:], "/")
		url := os.Getenv("IMAGE_SERVER") + relativePath
		//log.Debugf("%v", os.Getenv("IMAGE_SERVER")+relativePath)
		// 相对路径
		return c.Status(fiber.StatusOK).JSON(types.Success{Message: "上传成功", Data: url})
	}
}
