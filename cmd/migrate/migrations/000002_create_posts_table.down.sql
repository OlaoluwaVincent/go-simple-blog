DROP INDEX IF EXISTS idx_posts_search_vector;
ALTER TABLE posts DROP COLUMN IF EXISTS search_vector;

DROP TRIGGER IF EXISTS trg_posts_updated_at ON posts;
DROP INDEX IF EXISTS idx_posts_user_id;

DROP TABLE IF EXISTS posts;
