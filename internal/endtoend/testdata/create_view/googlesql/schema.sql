CREATE TABLE venues (
  id INT64 NOT NULL,
  name STRING(MAX) NOT NULL,
  city STRING(MAX),
) PRIMARY KEY (id);
CREATE VIEW venue_names SQL SECURITY INVOKER AS
SELECT venues.id, venues.name FROM venues;
CREATE OR REPLACE VIEW venue_names SQL SECURITY INVOKER AS
SELECT venues.id, venues.name, venues.city FROM venues;
CREATE VIEW dropped SQL SECURITY INVOKER AS SELECT venues.id FROM venues;
DROP VIEW dropped;
