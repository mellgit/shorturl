-- +goose Up
-- +goose StatementBegin
CREATE TABLE clicks (
    id UUID DEFAULT gen_random_uuid() NOT NULL,
    alias TEXT NOT NULL,
    ip TEXT NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    primary key(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table clicks;
-- +goose StatementEnd
