-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS comment (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content TEXT NOT NULL CHECK (char_length(content) <= 2000),
    post_id INT NOT NULL REFERENCES post(id),
    parent_id INT REFERENCES comment(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comment;
-- +goose StatementEnd
