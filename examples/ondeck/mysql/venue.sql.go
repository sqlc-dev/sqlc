// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: venue.sql

package ondeck

import (
	"context"
	"database/sql"
)

const createVenue = `-- name: CreateVenue :execresult
INSERT INTO venue (
    slug,
    name,
    city,
    created_at,
    spotify_playlist,
    status,
    statuses,
    tags
) VALUES (
    ?,
    ?,
    ?,
    NOW(),
    ?,
    ?,
    ?,
    ?
)
`

type CreateVenueParams struct {
	Slug            string         `json:"slug"`
	Name            string         `json:"name"`
	City            string         `json:"city"`
	SpotifyPlaylist string         `json:"spotify_playlist"`
	Status          VenuesStatus   `json:"status"`
	Statuses        sql.NullString `json:"statuses"`
	Tags            sql.NullString `json:"tags"`
}

func (q *Queries) CreateVenue(ctx context.Context, arg CreateVenueParams) (sql.Result, error) {
	return q.exec(ctx, q.createVenueStmt, createVenue,
		arg.Slug,
		arg.Name,
		arg.City,
		arg.SpotifyPlaylist,
		arg.Status,
		arg.Statuses,
		arg.Tags,
	)
}

const deleteVenue = `-- name: DeleteVenue :exec
DELETE FROM venue
WHERE slug = ? AND slug = ?
`

type DeleteVenueParams struct {
	Slug   string `json:"slug"`
	Slug_2 string `json:"slug_2"`
}

func (q *Queries) DeleteVenue(ctx context.Context, arg DeleteVenueParams) error {
	_, err := q.exec(ctx, q.deleteVenueStmt, deleteVenue, arg.Slug, arg.Slug_2)
	return err
}

const getVenue = `-- name: GetVenue :one
SELECT id, status, statuses, slug, name, city, spotify_playlist, songkick_id, tags, created_at
FROM venue
WHERE slug = ? AND city = ?
`

type GetVenueParams struct {
	Slug string `json:"slug"`
	City string `json:"city"`
}

func (q *Queries) GetVenue(ctx context.Context, arg GetVenueParams) (Venue, error) {
	row := q.queryRow(ctx, q.getVenueStmt, getVenue, arg.Slug, arg.City)
	var i Venue
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Statuses,
		&i.Slug,
		&i.Name,
		&i.City,
		&i.SpotifyPlaylist,
		&i.SongkickID,
		&i.Tags,
		&i.CreatedAt,
	)
	return i, err
}

const listVenues = `-- name: ListVenues :many
SELECT id, status, statuses, slug, name, city, spotify_playlist, songkick_id, tags, created_at
FROM venue
WHERE city = ?
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
		if err := rows.Scan(
			&i.ID,
			&i.Status,
			&i.Statuses,
			&i.Slug,
			&i.Name,
			&i.City,
			&i.SpotifyPlaylist,
			&i.SongkickID,
			&i.Tags,
			&i.CreatedAt,
		); err != nil {
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

const updateVenueName = `-- name: UpdateVenueName :exec
UPDATE venue
SET name = ?
WHERE slug = ?
`

type UpdateVenueNameParams struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (q *Queries) UpdateVenueName(ctx context.Context, arg UpdateVenueNameParams) error {
	_, err := q.exec(ctx, q.updateVenueNameStmt, updateVenueName, arg.Name, arg.Slug)
	return err
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
	City  string `json:"city"`
	Count int64  `json:"count"`
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
