package test_api

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewTestController),
	fx.Provide(NewTestRoutes),
)
