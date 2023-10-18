-- name: ListCaseIntentHistory :many
WITH RECURSIVE descendants AS
   ( SELECT case_intent_parent_id AS parent, case_intent_id AS child, 1 AS lvl
     FROM case_intent_parent_join
     UNION ALL
     SELECT d.parent as parent, p.case_intent_id as child, d.lvl + 1 as lvl
     FROM descendants d
              JOIN case_intent_parent_join p
                   ON d.child = p.case_intent_parent_id
   )
select distinct child, 'child' group_
from descendants
where parent = @case_intent_id
union
select distinct parent, 'parent' group_
from descendants
where child = @case_intent_id
ORDER BY child;
