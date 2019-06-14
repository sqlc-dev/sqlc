ALTER TABLE applications ADD COLUMN max_channels int DEFAULT 10;
ALTER TABLE applications ADD COLUMN total_channels int DEFAULT 0;
ALTER TABLE applications ADD CONSTRAINT application_channel_limit
  CHECK (total_channels <= max_channels);

UPDATE applications
SET total_channels = (
  SELECT COUNT(1) FROM channels
  WHERE channels.appid = applications.id
);

CREATE OR REPLACE FUNCTION increment_channel_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE applications SET total_channels = total_channels+1 WHERE id = NEW.appid;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER increment_channel
    BEFORE INSERT ON channels
    FOR EACH ROW
    EXECUTE PROCEDURE increment_channel_count();

CREATE OR REPLACE FUNCTION decrement_channel_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE applications SET total_channels = total_channels-1 WHERE id = OLD.appid;
    RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER decrement_channel
    BEFORE DELETE ON channels
    FOR EACH ROW
    EXECUTE PROCEDURE decrement_channel_count();
