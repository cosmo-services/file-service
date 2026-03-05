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

func (s *FileService) UploadAvatar(userId string, file File) (*FileMeta, error) {
	meta, err := s.createFile(userId, file, FileTypeImage, "avatar")
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (s *FileService) GetFile(userId string, fileName string, fileDir string) (File, error) {
	accessType, err := s.storage.GetAccessType(fileDir)
	if err != nil {
		return nil, err
	}
	if accessType != AccessTypePublic {
		return nil, ErrNoAccess
	}

	file, err := s.storage.Get(fileName, fileDir)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *FileService) DeleteFile(userId string, fileName string, fileDir string) error {
	accessType, err := s.storage.GetAccessType(fileDir)
	if err != nil {
		return err
	}
	if accessType != AccessTypePublic {
		return ErrNoAccess
	}

	storageErr := s.storage.Delete(fileName, fileDir)
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
		return []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/bmp",
			"image/webp",
		}
	}
	return nil
}

func (s *FileService) createFile(userId string, file File, fileType FileType, fileDir string) (*FileMeta, error) {
	accessType, err := s.storage.GetAccessType(fileDir)
	if err != nil {
		return nil, err
	}

	mimeType := file.MimeType()
	if err := s.ValidateAllowedFileType(mimeType, fileType); err != nil {
		return nil, err
	}

	fileName, err := s.storage.Save(file, fileDir)
	if err != nil {
		return nil, err
	}

	meta := &FileMeta{
		FileName:   fileName,
		FileType:   fileType,
		MimeType:   mimeType,
		AccessType: accessType,
		Directory:  fileDir,
		UserId:     userId,
	}

	if err := s.metaRepo.Create(meta); err != nil {
		s.storage.Delete(fileName, fileDir)
		return nil, err
	}

	return meta, nil
}
