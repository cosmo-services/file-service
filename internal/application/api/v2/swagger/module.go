package swagger_api

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewSwaggerRoutes),
)
