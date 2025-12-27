-- name: DropTrigger :exec
DROP TRIGGER IF EXISTS my_trigger ON accounts;

-- name: DropEventTrigger :exec
DROP EVENT TRIGGER IF EXISTS my_event_trigger;
