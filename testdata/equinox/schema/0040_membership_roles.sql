-- TODO: Eventually move these to enums
ALTER TABLE membership ADD COLUMN role TEXT DEFAULT 'admin';
UPDATE membership
SET role = 'owner'
FROM account
WHERE membership.account_id = account.id AND membership.user_id = account.owner_id;
ALTER TABLE account DROP COLUMN owner_id;
