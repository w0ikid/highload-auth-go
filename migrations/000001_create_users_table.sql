-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    -- Используем UUID вместо обычного ID для безопасности и масштабируемости
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Индексируемое поле для логина. 
    -- COLLATE "C" ускоряет текстовый поиск в Postgres
    email      VARCHAR(255) UNIQUE NOT NULL,
    
    -- Хэш пароля (Argon2id занимает около 90-100 символов)
    password_hash VARCHAR(255) NOT NULL,
    
    -- Статус аккаунта (bit или smallint экономят место)
    is_active  BOOLEAN DEFAULT true,
    
    -- Таймстампы для аудита
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Важный индекс для быстрого логина
CREATE INDEX idx_users_email_active ON users(email) WHERE is_active = true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email_active;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
