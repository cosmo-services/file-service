package api

import (
	health_api "main/internal/application/api/v2/health"
	swagger_api "main/internal/application/api/v2/swagger"
	test_api "main/internal/application/api/v2/test"

	"go.uber.org/fx"
)

type Route interface {
	Setup()
}

type Routes []Route

func NewRoutes(
	healthRoutes *health_api.HealthRoutes,
	swaggerRoutes *swagger_api.SwaggerRoutes,
	testRoutes *test_api.TestRoutes,
) Routes {
	return Routes{
		healthRoutes,
		swaggerRoutes,
		testRoutes,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}

var Module = fx.Options(
	fx.Provide(NewRoutes),
)
