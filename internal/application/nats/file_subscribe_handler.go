package nats

import (
	"encoding/json"
	"main/internal/domain/file"
	"main/pkg"
	"strings"

	"github.com/nats-io/nats.go"
)

type FileSubscribeHandler struct {
	logger      pkg.Logger
	fileService *file.FileService
}

func NewFileSubscribeHandler(fileService *file.FileService, logger pkg.Logger) *FileSubscribeHandler {
	return &FileSubscribeHandler{
		fileService: fileService,
		logger:      logger,
	}
}

func (p *FileSubscribeHandler) OnFileOrphaned(msg *nats.Msg) error {
	var event FileOrphanedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		p.logger.Error(err)

		return err
	}

	parts := strings.Split(event.FilePath, "/")
	fileName := parts[len(parts)-1]
	fileDir := parts[len(parts)-2]

	if err := p.fileService.DeleteFile(fileName, fileDir); err != nil {
		p.logger.Error(err)

		return err
	}

	p.logger.Infof("File deleted by file event: %s", event)

	return nil
}
