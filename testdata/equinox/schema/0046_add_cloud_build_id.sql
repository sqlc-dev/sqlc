ALTER TABLE cloud_builds RENAME COLUMN id TO gid;
ALTER TABLE cloud_builds ADD COLUMN sid varchar(64);
-- Not a true source of randomness, but good enough for the small number of builds
UPDATE cloud_builds SET sid = 'cb_' || md5(random()::text) WHERE sid IS NULL;
ALTER TABLE cloud_builds ALTER COLUMN sid SET NOT NULL;
ALTER TABLE cloud_builds ADD UNIQUE (sid);
