CREATE TABLE authors (
    author_id Serial,
    name Text NOT NULL DEFAULT '',
    PRIMARY KEY (author_id)
);

CREATE TABLE books (
    book_id Serial,
    author_id Int32 NOT NULL,
    isbn Text NOT NULL DEFAULT '',
    book_type Text NOT NULL DEFAULT 'FICTION',
    title Text NOT NULL DEFAULT '',
    year Int32 NOT NULL DEFAULT 2000,
    available Timestamp NOT NULL,
    tag Text NOT NULL DEFAULT '',
    PRIMARY KEY (book_id)
);

