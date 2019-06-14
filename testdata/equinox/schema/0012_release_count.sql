ALTER TABLE applications ADD COLUMN release_counter int DEFAULT 0;
ALTER TABLE releases ADD COLUMN counter_id int DEFAULT 0;

UPDATE applications
SET release_counter = (
  SELECT COUNT(1) FROM releases
  WHERE releases.appid = applications.id
);

CREATE OR REPLACE FUNCTION increment_release_count() RETURNS trigger AS $$
DECLARE
    counter INTEGER;
BEGIN
    UPDATE applications SET total_releases = total_releases+1 WHERE id = NEW.appid;
    UPDATE applications SET release_counter = release_counter+1 WHERE id = NEW.appid RETURNING release_counter INTO counter;
    NEW.counter_id := counter;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
