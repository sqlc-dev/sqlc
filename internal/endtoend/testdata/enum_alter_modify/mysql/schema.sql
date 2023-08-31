-- similar to issue https://github.com/sqlc-dev/sqlc/issues/1503

CREATE TABLE authors (
  id   bigint primary key,
  status enum("ok", "init") default "init" not null
);

-- remove this alter to see the change in models.go
ALTER TABLE authors MODIFY status enum('init', 'done', 'canceled', 'processing', 'waiting') default "init" not null;