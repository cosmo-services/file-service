package file_infrastructure

import (
	"database/sql"
	"errors"
	"main/internal/domain/file"
	"main/pkg"
	"time"
)

type fileMetaRepo struct {
	db pkg.PostgresDB
}

func NewFileMetaRepository(db pkg.PostgresDB) file.FileMetaRepository {
	return &fileMetaRepo{
		db: db,
	}
}

func (r *fileMetaRepo) Create(fileMeta *file.FileMeta) error {
	if fileMeta.CreatedAt.IsZero() {
		fileMeta.CreatedAt = time.Now()
	}

	var created file.FileMeta
	err := r.db.QueryRow(
		CreateFileQuery,
		fileMeta.FileName,
		fileMeta.FileType,
		fileMeta.AccessType,
		fileMeta.Directory,
		fileMeta.MimeType,
		fileMeta.UserId,
		fileMeta.CreatedAt,
	).Scan(
		&created.FileName,
		&created.FileType,
		&created.AccessType,
		&created.Directory,
		&created.MimeType,
		&created.UserId,
		&created.CreatedAt,
	)

	if err != nil {
		return err
	}

	*fileMeta = created
	return nil
}

func (r *fileMetaRepo) Delete(fileName string) error {
	result, err := r.db.Exec(DeleteFileQuery, fileName)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return file.ErrFileNotFound
	}

	return nil
}

func (r *fileMetaRepo) GetByName(fileName string) (*file.FileMeta, error) {
	var meta file.FileMeta

	err := r.db.QueryRow(GetFileByNameQuery, fileName).Scan(
		&meta.FileName,
		&meta.FileType,
		&meta.AccessType,
		&meta.Directory,
		&meta.MimeType,
		&meta.UserId,
		&meta.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, file.ErrFileNotFound
		}
		return nil, err
	}

	return &meta, nil
}

func (r *fileMetaRepo) GetByUserId(userId string) ([]*file.FileMeta, error) {
	rows, err := r.db.Query(GetFilesByUserIDQuery, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metas []*file.FileMeta

	for rows.Next() {
		var meta file.FileMeta
		err := rows.Scan(
			&meta.FileName,
			&meta.FileType,
			&meta.AccessType,
			&meta.Directory,
			&meta.MimeType,
			&meta.UserId,
			&meta.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		metas = append(metas, &meta)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metas, nil
}

func (r *fileMetaRepo) Exists(fileName string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(CheckFileExistsQuery, fileName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
