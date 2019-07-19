package ondeck

import (
	"context"
	"database/sql"
	"time"
)

type Status string

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "closed"
)

type City struct {
	Slug string
	Name string
}

type Venue struct {
	ID              int
	Status          Status
	Slug            string
	Name            string
	City            string
	SpotifyPlaylist string
	SongkickID      sql.NullString
	CreatedAt       time.Time
}

type dbtx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db dbtx
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

const createCity = `-- name: CreateCity :one
INSERT INTO city (
    name,
    slug
) VALUES (
    $1,
    $2
) RETURNING slug, name
`

type CreateCityParams struct {
	Name string
	Slug string
}

func (q *Queries) CreateCity(ctx context.Context, arg CreateCityParams) (City, error) {
	row := q.db.QueryRowContext(ctx, createCity, arg.Name, arg.Slug)
	var i City
	err := row.Scan(&i.Slug, &i.Name)
	return i, err
}

const createVenue = `-- name: CreateVenue :one
INSERT INTO venue (
    slug,
    name,
    city,
    created_at,
    spotify_playlist,
    status
) VALUES (
    $1,
    $2,
    $3,
    NOW(),
    $4,
    $5
) RETURNING id
`

type CreateVenueParams struct {
	Slug            string
	Name            string
	City            string
	SpotifyPlaylist string
	Status          Status
}

func (q *Queries) CreateVenue(ctx context.Context, arg CreateVenueParams) (int, error) {
	row := q.db.QueryRowContext(ctx, createVenue, arg.Slug, arg.Name, arg.City, arg.SpotifyPlaylist, arg.Status)
	var id int
	err := row.Scan(&id)
	return id, err
}

const deleteVenue = `-- name: DeleteVenue :exec
DELETE FROM venue
WHERE slug = $1 AND slug = $1
`

func (q *Queries) DeleteVenue(ctx context.Context, slug string) error {
	_, err := q.db.ExecContext(ctx, deleteVenue, slug)
	return err
}

const getCity = `-- name: GetCity :one
SELECT slug, name
FROM city
WHERE slug = $1
`

func (q *Queries) GetCity(ctx context.Context, slug string) (City, error) {
	row := q.db.QueryRowContext(ctx, getCity, slug)
	var i City
	err := row.Scan(&i.Slug, &i.Name)
	return i, err
}

const getVenue = `-- name: GetVenue :one
SELECT id, status, slug, name, city, spotify_playlist, songkick_id, created_at
FROM venue
WHERE slug = $1 AND city = $2
`

type GetVenueParams struct {
	Slug string
	City string
}

func (q *Queries) GetVenue(ctx context.Context, arg GetVenueParams) (Venue, error) {
	row := q.db.QueryRowContext(ctx, getVenue, arg.Slug, arg.City)
	var i Venue
	err := row.Scan(&i.ID, &i.Status, &i.Slug, &i.Name, &i.City, &i.SpotifyPlaylist, &i.SongkickID, &i.CreatedAt)
	return i, err
}

const listCities = `-- name: ListCities :many
SELECT slug, name
FROM city
ORDER BY name
`

func (q *Queries) ListCities(ctx context.Context) ([]City, error) {
	rows, err := q.db.QueryContext(ctx, listCities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []City
	for rows.Next() {
		var i City
		if err := rows.Scan(&i.Slug, &i.Name); err != nil {
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

const listVenues = `-- name: ListVenues :many
SELECT id, status, slug, name, city, spotify_playlist, songkick_id, created_at
FROM venue
WHERE city = $1
ORDER BY name
`

func (q *Queries) ListVenues(ctx context.Context, city string) ([]Venue, error) {
	rows, err := q.db.QueryContext(ctx, listVenues, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Venue
	for rows.Next() {
		var i Venue
		if err := rows.Scan(&i.ID, &i.Status, &i.Slug, &i.Name, &i.City, &i.SpotifyPlaylist, &i.SongkickID, &i.CreatedAt); err != nil {
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

const updateCityName = `-- name: UpdateCityName :exec
UPDATE city
SET name = $2
WHERE slug = $1
`

type UpdateCityNameParams struct {
	Slug string
	Name string
}

func (q *Queries) UpdateCityName(ctx context.Context, arg UpdateCityNameParams) error {
	_, err := q.db.ExecContext(ctx, updateCityName, arg.Slug, arg.Name)
	return err
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

func (q *Queries) UpdateVenueName(ctx context.Context, arg UpdateVenueNameParams) (int, error) {
	row := q.db.QueryRowContext(ctx, updateVenueName, arg.Slug, arg.Name)
	var id int
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
	Count int
}

func (q *Queries) VenueCountByCity(ctx context.Context) ([]VenueCountByCityRow, error) {
	rows, err := q.db.QueryContext(ctx, venueCountByCity)
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
