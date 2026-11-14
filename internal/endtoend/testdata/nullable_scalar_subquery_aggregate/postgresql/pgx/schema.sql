CREATE TABLE depots (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  note text
);

CREATE TABLE trucks (
  id       BIGSERIAL PRIMARY KEY,
  depot_id bigint    NOT NULL,
  plate    text      NOT NULL
);

CREATE TABLE trackers (
  id       BIGSERIAL PRIMARY KEY,
  truck_id bigint    NOT NULL,
  active   boolean   NOT NULL,
  serial   text
);
