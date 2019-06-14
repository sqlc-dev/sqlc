-- Channels
CREATE OR REPLACE FUNCTION increment_channel_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE applications SET total_channels = total_channels+1 WHERE sid = NEW.appsid;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION decrement_channel_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE applications SET total_channels = total_channels-1 WHERE sid = OLD.appsid;
    RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';


-- Releases
CREATE OR REPLACE FUNCTION increment_release_count() RETURNS trigger AS $$
DECLARE
    counter INTEGER;
BEGIN
    UPDATE applications SET total_releases = total_releases+1 WHERE sid = NEW.appsid;
    UPDATE applications SET release_counter = release_counter+1 WHERE sid = NEW.appsid RETURNING release_counter INTO counter;
    NEW.counter_id := counter;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION decrement_release_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE applications SET total_releases = total_releases-1 WHERE sid = OLD.appsid;
    RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';

-- Assets
CREATE OR REPLACE FUNCTION increment_asset_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE releases SET total_assets = total_assets+1 WHERE sid = NEW.releasesid;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION decrement_asset_count() RETURNS trigger AS $$
DECLARE
BEGIN
    UPDATE releases SET total_assets = total_assets-1 WHERE sid = OLD.releasesid;
    RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';
