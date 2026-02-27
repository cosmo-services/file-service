package health_api

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthController),
	fx.Provide(NewHealthRoutes),
)
