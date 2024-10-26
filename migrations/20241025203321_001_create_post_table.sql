-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS post (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL CHECK (char_length(title) <= 100),
    content TEXT NOT NULL CHECK (char_length(content) <= 2000),
    allow_comments BOOL NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS post;
-- +goose StatementEnd
