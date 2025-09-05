-- +goose Up
-- +goose StatementBegin
CREATE TABLE Chirps (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    user_id int NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Chirps;
-- +goose StatementEnd
