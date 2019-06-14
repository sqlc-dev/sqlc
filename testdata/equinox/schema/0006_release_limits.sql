ALTER TABLE applications ADD COLUMN max_releases int DEFAULT 1000;
ALTER TABLE applications ADD COLUMN total_releases int DEFAULT 0;
ALTER TABLE applications ADD CONSTRAINT application_release_limit
  CHECK (total_releases <= max_releases);

UPDATE applications
SET total_releases = (
  SELECT COUNT(1) FROM releases
  WHERE releases.appid = applications.id
);

CREATE OR REPLACE FUNCTION increment_release_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE applications SET total_releases = total_releases+1 WHERE id = NEW.appid;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER increment_release
    BEFORE INSERT ON releases
    FOR EACH ROW
    EXECUTE PROCEDURE increment_release_count();

CREATE OR REPLACE FUNCTION decrement_release_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE applications SET total_releases = total_releases-1 WHERE id = OLD.appid;
    RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER decrement_release
    BEFORE DELETE ON releases
    FOR EACH ROW
    EXECUTE PROCEDURE decrement_release_count();
