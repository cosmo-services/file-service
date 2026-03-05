CREATE TABLE IF NOT EXISTS files (
    file_name   VARCHAR(255) PRIMARY KEY,
    file_type   VARCHAR(50) NOT NULL,
    access_type VARCHAR(50) NOT NULL,
    mime_type   VARCHAR(100) NOT NULL,
    user_id     VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_files_user_id ON files(user_id);
CREATE INDEX idx_files_file_type ON files(file_type);
