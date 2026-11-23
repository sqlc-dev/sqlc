                                                                         
                       
                     
           
                

-- name: Shapes :many
SELECT
  (a * 2),
  (t.a * 2),
  ("t"."a" * 2),
  (42 * 2),
  (1.5 * 2),
  ('7' * 2),
  ($1 * 2),
  ($2 * 2),
  ($3 * 2),
  (abs(a) * 2),
  ((a + 1) * 2),
  ((a + 1) * 2),
  ((a::int) * 2),
  ((greatest(a, 1) + least(a, 2)) * 2),
  ((- a) * 2)
FROM t;
