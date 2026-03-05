package file_api

import (
	"errors"
	"main/internal/config"
	file_domain "main/internal/domain/file"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService *file_domain.FileService
	fileCfg     *config.FileStorageConfig
}

func NewFileController(
	fileService *file_domain.FileService,
	fileCfg *config.FileStorageConfig,
) *FileController {
	return &FileController{
		fileService: fileService,
		fileCfg:     fileCfg,
	}
}

func (c *FileController) GetFile(ctx *gin.Context) {
	directory := ctx.Param("directory")
	filePath := ctx.Param("filepath")

	if directory == "" || filePath == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "directory and file path are required",
		})
		return
	}

	dirInfo, err := c.fileCfg.GetDirectoryInfo(directory)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "directory not found",
		})
		return
	}

	switch dirInfo.Access {
	case "public":
		fullPath := filepath.Join(dirInfo.Path, filePath)
		ctx.File(fullPath)

	case "private":
		userId := ctx.GetString("user_id")
		if userId == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required for private files",
			})
			return
		}

		file, err := c.fileService.GetFile(userId, filePath, directory)
		if err != nil {
			if errors.Is(err, file_domain.ErrFileNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			} else if errors.Is(err, file_domain.ErrNoAccess) {
				ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.DataFromReader(http.StatusOK, file.Size(), file.MimeType(), file, nil)

	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid access type in config",
		})
	}
}
