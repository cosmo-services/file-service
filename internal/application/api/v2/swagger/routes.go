package swagger_api

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "main/docs"

	"main/pkg"
)

type SwaggerRoutes struct {
	handler pkg.RequestHandler
}

func NewSwaggerRoutes(handler pkg.RequestHandler) *SwaggerRoutes {
	return &SwaggerRoutes{
		handler: handler,
	}
}

func (r *SwaggerRoutes) Setup() {
	group := r.handler.Gin.Group("/api/v2/")

	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
