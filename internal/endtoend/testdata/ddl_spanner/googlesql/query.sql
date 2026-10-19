-- name: ListSingers :many
SELECT * FROM singers;

-- name: ListAlbums :many
SELECT * FROM albums WHERE singer_id = @singer_id;

-- name: ListTickets :many
SELECT * FROM tickets;
