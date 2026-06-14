-- +goose Up

CREATE TABLE votes (
    post_id   UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    value     SMALLINT NOT NULL CHECK (value IN (-1, 1)),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (post_id, user_id)
);
-- +goose Down
DROP TABLE IF EXISTS votes;