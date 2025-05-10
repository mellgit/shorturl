-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID DEFAULT gen_random_uuid() NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    token TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    primary key(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
