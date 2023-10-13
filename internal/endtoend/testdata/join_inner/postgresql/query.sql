-- name: SelectAllJoinedAlias :many
select e.* from events e
    inner join handled_events he
       on e.ID > he.last_handled_id
where he.handler = $1
    for update of he skip locked;

-- name: SelectAllJoined :many
select events.* from events
    inner join handled_events
       on events.ID > handled_events.last_handled_id
where handled_events.handler = $1
    for update of handled_events skip locked;
