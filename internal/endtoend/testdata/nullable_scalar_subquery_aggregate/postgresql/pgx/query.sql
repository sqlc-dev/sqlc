-- name: AggregateUnderCast :many
SELECT t.id, (SELECT count(*)::int FROM trackers k WHERE k.truck_id = t.id) AS trackers
FROM trucks t;

-- name: AggregateInExpression :many
SELECT t.id, (SELECT count(*)::int + 1 FROM trackers k WHERE k.truck_id = t.id) AS trackers
FROM trucks t;

-- name: GroupedAggregate :many
SELECT t.id, (SELECT count(*) FROM trackers k WHERE k.truck_id = t.id GROUP BY k.active LIMIT 1) AS trackers
FROM trucks t;

-- name: PlainFunction :many
SELECT t.id, (SELECT upper(d.name)::text FROM depots d WHERE d.id = t.depot_id) AS depot
FROM trucks t;

-- name: AggregateInNestedSubquery :many
SELECT t.id, (SELECT d.id::int + (SELECT count(*)::int FROM trackers) FROM depots d WHERE d.id = t.depot_id) AS n
FROM trucks t;

-- name: WindowFunction :many
SELECT t.id, (SELECT count(*) OVER () FROM depots d WHERE d.id = t.depot_id) AS n
FROM trucks t;
