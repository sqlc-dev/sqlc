-- name: GetAllParties :many
SELECT * FROM party;

-- name: GetAllPeople :many
SELECT * FROM person;

-- name: GetAllOrganisations :many
SELECT * FROM organisation;

-- name: GetOrganizationsByRegion :many
SELECT *
FROM organisation
WHERE
	region = 'us' AND
	rank = 'ensign';
