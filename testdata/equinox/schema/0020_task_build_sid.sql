ALTER TABLE assets ADD COLUMN releasesid varchar(64);
ALTER TABLE publishings ADD COLUMN releasesid varchar(64);
ALTER TABLE publishings ADD COLUMN channelsid varchar(64);
ALTER TABLE releases ADD COLUMN appsid varchar(64);
ALTER TABLE tasks ADD COLUMN buildsid varchar(64);
ALTER TABLE tasks ADD COLUMN rawassetsid varchar(64);
