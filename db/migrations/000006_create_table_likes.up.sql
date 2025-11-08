CREATE TABLE IF NOT EXISTS likes (
  id UUID PRIMARY KEY,
  post_id UUID NOT NULL,
  user_id UUID NOT NULL,
  created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_likes_post
    FOREIGN KEY(post_id)
      REFERENCES posts(id),
  CONSTRAINT fk_likes_user
    FOREIGN KEY(user_id)
      REFERENCES users(id)
);