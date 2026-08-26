CREATE TABLE authors (
  id integer PRIMARY KEY,
  name text NOT NULL,
  bio text
);

CREATE TABLE "Events" (
  id integer PRIMARY KEY,
  "EventName" text NOT NULL,
  "order" integer NOT NULL
);
