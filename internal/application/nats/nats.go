package nats

import (
	"main/internal/domain"
	"main/pkg"

	"go.uber.org/fx"
)

type Nats struct {
	natsClient           *pkg.NatsClient
	eventBus             *domain.EventBus
	fileEventHandler     *FileEventHandler
	fileSubscribeHandler *FileSubscribeHandler
}

func NewNats(
	natsClient *pkg.NatsClient,
	eventBus *domain.EventBus,
	fileEventHandler *FileEventHandler,
	fileSubscribeHandler *FileSubscribeHandler,
) *Nats {
	return &Nats{
		natsClient:           natsClient,
		eventBus:             eventBus,
		fileEventHandler:     fileEventHandler,
		fileSubscribeHandler: fileSubscribeHandler,
	}
}

func (n *Nats) SetupSubscribers() {
	n.natsClient.Subscribe("FILE_STREAM", "file.orphaned", n.fileSubscribeHandler.OnFileOrphaned)
}

func (n *Nats) SetupPublishers() {
	n.eventBus.On("user.avatar.uploaded", n.fileEventHandler.AvatarUploaded)
	n.eventBus.On("user.file.deleted", n.fileEventHandler.UserFileDeleted)
}

var Module = fx.Options(
	fx.Provide(NewNats),
	fx.Provide(NewFileSubscribeHandler),
)
