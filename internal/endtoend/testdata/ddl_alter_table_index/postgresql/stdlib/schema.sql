CREATE TABLE temp(a TEXT);

CREATE INDEX temp_idx ON temp(a);
ALTER INDEX temp_idx ATTACH PARTITION temp_partition_idx;

