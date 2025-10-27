-- name: SelectAllJoined :many
SELECT events.* FROM events
    INNER JOIN handled_events
       ON events.ID = handled_events.last_handled_id
WHERE handled_events.handler = $handler;

-- name: SelectAllJoinedAlias :many
SELECT e.* FROM events AS e
    INNER JOIN handled_events AS he
       ON e.ID = he.last_handled_id
WHERE he.handler = $handler;




