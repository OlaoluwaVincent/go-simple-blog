CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    tags TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_posts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT uq_user_post_title
        UNIQUE (user_id, title)
);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);

CREATE TRIGGER trg_posts_updated_at
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Full-Text Search (safe for re-runs)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'posts'
        AND column_name = 'search_vector'
    ) THEN
        ALTER TABLE posts
        ADD COLUMN search_vector tsvector
        GENERATED ALWAYS AS (
            to_tsvector('english', coalesce(title, '') || ' ' || coalesce(content, ''))
        ) STORED;
    END IF;
END$$;

CREATE INDEX IF NOT EXISTS idx_posts_search_vector
ON posts
USING GIN (search_vector);
