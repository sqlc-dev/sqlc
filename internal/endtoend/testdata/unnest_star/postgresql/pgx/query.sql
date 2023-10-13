-- name: GetPlanItems :many
SELECT p.plan_id, p.item_id
FROM (SELECT * FROM unnest(@ids::bigint[])) AS i(req_item_id),
LATERAL (
    SELECT plan_id, item_id
    FROM plan_items
    WHERE
        item_id = i.req_item_id AND
        (@after = 0 OR plan_id < @after)
    ORDER BY plan_id DESC
    LIMIT @limit_count
) p;