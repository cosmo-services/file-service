package cmd

import (
	"main/bootstrap"

	"context"

	"go.uber.org/fx"

	"main/pkg"

	"main/internal/application/api/v2"
	"main/internal/application/jobs"
	"main/internal/application/nats"
	"main/internal/config"
)

func SetupApp(
	lc fx.Lifecycle,
	env config.Env,
	logger pkg.Logger,
	handler pkg.RequestHandler,
	routes api.Routes,
	nats *nats.Nats,
	workers jobs.Workers,
) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(startCtx context.Context) error {
			go func() {
				nats.SetupSubscribers()
				nats.SetupPublishers()
			}()

			go func() {
				workers.Run(ctx)
			}()

			go func() {
				routes.Setup()

				if err := handler.Gin.Run(":" + env.Port); err != nil {
					logger.Error(err)
				}
			}()

			return nil
		},
		OnStop: func(stopCtx context.Context) error {
			cancel()

			return nil
		},
	})
}

func StartApp() {
	opts := fx.Options(
		fx.Invoke(SetupApp),
	)

	app := fx.New(
		bootstrap.CommonModules,
		opts,
	)

	app.Run()
}
