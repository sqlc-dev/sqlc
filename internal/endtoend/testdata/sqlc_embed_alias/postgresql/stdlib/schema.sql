CREATE SCHEMA IF NOT EXISTS baz;

CREATE TABLE users (
    id integer NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    age integer NULL
);

CREATE TABLE posts (
    id integer NOT NULL PRIMARY KEY,
    user_id integer NOT NULL,
    likes integer[] NOT NULL
);

CREATE TABLE baz.users (
    id integer NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE user_links (
  owner_id integer NOT NULL,
  consumer_id integer NOT NULL,
  PRIMARY KEY (owner_id, consumer_id)
);
