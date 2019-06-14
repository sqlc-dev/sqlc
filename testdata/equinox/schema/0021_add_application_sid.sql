ALTER TABLE channels ADD COLUMN appsid varchar(64);

UPDATE channels 
SET appsid = applications.sid
FROM applications
WHERE channels.appid = applications.id AND channels.appsid IS NULL;
