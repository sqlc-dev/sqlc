CREATE INDEX assets_get_by_checksum ON assets(checksum);
CREATE INDEX assets_get_by_spec ON assets(os, arch, goarm, archiveformat, releaseid);
