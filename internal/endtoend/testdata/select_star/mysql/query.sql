/* name: GetAll :many */
SELECT * FROM users;

/* name: GetIDAll :many */
SELECT * FROM (SELECT id FROM users) t;
