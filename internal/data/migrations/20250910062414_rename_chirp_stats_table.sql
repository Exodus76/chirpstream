-- +goose Up
-- +goose StatementBegin
ALTER TABLE Chirp_stats RENAME TO Chirp_likes;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Chirp_likes RENAME TO Chirp_stats;
-- +goose StatementEnd
