-- +goose Up
-- +goose StatementBegin
CREATE TABLE game (
   id SERIAL PRIMARY KEY,
   title VARCHAR(255) NOT NULL,
   description TEXT,
   genre VARCHAR(100),
   age VARCHAR(100),
   person VARCHAR(100),
   avg_time VARCHAR(100),
   difficulty VARCHAR(100),
   image VARCHAR(100),
   rules VARCHAR(100),
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS game;
-- +goose StatementEnd
