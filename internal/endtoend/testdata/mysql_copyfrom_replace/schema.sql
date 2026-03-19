CREATE TABLE locations (
  id        VARCHAR(512) PRIMARY KEY,
  name      TEXT NOT NULL,
  address   TEXT NOT NULL,
  latitude  FLOAT,
  longitude FLOAT
);
