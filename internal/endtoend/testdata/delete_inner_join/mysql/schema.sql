CREATE TABLE author (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE book (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL
);

CREATE TABLE author_book (
  author_id INT UNSIGNED NOT NULL,
  book_id INT UNSIGNED NOT NULL,
  CONSTRAINT `pk-author_book` PRIMARY KEY (author_id, book_id),
  CONSTRAINT `fk-author_book-author-id` FOREIGN KEY (author_id) REFERENCES author (id),
  CONSTRAINT `fk-author_book-book-id` FOREIGN KEY (book_id) REFERENCES book (id)
);

