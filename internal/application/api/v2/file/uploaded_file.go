package file_api

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type uploadedFile struct {
	multipart.File
	size        int64
	mimeType    string
	initialized bool
}

func (f *uploadedFile) init() error {
	if f.initialized {
		return nil
	}

	seeker, ok := f.File.(io.Seeker)
	if !ok {
		return errors.New("failed to readfile:")
	}

	currentPos, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("failed to get current position: %w", err)
	}

	size, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek end: %w", err)
	}
	f.size = size

	_, err = seeker.Seek(currentPos, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek back: %w", err)
	}

	buffer := make([]byte, 512)
	n, err := f.File.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	f.mimeType = http.DetectContentType(buffer[:n])

	_, err = seeker.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek to start: %w", err)
	}

	f.initialized = true
	return nil
}

func (f *uploadedFile) Size() int64 {
	return f.size
}

func (f *uploadedFile) MimeType() string {
	return f.mimeType
}
