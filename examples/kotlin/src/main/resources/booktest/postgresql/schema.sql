DROP TABLE IF EXISTS books CASCADE;
DROP TYPE IF EXISTS book_type CASCADE;
DROP TABLE IF EXISTS authors CASCADE;
DROP FUNCTION IF EXISTS say_hello(text) CASCADE;

CREATE TABLE authors (
          author_id SERIAL PRIMARY KEY,
          name text NOT NULL DEFAULT ''
);

CREATE INDEX authors_name_idx ON authors(name);

CREATE TYPE book_type AS ENUM (
          'FICTION',
          'NONFICTION'
);

CREATE TABLE books (
          book_id SERIAL PRIMARY KEY,
          author_id integer NOT NULL REFERENCES authors(author_id),
          isbn text NOT NULL DEFAULT '' UNIQUE,
          booktype book_type NOT NULL DEFAULT 'FICTION',
          title text NOT NULL DEFAULT '',
          year integer NOT NULL DEFAULT 2000,
          available timestamp with time zone NOT NULL DEFAULT 'NOW()',
          tags varchar[] NOT NULL DEFAULT '{}'
);

CREATE INDEX books_title_idx ON books(title, year);

CREATE FUNCTION say_hello(text) RETURNS text AS $$
BEGIN
          RETURN CONCAT('hello ', $1);
END;
$$ LANGUAGE plpgsql;

CREATE INDEX books_title_lower_idx ON books(title);
