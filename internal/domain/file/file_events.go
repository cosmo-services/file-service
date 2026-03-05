package file

import "time"

type AvatarUploadedEvent struct {
	Directory  string    `json:"directory"`
	FileName   string    `json:"file_name"`
	UserID     string    `json:"user_id"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type UserFileDeletedEvent struct {
	Directory string    `json:"directory"`
	FileName  string    `json:"file_name"`
	UserID    string    `json:"user_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
