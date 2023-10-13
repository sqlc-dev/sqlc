-- https://github.com/sqlc-dev/sqlc/issues/604
CREATE TABLE users (
  user_id    INT PRIMARY KEY,
  city_id    INT -- nullable
);
CREATE TABLE cities (
  city_id    INT PRIMARY KEY,
  mayor_id   INT NOT NULL
);
CREATE TABLE mayors (
  mayor_id   INT PRIMARY KEY,
  full_name  TEXT NOT NULL
);

-- https://github.com/sqlc-dev/sqlc/issues/1334
CREATE TABLE authors (
  id        INT PRIMARY KEY,
  name      TEXT NOT NULL,
  parent_id INT -- nullable
);

CREATE TABLE super_authors (
  super_id        INT PRIMARY KEY,
  super_name      TEXT NOT NULL,
  super_parent_id INT -- nullable
);

-- https://github.com/sqlc-dev/sqlc/issues/1334
CREATE TABLE users_2 (
    user_id           INT PRIMARY KEY,
    user_nickname     VARCHAR(30) UNIQUE NOT NULL,
    user_email        TEXT UNIQUE        NOT NULL,
    user_display_name TEXT               NOT NULL,
    user_password     TEXT               ,
    user_google_id    TEXT UNIQUE        ,
    user_apple_id     TEXT UNIQUE        ,
    user_bio          VARCHAR(160)       NOT NULL DEFAULT '',
    user_created_at   TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_avatar_id    INT UNIQUE        
);

CREATE TABLE media (
    media_id         INT PRIMARY KEY,
    media_created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    media_hash       TEXT      NOT NULL,
    media_directory  TEXT      NOT NULL,
    media_author_id  INT       NOT NULL,
    media_width      INT       NOT NULL,
    media_height     INT       NOT NULL
);
