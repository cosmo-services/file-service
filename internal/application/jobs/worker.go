package jobs

import (
	"context"

	"go.uber.org/fx"
)

type Worker interface {
	Run(ctx context.Context)
}

type Workers []Worker

func NewWorkers(
	testWorker *TestWorker,
) Workers {
	return Workers{
		testWorker,
	}
}

func (w Workers) Run(ctx context.Context) {
	for _, worker := range w {
		go worker.Run(ctx)
	}
}

var Module = fx.Options(
	fx.Provide(NewWorkers),
	fx.Provide(NewTestWorker),
)
