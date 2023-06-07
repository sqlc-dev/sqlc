CREATE TABLE authors (
          author_id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
          name text NOT NULL
);

CREATE INDEX authors_name_idx ON authors(name);

CREATE TABLE books (
          book_id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
          author_id integer NOT NULL,
          isbn varchar(255) NOT NULL DEFAULT '' UNIQUE,
          book_type text NOT NULL DEFAULT 'FICTION',
          title text NOT NULL,
          yr integer NOT NULL DEFAULT 2000,
          available datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
          tag text NOT NULL,
          CHECK (book_type = 'FICTION' OR book_type = 'NONFICTION')
);

CREATE INDEX books_title_idx ON books(title, yr);
