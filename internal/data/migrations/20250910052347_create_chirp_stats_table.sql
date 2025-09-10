-- +goose Up
-- +goose StatementBegin
CREATE TABLE Chirp_stats (
    user_id BIGINT NOT NULL REFERENCES Users(id) ON DELETE CASCADE,
    chirp_id BIGINT NOT NULL REFERENCES Chirps(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, chirp_id)
    -- CONSTRAINT chirpid FOREIGN KEY (chirp_id) REFERENCES Chirps(id),
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Chirp_stats;
-- +goose StatementEnd
