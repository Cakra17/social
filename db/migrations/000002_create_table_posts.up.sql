CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY,
  caption TEXT NULL,
  media BYTEA DEFAULT NULL,
  user_id UUID NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user_post
    FOREIGN KEY(user_id)
      REFERENCES users(id)
);