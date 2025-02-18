CREATE TABLE books (
                       id   integer PRIMARY KEY,
                       title text      NOT NULL,
                       author text     NOT NULL,
                       pages integer   NOT NULL,
                       score float     NOT NULL,
                       price decimal(4, 2) NOT NULL,
                       avg_word_length double NOT NULL
);
