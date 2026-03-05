package file_api

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewFileController),
	fx.Provide(NewFileRoutes),
)
