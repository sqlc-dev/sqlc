-- :name GetAccountBySlug
-- :result :one
SELECT *
FROM account
WHERE slug = $1;

-- :name GetAccountByID
-- :result :one
SELECT *
FROM account
WHERE id = $1;

-- :name GetAccountByUser
-- :result :one
SELECT *
FROM account
WHERE slug = $1
AND EXISTS (
	SELECT user_id
	FROM membership
	WHERE account_id = account.id AND user_id = $2
);

-- :name GetDefaultAccountForUser
-- :result :one
SELECT *
FROM account
WHERE id IN (
    SELECT account_id
	FROM "user", membership
	WHERE "user".default_membership_id = membership.id
    AND "user".id = $1
);
