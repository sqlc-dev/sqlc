ALTER TABLE builds ADD COLUMN channelsid varchar(64);
ALTER TABLE builds ADD COLUMN releasesid varchar(64);

UPDATE builds 
SET releasesid = releases.sid
FROM releases
WHERE builds.releaseid = releases.id AND builds.releasesid IS NULL;

UPDATE builds 
SET channelsid = channels.sid
FROM channels
WHERE builds.channelid = channels.id AND builds.channelsid IS NULL;
