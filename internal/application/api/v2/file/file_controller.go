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

// GetFile retrieves a file by its name and directory
// @Summary Get file
// @Description Returns a file from the specified directory. Public files don't require authentication, private files require JWT token
// @Tags files
// @Accept json
// @Produce octet-stream
// @Param directory path string true "Logical directory (avatars, attachments, documents)"
// @Param filename path string true "File name"
// @Success 200 {file} file "File content"
// @Failure 400 {object} map[string]string "directory and file name are required"
// @Failure 401 {object} map[string]string "authentication required for private files"
// @Failure 403 {object} map[string]string "access denied"
// @Failure 404 {object} map[string]string "directory not found"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /{directory}/{filename} [get]
func (c *FileController) GetFile(ctx *gin.Context) {
	directory := ctx.Param("directory")
	fileName := ctx.Param("filename")

	if directory == "" || fileName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "directory and file name are required",
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
		fullPath := filepath.Join(dirInfo.Path, fileName)
		ctx.File(fullPath)

	case "private":
		userId := ctx.GetString("user_id")
		if userId == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required for private files",
			})
			return
		}

		file, err := c.fileService.GetFile(userId, fileName, directory)
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
		defer file.Close()

		ctx.DataFromReader(http.StatusOK, file.Size(), file.MimeType(), file, nil)

	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid access type in config",
		})
	}
}

// UploadAvatar handles avatar upload for a user
// @Summary Upload user avatar
// @Description Uploads and sets a new avatar for the authenticated user. Only image files are allowed.
// @Tags avatars
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Avatar image file (jpeg, png, gif)"
// @Success 200 {object} file_domain.FileMeta "Uploaded file metadata"
// @Failure 400 {object} map[string]string "no file provided or invalid file type"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 403 {object} map[string]string "access denied"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /avatar [post]
func (c *FileController) UploadAvatar(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	multipartFile, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}
	defer multipartFile.Close()

	file := &uploadedFile{
		File: multipartFile,
	}

	if err := file.init(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	result, err := c.fileService.UploadAvatar(userID, file)
	if err != nil {
		switch {
		case errors.Is(err, file_domain.ErrNoAccess):
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, file_domain.ErrFileTypeNotAllowed):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// DeleteFile deletes a file by its name and directory
// @Summary Delete file
// @Description Deletes a file from the specified directory. User must be the owner of the file.
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param directory path string true "Logical directory (avatars, attachments, documents)"
// @Param filename path string true "File name"
// @Success 200 {object} map[string]string "file deleted successfully"
// @Failure 400 {object} map[string]string "directory and file name are required"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 403 {object} map[string]string "access denied"
// @Failure 404 {object} map[string]string "file not found"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /{directory}/{filename} [delete]
func (c *FileController) DeleteFile(ctx *gin.Context) {
	directory := ctx.Param("directory")
	fileName := ctx.Param("filename")
	userID := ctx.GetString("user_id")

	if directory == "" || fileName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "directory and file name are required",
		})
		return
	}

	if err := c.fileService.DeleteFileByUser(userID, fileName, directory); err != nil {
		switch {
		case errors.Is(err, file_domain.ErrFileNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		case errors.Is(err, file_domain.ErrNoAccess):
			ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "file deleted successfully",
	})
}
