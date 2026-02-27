package jobs

import (
	"context"
	"main/pkg"
	"time"
)

type TestWorker struct {
	logger pkg.Logger
}

func NewTestWorker(
	logger pkg.Logger,
) *TestWorker {
	return &TestWorker{
		logger: logger,
	}
}

func (w *TestWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Token clear worker stoped.")
			return
		case <-ticker.C:
			w.logger.Info("Test worker is working hard!")
		}
	}
}
