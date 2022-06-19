// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
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
    CURRENT_TIMESTAMP,
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
	Status          string         `json:"status"`
	Statuses        sql.NullString `json:"statuses"`
	Tags            sql.NullString `json:"tags"`
}

func (q *Queries) CreateVenue(ctx context.Context, arg CreateVenueParams) (sql.Result, error) {
	ctx, done := q.observer(ctx, "CreateVenue")
	result, err := q.exec(ctx, q.createVenueStmt, createVenue,
		arg.Slug,
		arg.Name,
		arg.City,
		arg.SpotifyPlaylist,
		arg.Status,
		arg.Statuses,
		arg.Tags,
	)
	return result, done(err)
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
	ctx, done := q.observer(ctx, "DeleteVenue")
	_, err := q.exec(ctx, q.deleteVenueStmt, deleteVenue, arg.Slug, arg.Slug_2)
	return done(err)
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
	ctx, done := q.observer(ctx, "GetVenue")
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
	return i, done(err)
}

const listVenues = `-- name: ListVenues :many
SELECT id, status, statuses, slug, name, city, spotify_playlist, songkick_id, tags, created_at
FROM venue
WHERE city = ?
ORDER BY name
`

func (q *Queries) ListVenues(ctx context.Context, city string) ([]Venue, error) {
	ctx, done := q.observer(ctx, "ListVenues")
	rows, err := q.query(ctx, q.listVenuesStmt, listVenues, city)
	if err != nil {
		return nil, done(err)
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
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
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
	ctx, done := q.observer(ctx, "UpdateVenueName")
	_, err := q.exec(ctx, q.updateVenueNameStmt, updateVenueName, arg.Name, arg.Slug)
	return done(err)
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
	ctx, done := q.observer(ctx, "VenueCountByCity")
	rows, err := q.query(ctx, q.venueCountByCityStmt, venueCountByCity)
	if err != nil {
		return nil, done(err)
	}
	defer rows.Close()
	var items []VenueCountByCityRow
	for rows.Next() {
		var i VenueCountByCityRow
		if err := rows.Scan(&i.City, &i.Count); err != nil {
			return nil, done(err)
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, done(err)
	}
	if err := rows.Err(); err != nil {
		return nil, done(err)
	}
	return items, done(nil)
}
