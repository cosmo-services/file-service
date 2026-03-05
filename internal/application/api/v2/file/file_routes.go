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

	v2.GET("/file/:directory/*filepath", r.authMiddleware.OptionalAuth(), r.fileController.GetFile)

	files := v2.Group("/file")
	files.Use(r.authMiddleware.RequireAuth())
	{

		//files.POST("/:directory", r.fileController.UploadFile)
		//files.DELETE("/:file_id", r.fileController.DeleteFile)
	}
}
