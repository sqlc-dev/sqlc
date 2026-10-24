CREATE TABLE items (
  id   INTEGER NOT NULL,
  info STRUCT(label VARCHAR, score DOUBLE, tags VARCHAR[])
);
