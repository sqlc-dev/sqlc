ALTER TABLE channels DROP CONSTRAINT "channels_appid_fkey";
ALTER TABLE channels ALTER COLUMN appid DROP NOT NULL;
ALTER TABLE releases DROP CONSTRAINT "releases_appid_fkey";
ALTER TABLE releases ALTER COLUMN appid DROP NOT NULL;
