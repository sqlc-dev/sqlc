ATTACH 'baz.db' AS baz;

CREATE TABLE users (
    id integer PRIMARY KEY,
    name text NOT NULL,
    age integer
);

CREATE TABLE posts (
    id integer PRIMARY KEY,
    user_id integer NOT NULL
);

CREATE TABLE baz.users (
    id integer PRIMARY KEY,
    name text NOT NULL
);


CREATE TABLE user_links (
  owner_id integer NOT NULL,
  consumer_id integer NOT NULL,
  PRIMARY KEY (owner_id, consumer_id)
);


