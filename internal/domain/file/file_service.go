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

	fileDirectory := s.getDirectoryByAccessType(AccessTypePublic)
	fileName, err := s.storage.Save(file, fileDirectory)
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
		s.storage.Delete(fileName, fileDirectory)
		return nil, err
	}

	return meta, nil
}

func (s *FileService) DeleteFile(fileName string) error {
	fileDirectory := s.getDirectoryByAccessType(AccessTypePublic)
	storageErr := s.storage.Delete(fileName, fileDirectory)
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
	case FileTypeImage:
		return []string{"jpg", "jpeg", "png", "gif", "bmp"}
	}
	return nil
}

func (s *FileService) getDirectoryByAccessType(accessType AccessType) string {
	switch accessType {
	case AccessTypePrivate:
		return "private"
	case AccessTypePublic:
		return "public"
	}
	return ""
}
