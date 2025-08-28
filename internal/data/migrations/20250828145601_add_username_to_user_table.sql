-- +goose Up
-- +goose StatementBegin
ALTER TABLE Users ADD COLUMN username varchar(20);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Users DROP COLUMN username;
-- +goose StatementEnd
