-- +goose Up
-- +goose StatementBegin
CREATE TABLE clicks (
    id SERIAL PRIMARY KEY,
    alias TEXT NOT NULL,
    ip TEXT NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table clicks;
-- +goose StatementEnd
