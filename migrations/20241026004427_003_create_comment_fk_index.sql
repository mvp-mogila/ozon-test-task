-- +goose Up
-- +goose StatementBegin
CREATE INDEX comment_parent_id_idx ON comment(parent_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS comment_parent_id_idx;
-- +goose StatementEnd
