-- +goose Up
-- +goose StatementBegin
CREATE TABLE collection (
     id SERIAL PRIMARY KEY,
     user_id BIGINT NOT NULL,
     name VARCHAR(255) NOT NULL,
     pinned BOOLEAN NOT NULL DEFAULT false,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_collections_user_id ON collection(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS collection;
-- +goose StatementEnd
