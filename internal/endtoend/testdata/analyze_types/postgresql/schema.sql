CREATE EXTENSION hstore;
CREATE SCHEMA myschema;
CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy');
CREATE TYPE myschema.mood AS ENUM ('x', 'y');
CREATE DOMAIN posint AS integer CHECK (VALUE > 0);
CREATE DOMAIN shortname AS varchar(20) NOT NULL;
CREATE TYPE point2 AS (x float8, y float8);
CREATE TYPE floatrange AS RANGE (subtype = float8);

CREATE TABLE things (
  id      bigserial PRIMARY KEY,
  price   numeric(10,2) NOT NULL,
  amount  numeric,
  title   varchar(255),
  code    character varying(10),
  tag     char(5),
  raw     bpchar,
  count   pg_catalog.int4,
  ints    int[],
  grid    int[][],
  bounded int[3],
  words   text[] NOT NULL,
  prices  numeric(10,2)[],
  ts      timestamp(3),
  tstz    timestamptz NOT NULL,
  ttz     time with time zone,
  iv      interval,
  ivd     interval day to second,
  iv3     interval(3),
  m       mood,
  mm      myschema.mood,
  p       posint,
  sn      shortname,
  pt      point2,
  pts     point2[],
  fr      floatrange,
  ir      int4range,
  b       bit(8),
  vb      varbit(16),
  js      jsonb,
  u       uuid,
  h       hstore,
  moods   mood[],
  dp      double precision,
  ts2     timestamp without time zone
);
