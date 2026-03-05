package nats

import (
	"main/pkg"

	"go.uber.org/fx"
)

type Nats struct {
	natsClient           *pkg.NatsClient
	fileSubscribeHandler *FileSubscribeHandler
}

func NewNats(
	natsClient *pkg.NatsClient,
	fileSubscribeHandler *FileSubscribeHandler,
) *Nats {
	return &Nats{
		natsClient:           natsClient,
		fileSubscribeHandler: fileSubscribeHandler,
	}
}

func (n *Nats) SetupSubscribers() {
	n.natsClient.Subscribe("FILE_STREAM", "file.orphaned", n.fileSubscribeHandler.OnFileOrphaned)
}

var Module = fx.Options(
	fx.Provide(NewNats),
	fx.Provide(NewFileSubscribeHandler),
)
