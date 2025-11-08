CREATE TABLE IF NOT EXISTS followers (
  id UUID PRIMARY KEY,
  followers_id UUID NOT NULL,
  followee_id UUID NOT NULL,
  CONSTRAINT fk_followers
    FOREIGN KEY(followers_id)
      REFERENCES users(id),
  CONSTRAINT fk_followee
    FOREIGN KEY(followee_id)
      REFERENCES users(id)
);