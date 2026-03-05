package file

import "time"

type FileType string

const (
	FileTypeImage FileType = "image"
)

type AccessType string

const (
	AccessTypePublic  AccessType = "public"
	AccessTypePrivate AccessType = "private"
	AccessTypeNone    AccessType = "none"
)

type FileMeta struct {
	FileName   string     `json:"file_name"`
	FileType   FileType   `json:"file_type"`
	AccessType AccessType `json:"access_type"`
	Directory  string     `json:"directory"`
	MimeType   string     `json:"mime_type"`
	UserId     string     `json:"user_id"`
	CreatedAt  time.Time  `json:"created_at"`
}

type FileMetaRepository interface {
	Create(fileMeta *FileMeta) error
	Delete(fileName string) error
	GetByName(fileName string) (*FileMeta, error)
	GetByUserId(userId string) ([]*FileMeta, error)
	Exists(fileName string) (bool, error)
}

type File interface {
	Read(p []byte) (n int, err error)
	Close() error
	Size() int64
	MimeType() string
}

type FileStorage interface {
	Save(file File, directory string) (fileName string, err error)
	Delete(fileName string, directory string) (err error)
	Get(fileName string, directory string) (file File, err error)
	Exists(fileName string, directory string) (exists bool, err error)
	GetAccessType(directory string) (AccessType, error)
}
