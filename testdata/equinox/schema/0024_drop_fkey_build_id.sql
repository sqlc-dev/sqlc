ALTER TABLE tasks DROP CONSTRAINT "tasks_buildid_fkey";
ALTER TABLE tasks DROP CONSTRAINT "tasks_rawassetid_fkey";
ALTER TABLE builds DROP CONSTRAINT "builds_channelid_fkey";
ALTER TABLE builds DROP CONSTRAINT "builds_releaseid_fkey";
