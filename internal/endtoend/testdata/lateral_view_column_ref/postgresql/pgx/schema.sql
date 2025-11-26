CREATE TABLE foo (
  id integer PRIMARY KEY,
  val text NOT NULL
);

-- Reproduces issue #4182: LATERAL subquery referencing outer column
CREATE VIEW foo_lateral AS
SELECT t.val, sub.result
FROM foo t
CROSS JOIN LATERAL (
  SELECT t.val AS result
  FROM foo LIMIT 1
) sub;
