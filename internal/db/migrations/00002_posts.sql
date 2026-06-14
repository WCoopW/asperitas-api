-- +goose Up
CREATE TYPE post_type AS ENUM ('link', 'text');

CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES users(id),
    type post_type NOT NULL,
    title TEXT NOT NULL,
    url TEXT,
    text TEXT, 
    category TEXT NOT NULL,
    score INT NOT NULL DEFAULT 0,
    views INT NOT NULL DEFAULT 0,
    upvote_percentage INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_category ON posts(category);

-- +goose Down
DROP TABLE IF EXISTS posts;
DROP TYPE IF EXISTS post_type;