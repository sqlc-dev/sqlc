CREATE TABLE authors (
  id   int PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);
CREATE TABLE book1 (
  author_id  int PRIMARY KEY,
  name text,
  FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`) ON DELETE CASCADE
);
CREATE TABLE book2 (
  author_id  int PRIMARY KEY,
  name text,
  FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`) ON DELETE CASCADE
);
