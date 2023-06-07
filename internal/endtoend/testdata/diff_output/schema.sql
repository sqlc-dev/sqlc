CREATE TABLE authors (
          id   BIGSERIAL PRIMARY KEY,
          name text      NOT NULL,
          bio  text
);

CREATE TABLE books (
          id    BIGSERIAL PRIMARY KEY,
          title text      NOT NULL
);

