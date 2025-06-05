-- +goose Up
-- +goose StatementBegin
CREATE TABLE collection_game (
    collection_id BIGINT NOT NULL,
    game_id BIGINT NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (collection_id, game_id)
);

CREATE INDEX idx_collection_games_collection_id ON collection_game(collection_id);
CREATE INDEX idx_collection_games_game_id ON collection_game(game_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS collection_game;
-- +goose StatementEnd
