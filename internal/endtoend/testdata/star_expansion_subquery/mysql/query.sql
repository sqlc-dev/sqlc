/* name: StarExpansionSubquery :many */
SELECT * FROM foo WHERE EXISTS (SELECT * FROM foo);
