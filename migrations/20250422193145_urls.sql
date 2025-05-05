-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id UUID DEFAULT gen_random_uuid() NOT NULL,
    user_id UUID NOT NULL,
    original TEXT NOT NULL,
    alias TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    primary key(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table urls;
-- +goose StatementEnd
