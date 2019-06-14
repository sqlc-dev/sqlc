-- These commands will lock the table; this is not a problem because the tables are small
ALTER TABLE releases ADD CONSTRAINT app_sid_version UNIQUE(appsid, version);
ALTER TABLE channels ADD CONSTRAINT app_sid_name UNIQUE(appsid, name);
ALTER TABLE publishings ADD CONSTRAINT channel_sid_release_sid UNIQUE(channelsid, releasesid);
