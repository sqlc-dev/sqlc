/* name: CreateT1 :execresult */
INSERT INTO
  t1 (user_id, name)
VALUES
  (?, ?);

/* name: CreateT2 :execresult */
INSERT INTO
  t2 (email, name)
VALUES
  (?, ?);

/* name: CreateT3 :execresult */
INSERT INTO
  t3 (user_id, email)
VALUES
  (?, ?);

/* name: UpdateAll :exec */
UPDATE
  t1
  INNER JOIN t3 ON t3.user_id = t1.user_id
  INNER JOIN t2 ON t2.email = t3.email
SET
  t1.name = t2.name
WHERE
  t1.name = '';

/* name: GetT1 :one */
SELECT
  *
FROM
  t1
WHERE
  user_id = ?
LIMIT
  1;