CREATE TABLE users 
(
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE pastes (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,             -- публичный идентификатор 
    user_id BIGINT REFERENCES users(id),   -- null если paste анонимная

    title TEXT,
    content TEXT NOT NULL,

    syntax TEXT,                           -- "go", "js", "none" — подсветка
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    is_burn_after_read BOOLEAN NOT NULL DEFAULT FALSE,

    expire_at TIMESTAMPTZ,                 -- null → бессрочно
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pastes_uuid ON pastes(uuid);
CREATE INDEX idx_pastes_user_id ON pastes(user_id);
CREATE INDEX idx_pastes_expire_at ON pastes(expire_at);