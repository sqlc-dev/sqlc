ALTER TABLE venues RENAME TO venue;
ALTER TABLE venue DROP COLUMN dropped;
ALTER TABLE venue ADD COLUMN created_at Timestamp;

