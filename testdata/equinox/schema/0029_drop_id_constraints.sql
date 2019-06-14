ALTER TABLE releases DROP CONSTRAINT app_id_version;
ALTER TABLE channels DROP CONSTRAINT app_id_name;
ALTER TABLE publishings DROP CONSTRAINT channel_id_release_id;
