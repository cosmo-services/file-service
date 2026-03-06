package nats

import (
	"context"
	"main/internal/domain"
	"main/pkg"
	"time"
)

type FileEventHandler struct {
	natsClient *pkg.NatsClient
	logger     pkg.Logger
}

func NewFileEventHandler(
	natsClient *pkg.NatsClient,
	logger pkg.Logger,
) *FileEventHandler {
	return &FileEventHandler{
		natsClient: natsClient,
		logger:     logger,
	}
}

func (p *FileEventHandler) AvatarUploaded(event domain.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := p.natsClient.PublishJSON(ctx, "file.avatar.uploaded", event)
	if err != nil {
		p.logger.Error(err)

		return err
	}

	p.logger.Infof("File event successfully delivered: %s", event)

	return nil
}

func (p *FileEventHandler) UserFileDeleted(event domain.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := p.natsClient.PublishJSON(ctx, "file.deleted", event)
	if err != nil {
		p.logger.Error(err)

		return err
	}

	p.logger.Infof("File event successfully delivered: %s", event)

	return nil
}
