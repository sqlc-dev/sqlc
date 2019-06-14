ALTER TABLE releases ADD COLUMN max_assets int DEFAULT 50;
ALTER TABLE releases ADD COLUMN total_assets int DEFAULT 0;
ALTER TABLE releases ADD CONSTRAINT release_asset_limit
  CHECK (total_assets <= max_assets);

UPDATE releases
SET total_assets = (
  SELECT COUNT(1) FROM assets
  WHERE assets.releaseid = releases.id
);

CREATE OR REPLACE FUNCTION increment_asset_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE releases SET total_assets = total_assets+1 WHERE id = NEW.releaseid;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER increment_asset
    BEFORE INSERT ON assets
    FOR EACH ROW
    EXECUTE PROCEDURE increment_asset_count();

CREATE OR REPLACE FUNCTION decrement_asset_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE releases SET total_assets = total_assets-1 WHERE id = OLD.releaseid;
    RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER decrement_asset
    BEFORE DELETE ON assets
    FOR EACH ROW
    EXECUTE PROCEDURE decrement_asset_count();

