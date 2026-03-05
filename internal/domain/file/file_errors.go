package file

import "errors"

var (
	ErrFileTypeNotAllowed = errors.New("FILE_TYPE_NOT_ALLOWED")
	ErrFileNotFound       = errors.New("FILE_NOT_FOUND")
)
