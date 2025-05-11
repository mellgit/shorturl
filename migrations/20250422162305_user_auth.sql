-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID DEFAULT gen_random_uuid() NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    primary key(id)
);

CREATE TABLE refresh_tokens (
    id UUID DEFAULT gen_random_uuid() NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    primary key(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table refresh_tokens;
drop table users;
-- +goose StatementEnd
