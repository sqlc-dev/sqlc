CREATE TABLE users (
  id    INTEGER PRIMARY KEY,
  name  TEXT NOT NULL,
  bio   TEXT,
  score REAL NOT NULL DEFAULT 0
);

CREATE TABLE posts (
  id      INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  title   TEXT,
  created TEXT NOT NULL
);

CREATE INDEX posts_user ON posts(user_id);
