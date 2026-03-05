package file_infrastructure

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"main/internal/config"
	domain "main/internal/domain/file"
)

type localFileStorage struct {
	config *config.FileStorageConfig
}

func NewLocalFileStorage(cfg *config.FileStorageConfig) domain.FileStorage {
	return &localFileStorage{
		config: cfg,
	}
}

func (s *localFileStorage) Save(f domain.File, directory string) (string, error) {
	dirInfo, err := s.config.GetDirectoryInfo(directory)
	if err != nil {
		return "", fmt.Errorf("invalid directory: %w", err)
	}

	fileName, err := s.generateFileName(f)
	if err != nil {
		return "", fmt.Errorf("failed to generate filename: %w", err)
	}

	fullPath := filepath.Join(dirInfo.Path, fileName)

	if err := os.MkdirAll(dirInfo.Path, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	if err := s.writeFile(fullPath, f); err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *localFileStorage) Delete(fileName string, directory string) error {
	dirInfo, err := s.config.GetDirectoryInfo(directory)
	if err != nil {
		return fmt.Errorf("invalid directory: %w", err)
	}

	fullPath := filepath.Join(dirInfo.Path, fileName)
	if err := s.validatePath(fullPath, dirInfo.Path); err != nil {
		return err
	}

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return domain.ErrFileNotFound
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *localFileStorage) Get(fileName string, directory string) (domain.File, error) {
	dirInfo, err := s.config.GetDirectoryInfo(directory)
	if err != nil {
		return nil, fmt.Errorf("invalid directory: %w", err)
	}

	fullPath := filepath.Join(dirInfo.Path, fileName)
	if err := s.validatePath(fullPath, dirInfo.Path); err != nil {
		return nil, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &localFile{
		File:      file,
		mimeType:  s.detectMimeType(fileName),
		closeFunc: file.Close,
	}, nil
}

func (s *localFileStorage) Exists(fileName string, directory string) (bool, error) {
	dirInfo, err := s.config.GetDirectoryInfo(directory)
	if err != nil {
		return false, fmt.Errorf("invalid directory: %w", err)
	}

	fullPath := filepath.Join(dirInfo.Path, fileName)
	if err := s.validatePath(fullPath, dirInfo.Path); err != nil {
		return false, err
	}

	_, err = os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check file: %w", err)
}

func (s *localFileStorage) generateFileName(f domain.File) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	ext := s.getExtension(f.MimeType())

	timestamp := time.Now().UnixNano()
	uuid := hex.EncodeToString(bytes)

	return fmt.Sprintf("%d_%s%s", timestamp, uuid, ext), nil
}

func (s *localFileStorage) getExtension(mimeType string) string {
	extMap := map[string]string{
		"image/jpeg":       ".jpg",
		"image/png":        ".png",
		"image/gif":        ".gif",
		"image/webp":       ".webp",
		"application/pdf":  ".pdf",
		"text/plain":       ".txt",
		"application/json": ".json",
	}

	if ext, ok := extMap[mimeType]; ok {
		return ext
	}
	return ".bin"
}

func (s *localFileStorage) detectMimeType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))

	mimeMap := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".json": "application/json",
	}

	if mime, ok := mimeMap[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}

func (s *localFileStorage) writeFile(path string, f domain.File) error {
	tempPath := path + ".tmp"

	dst, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, f); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	if err := dst.Sync(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to sync file: %w", err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}

func (s *localFileStorage) validatePath(fullPath string, basePath string) error {
	cleanFull := filepath.Clean(fullPath)
	cleanBase := filepath.Clean(basePath)

	if !strings.HasPrefix(cleanFull, cleanBase) {
		return os.ErrPermission
	}

	rel, err := filepath.Rel(cleanBase, cleanFull)
	if err != nil {
		return os.ErrPermission
	}

	if strings.HasPrefix(rel, "..") {
		return os.ErrPermission
	}

	return nil
}

type localFile struct {
	*os.File
	mimeType  string
	closeFunc func() error
}

func (f *localFile) MimeType() string {
	return f.mimeType
}

func (f *localFile) Close() error {
	if f.closeFunc != nil {
		return f.closeFunc()
	}
	return f.File.Close()
}
