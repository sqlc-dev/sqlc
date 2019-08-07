package prepared

import (
	"context"

	"github.com/lib/pq"
)

const createVenue = `-- name: CreateVenue :one
INSERT INTO venue (
    slug,
    name,
    city,
    created_at,
    spotify_playlist,
    status,
    tags
) VALUES (
    $1,
    $2,
    $3,
    NOW(),
    $4,
    $5,
    $6
) RETURNING id
`

type CreateVenueParams struct {
	Slug            string
	Name            string
	City            string
	SpotifyPlaylist string
	Status          Status
	Tags            []string
}

func (q *Queries) CreateVenue(ctx context.Context, arg CreateVenueParams) (int32, error) {
	row := q.queryRow(ctx, q.createVenueStmt, createVenue, arg.Slug, arg.Name, arg.City, arg.SpotifyPlaylist, arg.Status, pq.Array(arg.Tags))
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteVenue = `-- name: DeleteVenue :exec
DELETE FROM venue
WHERE slug = $1 AND slug = $1
`

func (q *Queries) DeleteVenue(ctx context.Context, slug string) error {
	_, err := q.exec(ctx, q.deleteVenueStmt, deleteVenue, slug)
	return err
}

const getVenue = `-- name: GetVenue :one
SELECT id, status, slug, name, city, spotify_playlist, songkick_id, tags, created_at
FROM venue
WHERE slug = $1 AND city = $2
`

type GetVenueParams struct {
	Slug string
	City string
}

func (q *Queries) GetVenue(ctx context.Context, arg GetVenueParams) (Venue, error) {
	row := q.queryRow(ctx, q.getVenueStmt, getVenue, arg.Slug, arg.City)
	var i Venue
	err := row.Scan(&i.ID, &i.Status, &i.Slug, &i.Name, &i.City, &i.SpotifyPlaylist, &i.SongkickID, pq.Array(&i.Tags), &i.CreatedAt)
	return i, err
}

const listVenues = `-- name: ListVenues :many
SELECT id, status, slug, name, city, spotify_playlist, songkick_id, tags, created_at
FROM venue
WHERE city = $1
ORDER BY name
`

func (q *Queries) ListVenues(ctx context.Context, city string) ([]Venue, error) {
	rows, err := q.query(ctx, q.listVenuesStmt, listVenues, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Venue
	for rows.Next() {
		var i Venue
		if err := rows.Scan(&i.ID, &i.Status, &i.Slug, &i.Name, &i.City, &i.SpotifyPlaylist, &i.SongkickID, pq.Array(&i.Tags), &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateVenueName = `-- name: UpdateVenueName :one
UPDATE venue
SET name = $2
WHERE slug = $1
RETURNING id
`

type UpdateVenueNameParams struct {
	Slug string
	Name string
}

func (q *Queries) UpdateVenueName(ctx context.Context, arg UpdateVenueNameParams) (int32, error) {
	row := q.queryRow(ctx, q.updateVenueNameStmt, updateVenueName, arg.Slug, arg.Name)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const venueCountByCity = `-- name: VenueCountByCity :many
SELECT
    city,
    count(*)
FROM venue
GROUP BY 1
ORDER BY 1
`

type VenueCountByCityRow struct {
	City  string
	Count int64
}

func (q *Queries) VenueCountByCity(ctx context.Context) ([]VenueCountByCityRow, error) {
	rows, err := q.query(ctx, q.venueCountByCityStmt, venueCountByCity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []VenueCountByCityRow
	for rows.Next() {
		var i VenueCountByCityRow
		if err := rows.Scan(&i.City, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
