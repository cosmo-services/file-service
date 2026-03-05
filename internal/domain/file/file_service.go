package file

import "slices"

type FileService struct {
	metaRepo FileMetaRepository
	storage  FileStorage
}

func NewFileService(
	metaRepo FileMetaRepository,
	storage FileStorage,
) *FileService {
	return &FileService{
		metaRepo: metaRepo,
		storage:  storage,
	}
}

func (s *FileService) CreateFile(file File, userId string, fileType FileType) (*FileMeta, error) {
	mimeType := file.MimeType()
	if err := s.ValidateAllowedFileType(mimeType, fileType); err != nil {
		return nil, err
	}

	fileName, err := s.storage.Save(file)
	if err != nil {
		return nil, err
	}

	meta := &FileMeta{
		FileName: fileName,
		FileType: fileType,
		MimeType: mimeType,
		UserId:   userId,
	}

	if err := s.metaRepo.Create(meta); err != nil {
		s.storage.Delete(fileName)
		return nil, err
	}

	return meta, nil
}

func (s *FileService) DeleteFile(fileName string) error {
	storageErr := s.storage.Delete(fileName)
	metaErr := s.metaRepo.Delete(fileName)

	if storageErr != nil {
		return storageErr
	}

	if metaErr != nil {
		return metaErr
	}

	return nil
}

func (s *FileService) ValidateAllowedFileType(mimeType string, fileType FileType) error {
	allowed := s.GetAllowedMimeType(fileType)
	if slices.Contains(allowed, mimeType) {
		return nil
	}
	return ErrFileTypeNotAllowed
}

func (s *FileService) GetAllowedMimeType(fileType FileType) []string {
	switch fileType {
	case FileTypeAvatar:
		return []string{"jpg", "jpeg", "png", "gif", "bmp"}
	}
	return nil
}
