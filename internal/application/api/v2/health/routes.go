package health_api

import (
	"main/pkg"
)

type HealthRoutes struct {
	handler          pkg.RequestHandler
	healthController HealthController
}

func NewHealthRoutes(
	pingController HealthController,
	handler pkg.RequestHandler,
) *HealthRoutes {
	return &HealthRoutes{
		healthController: pingController,
		handler:          handler,
	}
}

func (r *HealthRoutes) Setup() {
	group := r.handler.Gin.Group("/api/v2/")

	group.GET("/health", r.healthController.Health)
}
