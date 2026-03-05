package file_api

import (
	"main/internal/application/api/v2/auth"
	"main/internal/config"
	"main/pkg"
)

type FileRoutes struct {
	handler        pkg.RequestHandler
	fileController *FileController
	fileCfg        *config.FileStorageConfig
	authMiddleware *auth.AuthMiddleware
}

func NewFileRoutes(
	fileController *FileController,
	handler pkg.RequestHandler,
	fileCfg *config.FileStorageConfig,
	authMiddleware *auth.AuthMiddleware,
) *FileRoutes {
	return &FileRoutes{
		fileController: fileController,
		fileCfg:        fileCfg,
		handler:        handler,
		authMiddleware: authMiddleware,
	}
}

func (r *FileRoutes) Setup() {
	v2 := r.handler.Gin.Group("/api/v2")

	v2.GET("/file/:directory/:filename", r.authMiddleware.OptionalAuth(), r.fileController.GetFile)

	files := v2.Group("/file")
	files.Use(r.authMiddleware.RequireAuth())
	{

		files.POST("/avatar", r.fileController.UploadAvatar)
		//files.DELETE("/:file_id", r.fileController.DeleteFile)
	}
}
