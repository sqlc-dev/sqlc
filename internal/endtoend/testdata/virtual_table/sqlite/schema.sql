CREATE TABLE tbl(a INTEGER PRIMARY KEY, b TEXT, c TEXT, d TEXT, e INTEGER);

CREATE VIRTUAL TABLE tbl_ft USING fts5(b, c UNINDEXED, content='tbl', content_rowid='a');

CREATE VIRTUAL TABLE ft USING fts5(b);

CREATE TRIGGER tbl_ai AFTER INSERT ON tbl BEGIN
  INSERT INTO tbl_ft(rowid, b, c) VALUES (new.a, new.b, new.c);
END;

INSERT INTO tbl VALUES(1, 'xx yy cc', 't', 'a', 11);
INSERT INTO tbl VALUES(2, 'aa bb', 't', 'a', 22);

INSERT INTO ft VALUES('xx cc');
INSERT INTO ft VALUES('cc bb');

CREATE TABLE weather (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL
);


CREATE VIRTUAL TABLE weather_rtree USING rtree(
  id,
  min_lang, max_long,
  min_lat, max_lat
);

CREATE TRIGGER weather_insert
AFTER INSERT
  ON weather BEGIN
INSERT INTO
  weather_rtree (
    id,
    min_lang,
    max_long,
    min_lat,
    max_lat
  )
VALUES
  (
    NEW.id,
    NEW.latitude,
    NEW.latitude,
    NEW.longitude,
    NEW.longitude
  );

END;
