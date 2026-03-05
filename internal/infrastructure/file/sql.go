package file_infrastructure

const (
    CreateFileQuery = `
        INSERT INTO files (
            file_name, 
            file_type, 
            access_type, 
            mime_type, 
            user_id, 
            created_at
        ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING file_name, file_type, access_type, mime_type, user_id, created_at
    `

    DeleteFileQuery = `
        DELETE FROM files 
        WHERE file_name = $1
    `

    GetFileByNameQuery = `
        SELECT 
            file_name, 
            file_type, 
            access_type, 
            mime_type, 
            user_id, 
            created_at
        FROM files 
        WHERE file_name = $1
    `

    GetFilesByUserIDQuery = `
        SELECT 
            file_name, 
            file_type, 
            access_type, 
            mime_type, 
            user_id, 
            created_at
        FROM files 
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

    CheckFileExistsQuery = `
        SELECT EXISTS(
            SELECT 1 FROM files 
            WHERE file_name = $1
        )
    `
)