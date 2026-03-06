package nats

import "time"

type FileOrphanedEvent struct {
	FilePath   string    `json:"file_path"`
	OrphanedAt time.Time `json:"orphaned_at"`
}
